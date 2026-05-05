package file

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"server/global"
	"server/model"
	modelrequest "server/model/request"
	modelresponse "server/model/response"
	"server/service/configsvc"
	"server/service/oss"
	"server/service/storagesvc"
)

type FileService struct{}

var Default = &FileService{}

type fileReferenceCountRow struct {
	FileID uint
	Count  int64
}

const (
	FileDeleteModeConfigKey = "file_delete_mode"
	FileDeleteModeLogical   = "logical"
	FileDeleteModePhysical  = "physical"
)

// GetFileList 获取文件列表
func (s *FileService) GetFileList(page, pageSize int, name, ext string, referenced *bool) ([]model.SysFile, int64, error) {
	var files []model.SysFile
	var total int64

	db := global.DB.Model(&model.SysFile{}).Where("status = ?", 1)
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if ext != "" {
		exts := strings.Split(ext, ",")
		if len(exts) == 1 {
			db = db.Where("ext = ?", ext)
		} else {
			db = db.Where("ext IN ?", exts)
		}
	}

	if referenced == nil {
		db.Count(&total)
		err := db.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&files).Error
		if err != nil {
			return nil, 0, err
		}
		if err := s.fillFileReferenceCounts(files); err != nil {
			return nil, 0, err
		}
		return files, total, nil
	}

	var fileIDs []uint
	if err := db.Order("id desc").Pluck("id", &fileIDs).Error; err != nil {
		return nil, 0, err
	}
	if len(fileIDs) == 0 {
		return []model.SysFile{}, 0, nil
	}

	referenceCounts, err := s.getFileReferenceCounts(fileIDs)
	if err != nil {
		return nil, 0, err
	}

	filteredIDs := make([]uint, 0, len(fileIDs))
	for _, fileID := range fileIDs {
		isReferenced := referenceCounts[fileID] > 0
		if isReferenced == *referenced {
			filteredIDs = append(filteredIDs, fileID)
		}
	}

	total = int64(len(filteredIDs))
	if total == 0 {
		return []model.SysFile{}, 0, nil
	}

	start := (page - 1) * pageSize
	if start >= len(filteredIDs) {
		return []model.SysFile{}, total, nil
	}

	end := start + pageSize
	if end > len(filteredIDs) {
		end = len(filteredIDs)
	}
	pageIDs := filteredIDs[start:end]
	if err := global.DB.Where("id IN ?", pageIDs).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	fileMap := make(map[uint]model.SysFile, len(files))
	for _, file := range files {
		file.ReferenceCount = referenceCounts[file.ID]
		fileMap[file.ID] = file
	}

	ordered := make([]model.SysFile, 0, len(pageIDs))
	for _, fileID := range pageIDs {
		if file, ok := fileMap[fileID]; ok {
			ordered = append(ordered, file)
		}
	}

	return ordered, total, nil
}

// GetFileByID 根据ID获取文件
func (s *FileService) GetFileByID(id uint) (*model.SysFile, error) {
	var file model.SysFile
	if err := global.DB.First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// GetFileByMD5 根据MD5获取文件（用于秒传）
func (s *FileService) GetFileByMD5(md5 string) (*model.SysFile, error) {
	var file model.SysFile
	if err := global.DB.Where("md5 = ? AND status = ?", md5, 1).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// CreateFile 创建文件记录
func (s *FileService) CreateFile(file *model.SysFile) error {
	return global.DB.Create(file).Error
}

func (s *FileService) resolveFileStorage(file model.SysFile) (*model.StorageProfile, error) {
	if strings.TrimSpace(file.StorageType) != "" {
		return storagesvc.Default.GetStorageByType(model.StorageType(file.StorageType))
	}
	return storagesvc.Default.GetDefaultStorage()
}

func (s *FileService) createFileRecord(filename, key, url, md5 string, fileSize int64, uploaderID uint, storage *model.StorageProfile) *model.SysFile {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	return &model.SysFile{
		Name:          filename,
		Path:          key,
		URL:           url,
		Size:          fileSize,
		Ext:           ext,
		MimeType:      getMimeType(ext),
		MD5:           md5,
		StorageType:   string(storage.Type),
		StorageBucket: storage.GetBucketName(),
		UploaderID:    uploaderID,
		Status:        1,
	}
}

func (s *FileService) getDeleteMode() string {
	config, err := configsvc.Default.GetConfigByKey(FileDeleteModeConfigKey)
	if err != nil || strings.TrimSpace(config.Value) == "" {
		return FileDeleteModeLogical
	}
	if config.Value == FileDeleteModePhysical {
		return FileDeleteModePhysical
	}
	return FileDeleteModeLogical
}

func (s *FileService) ensureFileNotReferenced(file model.SysFile) error {
	var avatarCount int64
	if err := global.DB.Model(&model.SysUser{}).Where("avatar_file_id = ?", file.ID).Count(&avatarCount).Error; err != nil {
		return err
	}
	if avatarCount > 0 {
		return errors.New("文件正在被引用：用户头像正在使用，无法删除")
	}

	if configLabel, err := s.findConfigFileReference(file.ID); err != nil {
		return err
	} else if configLabel != "" {
		return fmt.Errorf("文件正在被引用：系统配置[%s]正在使用，无法删除", configLabel)
	}

	var messages []model.AIMessage
	if err := global.DB.Select("id", "file_ids").Where("file_ids IS NOT NULL AND file_ids <> ''").Find(&messages).Error; err != nil {
		return err
	}
	for _, msg := range messages {
		var fileIDs []uint
		if err := json.Unmarshal([]byte(msg.FileIDs), &fileIDs); err != nil {
			continue
		}
		for _, fileID := range fileIDs {
			if fileID == file.ID {
				return errors.New("文件正在被引用：AI对话附件正在使用，无法删除")
			}
		}
	}

	return nil
}

func parseConfigFileID(rawValue string) (uint, bool) {
	fileID, err := strconv.ParseUint(strings.TrimSpace(rawValue), 10, 64)
	if err != nil || fileID == 0 {
		return 0, false
	}
	return uint(fileID), true
}

func (s *FileService) findConfigFileReference(fileID uint) (string, error) {
	var configs []model.SysConfig
	if err := global.DB.Select("key", "value").
		Where("`key` IN ? AND value IS NOT NULL AND value <> ''", configsvc.ImageFileReferenceConfigKeys()).
		Find(&configs).Error; err != nil {
		return "", err
	}

	for _, config := range configs {
		boundFileID, ok := parseConfigFileID(config.Value)
		if !ok || boundFileID != fileID {
			continue
		}
		return configsvc.ImageFileReferenceLabel(config.Key), nil
	}

	return "", nil
}

func (s *FileService) addConfigFileReferenceCounts(fileIDSet map[uint]struct{}, referenceCounts map[uint]int64) error {
	var configs []model.SysConfig
	if err := global.DB.Select("key", "value").
		Where("`key` IN ? AND value IS NOT NULL AND value <> ''", configsvc.ImageFileReferenceConfigKeys()).
		Find(&configs).Error; err != nil {
		return err
	}

	for _, config := range configs {
		fileID, ok := parseConfigFileID(config.Value)
		if !ok {
			continue
		}
		if _, exists := fileIDSet[fileID]; !exists {
			continue
		}
		referenceCounts[fileID]++
	}

	return nil
}

func (s *FileService) fillFileReferenceCounts(files []model.SysFile) error {
	referenceCounts, err := s.getFileReferenceCounts(s.collectFileIDs(files))
	if err != nil {
		return err
	}
	for i := range files {
		files[i].ReferenceCount = referenceCounts[files[i].ID]
	}
	return nil
}

func (s *FileService) getFileReferenceCounts(fileIDs []uint) (map[uint]int64, error) {
	referenceCounts := make(map[uint]int64, len(fileIDs))
	if len(fileIDs) == 0 {
		return referenceCounts, nil
	}

	fileIDSet := make(map[uint]struct{}, len(fileIDs))
	for _, fileID := range fileIDs {
		fileIDSet[fileID] = struct{}{}
	}

	var avatarCounts []fileReferenceCountRow
	if err := global.DB.Model(&model.SysUser{}).
		Select("avatar_file_id AS file_id, COUNT(*) AS count").
		Where("avatar_file_id IN ?", fileIDs).
		Group("avatar_file_id").
		Scan(&avatarCounts).Error; err != nil {
		return nil, err
	}
	for _, item := range avatarCounts {
		referenceCounts[item.FileID] += item.Count
	}

	if err := s.addConfigFileReferenceCounts(fileIDSet, referenceCounts); err != nil {
		return nil, err
	}

	var messages []model.AIMessage
	if err := global.DB.Select("id", "file_ids").Where("file_ids IS NOT NULL AND file_ids <> ''").Find(&messages).Error; err != nil {
		return nil, err
	}
	for _, msg := range messages {
		var ids []uint
		if err := json.Unmarshal([]byte(msg.FileIDs), &ids); err != nil {
			continue
		}

		seenInMessage := make(map[uint]struct{}, len(ids))
		for _, fileID := range ids {
			if _, exists := fileIDSet[fileID]; !exists {
				continue
			}
			if _, duplicated := seenInMessage[fileID]; duplicated {
				continue
			}
			seenInMessage[fileID] = struct{}{}
			referenceCounts[fileID]++
		}
	}

	return referenceCounts, nil
}

func (s *FileService) cleanupChunkArtifacts(file model.SysFile) error {
	query := global.DB.Model(&model.SysFileChunk{})
	if strings.TrimSpace(file.MD5) != "" {
		query = query.Where("file_hash = ? OR storage_path = ?", file.MD5, file.Path)
	} else {
		query = query.Where("storage_path = ?", file.Path)
	}

	var chunks []model.SysFileChunk
	if err := query.Find(&chunks).Error; err != nil {
		return err
	}
	if len(chunks) == 0 {
		return nil
	}

	seenUploads := make(map[string]struct{})
	for _, chunk := range chunks {
		if strings.TrimSpace(chunk.UploadID) == "" {
			continue
		}

		storageType := model.StorageType(chunk.StorageType)
		var storage *model.StorageProfile
		var err error
		if storageType != "" {
			storage, err = storagesvc.Default.GetStorageByType(storageType)
		} else {
			storage, err = s.resolveFileStorage(file)
		}
		if err != nil {
			return err
		}

		cacheKey := strings.Join([]string{string(storage.Type), chunk.UploadID, chunk.StoragePath}, "|")
		if _, ok := seenUploads[cacheKey]; ok {
			continue
		}
		seenUploads[cacheKey] = struct{}{}

		client, err := oss.GetClient(storage)
		if err != nil {
			return err
		}
		if err := client.AbortMultipartUpload(context.Background(), chunk.StoragePath, chunk.UploadID); err != nil {
			return err
		}
	}

	chunkIDs := make([]uint, 0, len(chunks))
	for _, chunk := range chunks {
		chunkIDs = append(chunkIDs, chunk.ID)
	}
	return global.DB.Unscoped().Where("id IN ?", chunkIDs).Delete(&model.SysFileChunk{}).Error
}

func (s *FileService) clearSoftDeletedAvatarReferences(fileID uint) error {
	return global.DB.Unscoped().
		Model(&model.SysUser{}).
		Where("avatar_file_id = ? AND deleted_at IS NOT NULL", fileID).
		Update("avatar_file_id", nil).Error
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(id uint) error {
	var file model.SysFile
	if err := global.DB.First(&file, id).Error; err != nil {
		return err
	}

	if err := s.ensureFileNotReferenced(file); err != nil {
		return err
	}

	if s.getDeleteMode() == FileDeleteModeLogical {
		return global.DB.Model(&model.SysFile{}).Where("id = ?", id).Update("status", 0).Error
	}

	if err := s.cleanupChunkArtifacts(file); err != nil {
		return err
	}

	if err := s.clearSoftDeletedAvatarReferences(file.ID); err != nil {
		return err
	}

	storage, err := s.resolveFileStorage(file)
	if err != nil {
		return err
	}

	client, err := oss.GetClient(storage)
	if err != nil {
		return err
	}
	if err := client.Delete(context.Background(), file.Path); err != nil {
		return err
	}

	return global.DB.Unscoped().Delete(&model.SysFile{}, id).Error
}

// BatchDeleteFiles 批量删除文件
func (s *FileService) BatchDeleteFiles(ids []uint) (int, []string) {
	var successCount int
	var failedMsgs []string

	for _, id := range ids {
		if err := s.DeleteFile(id); err != nil {
			failedMsgs = append(failedMsgs, fmt.Sprintf("ID %d: %s", id, err.Error()))
		} else {
			successCount++
		}
	}

	return successCount, failedMsgs
}

// PreviewFileMigration 预览文件迁移
func (s *FileService) PreviewFileMigration(req modelrequest.FileMigrationRequest) (*modelresponse.FileMigrationResult, error) {
	precheck, err := s.buildFileMigrationPrecheck(req)
	if err != nil {
		return nil, err
	}
	return precheck.result, nil
}

// ExecuteFileMigration 执行文件迁移
func (s *FileService) ExecuteFileMigration(req modelrequest.FileMigrationRequest) (*modelresponse.FileMigrationResult, error) {
	precheck, err := s.buildFileMigrationPrecheck(req)
	if err != nil {
		return nil, err
	}

	result := s.newFileMigrationExecutionResult(precheck.result)

	for _, candidate := range precheck.candidates {
		action, message := s.resolveFileMigrationCandidateAction(candidate)
		if action != "PENDING" {
			result.Items = append(result.Items, s.toFileMigrationItem(candidate, action, message))
			continue
		}

		warning, err := s.executeFileMigrationCandidate(candidate)
		switch {
		case err != nil:
			result.FailedCount++
			result.Items = append(result.Items, s.toFileMigrationItem(candidate, "FAILED", err.Error()))
		case warning != "":
			result.MigratedCount++
			result.WarningCount++
			result.Items = append(result.Items, s.toFileMigrationItem(candidate, "WARNING", warning))
		default:
			result.MigratedCount++
			result.Items = append(result.Items, s.toFileMigrationItem(candidate, "MIGRATED", "迁移成功"))
		}
	}

	return result, nil
}

// StartFileMigrationTask 启动文件迁移任务
func (s *FileService) StartFileMigrationTask(req modelrequest.FileMigrationRequest) (*modelresponse.FileMigrationTaskStatus, error) {
	return fileMigrationTasks.start(req, func(task *fileMigrationTask) {
		s.runFileMigrationTask(task, req)
	})
}

// GetCurrentFileMigrationTask 获取当前文件迁移任务状态
func (s *FileService) GetCurrentFileMigrationTask() *modelresponse.FileMigrationTaskStatus {
	return fileMigrationTasks.snapshot()
}

func (s *FileService) runFileMigrationTask(task *fileMigrationTask, req modelrequest.FileMigrationRequest) {
	precheck, err := s.buildFileMigrationPrecheck(req)
	if err != nil {
		task.finishFailed(err)
		return
	}

	task.setPrecheck(precheck.result)
	task.markRunning()
	for _, candidate := range precheck.candidates {
		task.setCurrentFile(candidate.file.ID, candidate.file.Name)

		action, message := s.resolveFileMigrationCandidateAction(candidate)
		if action != "PENDING" {
			task.recordHandled(s.toFileMigrationItem(candidate, action, message), candidate.file.Size)
			continue
		}

		warning, err := s.executeFileMigrationCandidate(candidate)
		switch {
		case err != nil:
			task.recordHandled(s.toFileMigrationItem(candidate, "FAILED", err.Error()), candidate.file.Size)
		case warning != "":
			task.recordHandled(s.toFileMigrationItem(candidate, "WARNING", warning), candidate.file.Size)
		default:
			task.recordHandled(s.toFileMigrationItem(candidate, "MIGRATED", "迁移成功"), candidate.file.Size)
		}
	}

	status := task.snapshot()
	if status.FailedCount > 0 || status.WarningCount > 0 {
		task.finishSuccess("迁移完成，请检查结果详情")
		return
	}
	task.finishSuccess("迁移完成")
}

func (s *FileService) buildFileMigrationPrecheck(req modelrequest.FileMigrationRequest) (*fileMigrationPrecheck, error) {
	candidates, err := s.buildFileMigrationCandidates(req)
	if err != nil {
		return nil, err
	}

	result := &modelresponse.FileMigrationResult{
		TargetStorageType: req.TargetStorageType,
		Items:             make([]modelresponse.FileMigrationItem, 0, len(candidates)),
	}

	for _, candidate := range candidates {
		result.TotalCount++
		result.TotalSize += candidate.file.Size

		action, message := s.resolveFileMigrationCandidateAction(candidate)
		switch action {
		case "PENDING":
			result.PendingCount++
			result.PendingSize += candidate.file.Size
		case "CONFLICT":
			result.ConflictCount++
			result.ConflictSize += candidate.file.Size
		case "MISSING_SOURCE":
			result.MissingSourceCount++
			result.MissingSourceSize += candidate.file.Size
		default:
			result.SkippedCount++
			result.SkippedSize += candidate.file.Size
		}

		result.Items = append(result.Items, s.toFileMigrationItem(candidate, action, message))
	}

	return &fileMigrationPrecheck{
		result:     result,
		candidates: candidates,
	}, nil
}

func (s *FileService) newFileMigrationExecutionResult(precheck *modelresponse.FileMigrationResult) *modelresponse.FileMigrationResult {
	return &modelresponse.FileMigrationResult{
		TargetStorageType:  precheck.TargetStorageType,
		TotalCount:         precheck.TotalCount,
		TotalSize:          precheck.TotalSize,
		PendingCount:       precheck.PendingCount,
		PendingSize:        precheck.PendingSize,
		SkippedCount:       precheck.SkippedCount,
		SkippedSize:        precheck.SkippedSize,
		ConflictCount:      precheck.ConflictCount,
		ConflictSize:       precheck.ConflictSize,
		MissingSourceCount: precheck.MissingSourceCount,
		MissingSourceSize:  precheck.MissingSourceSize,
		Items:              make([]modelresponse.FileMigrationItem, 0, len(precheck.Items)),
	}
}

func (s *FileService) resolveFileMigrationCandidateAction(candidate fileMigrationCandidate) (string, string) {
	if candidate.conflictMessage != "" {
		return "CONFLICT", candidate.conflictMessage
	}
	if candidate.sourceMissing {
		return "MISSING_SOURCE", candidate.skipMessage
	}
	if candidate.skipMessage != "" {
		return "SKIP", candidate.skipMessage
	}
	return "PENDING", "待迁移"
}

func (s *FileService) buildFileMigrationCandidates(req modelrequest.FileMigrationRequest) ([]fileMigrationCandidate, error) {
	scope := s.normalizeFileMigrationScope(req.Scope)
	sourceStorageType := model.StorageType(strings.TrimSpace(req.SourceStorageType))
	targetStorageType := model.StorageType(strings.TrimSpace(req.TargetStorageType))
	if sourceStorageType == "" {
		return nil, errors.New("请选择源存储")
	}
	if targetStorageType == "" {
		return nil, errors.New("请选择目标存储")
	}
	sameTypeMode := sourceStorageType == targetStorageType
	if sameTypeMode && strings.TrimSpace(req.SourceConfig) == "" {
		return nil, errors.New("同类型迁移需要提供源存储配置")
	}
	if !s.isSupportedMigrationStorageType(sourceStorageType) {
		return nil, fmt.Errorf("不支持的源存储类型: %s", sourceStorageType)
	}
	if !s.isSupportedMigrationStorageType(targetStorageType) {
		return nil, fmt.Errorf("不支持的目标存储类型: %s", targetStorageType)
	}
	if scope == fileMigrationScopeSelected && len(req.IDs) == 0 {
		return nil, errors.New("请选择要迁移的文件")
	}

	targetStorage, err := storagesvc.Default.GetStorageByType(targetStorageType)
	if err != nil {
		return nil, fmt.Errorf("获取目标存储配置失败: %v", err)
	}

	targetClient, err := oss.GetClient(targetStorage)
	if err != nil {
		return nil, fmt.Errorf("创建目标存储客户端失败: %v", err)
	}

	var sourceConfigOverride *model.StorageProfile
	if sameTypeMode {
		profile, err := storagesvc.Default.BuildStorageProfile(sourceStorageType, req.SourceConfig)
		if err != nil {
			return nil, fmt.Errorf("解析源存储配置失败: %v", err)
		}
		if profile.GetBucketName() == targetStorage.GetBucketName() {
			return nil, errors.New("源桶与目标桶相同，无需迁移")
		}
		sourceConfigOverride = profile
	}

	files, err := s.collectFilesForMigration(req, scope, sourceStorageType)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.New("未找到可迁移的文件")
	}

	selectedPathCount := make(map[string]int, len(files))
	targetPaths := make([]string, 0, len(files))
	for _, file := range files {
		path := strings.TrimSpace(file.Path)
		if path == "" {
			continue
		}
		selectedPathCount[path]++
		targetPaths = append(targetPaths, path)
	}

	conflictPathMap := make(map[string]struct{})
	if len(targetPaths) > 0 {
		var existingFiles []model.SysFile
		if err := global.DB.
			Where("status = ? AND storage_type = ? AND path IN ? AND id NOT IN ?", 1, string(targetStorageType), targetPaths, s.collectFileIDs(files)).
			Find(&existingFiles).Error; err != nil {
			return nil, err
		}
		for _, file := range existingFiles {
			conflictPathMap[file.Path] = struct{}{}
		}
	}

	candidates := make([]fileMigrationCandidate, 0, len(files))
	for _, file := range files {
		candidate := fileMigrationCandidate{
			file:          file,
			targetStorage: targetStorage,
			targetURL:     targetClient.GetURL(file.Path),
		}

		path := strings.TrimSpace(file.Path)
		if path == "" {
			candidate.skipMessage = "文件路径为空，无法迁移"
			candidates = append(candidates, candidate)
			continue
		}

		if selectedPathCount[path] > 1 {
			candidate.conflictMessage = "选中的文件存在重复存储路径，无法安全迁移"
			candidates = append(candidates, candidate)
			continue
		}

		if _, ok := conflictPathMap[path]; ok {
			candidate.conflictMessage = "目标存储已存在其他文件占用相同路径，无法覆盖"
			candidates = append(candidates, candidate)
			continue
		}

		var sourceStorage *model.StorageProfile
		if sourceConfigOverride != nil {
			sourceStorage = sourceConfigOverride
		} else {
			sourceStorage, err = s.resolveFileStorage(file)
			if err != nil {
				candidate.skipMessage = fmt.Sprintf("解析源存储失败: %v", err)
				candidates = append(candidates, candidate)
				continue
			}
		}
		candidate.sourceStorage = sourceStorage

		if !sameTypeMode && sourceStorage.Type == targetStorageType {
			candidate.skipMessage = "文件已在目标存储，无需迁移"
			candidates = append(candidates, candidate)
			continue
		}

		exists, err := s.fileObjectExists(sourceStorage, file.Path)
		if err != nil {
			candidate.skipMessage = fmt.Sprintf("检查源文件失败: %v", err)
			candidates = append(candidates, candidate)
			continue
		}
		if !exists {
			candidate.skipMessage = "源文件不存在，无法迁移"
			candidate.sourceMissing = true
			candidates = append(candidates, candidate)
			continue
		}

		targetConflictMessage, err := s.resolveTargetObjectConflict(candidate)
		if err != nil {
			candidate.skipMessage = err.Error()
			candidates = append(candidates, candidate)
			continue
		}
		if targetConflictMessage != "" {
			candidate.conflictMessage = targetConflictMessage
		}

		candidates = append(candidates, candidate)
	}

	return candidates, nil
}

func (s *FileService) collectFilesForMigration(req modelrequest.FileMigrationRequest, scope string, sourceStorageType model.StorageType) ([]model.SysFile, error) {
	db := global.DB.Model(&model.SysFile{}).
		Where("status = ? AND storage_type = ?", 1, string(sourceStorageType))

	switch scope {
	case fileMigrationScopeSelected:
		db = db.Where("id IN ?", req.IDs)
	case fileMigrationScopeFilter:
		if strings.TrimSpace(req.Filters.Name) != "" {
			db = db.Where("name LIKE ?", "%"+strings.TrimSpace(req.Filters.Name)+"%")
		}
		if strings.TrimSpace(req.Filters.Ext) != "" {
			exts := strings.Split(req.Filters.Ext, ",")
			if len(exts) == 1 {
				db = db.Where("ext = ?", strings.TrimSpace(req.Filters.Ext))
			} else {
				db = db.Where("ext IN ?", exts)
			}
		}
	}

	var files []model.SysFile
	if err := db.Order("id asc").Find(&files).Error; err != nil {
		return nil, err
	}

	if req.Filters.Referenced != nil {
		referenceCounts, err := s.getFileReferenceCounts(s.collectFileIDs(files))
		if err != nil {
			return nil, err
		}

		filtered := make([]model.SysFile, 0, len(files))
		for _, file := range files {
			isReferenced := referenceCounts[file.ID] > 0
			if isReferenced == *req.Filters.Referenced {
				filtered = append(filtered, file)
			}
		}
		files = filtered
	}

	return files, nil
}

func (s *FileService) normalizeFileMigrationScope(scope string) string {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case fileMigrationScopeFilter:
		return fileMigrationScopeFilter
	case fileMigrationScopeSelected:
		return fileMigrationScopeSelected
	default:
		return fileMigrationScopeAll
	}
}

func (s *FileService) collectFileIDs(files []model.SysFile) []uint {
	ids := make([]uint, 0, len(files))
	for _, file := range files {
		ids = append(ids, file.ID)
	}
	return ids
}

func (s *FileService) fileObjectExists(storage *model.StorageProfile, key string) (bool, error) {
	client, err := oss.GetClient(storage)
	if err != nil {
		return false, err
	}
	return client.Exists(context.Background(), key)
}

func (s *FileService) resolveTargetObjectConflict(candidate fileMigrationCandidate) (string, error) {
	client, err := oss.GetClient(candidate.targetStorage)
	if err != nil {
		return "", fmt.Errorf("检查目标文件失败: %v", err)
	}

	exists, err := client.Exists(context.Background(), candidate.file.Path)
	if err != nil {
		return "", fmt.Errorf("检查目标文件失败: %v", err)
	}
	if !exists {
		return "", nil
	}

	if strings.TrimSpace(candidate.file.MD5) == "" {
		return "目标存储已存在同路径文件，且当前文件缺少 MD5，无法安全迁移", nil
	}

	matches, err := s.fileObjectMatchesMD5(context.Background(), client, candidate.file.Path, candidate.file.MD5)
	if err != nil {
		return "", fmt.Errorf("校验目标文件失败: %v", err)
	}
	if !matches {
		return "目标存储已存在同路径文件，但内容不一致", nil
	}

	return "", nil
}

func (s *FileService) executeFileMigrationCandidate(candidate fileMigrationCandidate) (string, error) {
	if candidate.sourceStorage == nil {
		return "", errors.New("源存储配置不存在")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	sourceClient, err := oss.GetClient(candidate.sourceStorage)
	if err != nil {
		return "", fmt.Errorf("创建源存储客户端失败: %v", err)
	}
	targetClient, err := oss.GetClient(candidate.targetStorage)
	if err != nil {
		return "", fmt.Errorf("创建目标存储客户端失败: %v", err)
	}

	path := candidate.file.Path
	targetExists, err := targetClient.Exists(ctx, path)
	if err != nil {
		return "", fmt.Errorf("检查目标文件失败: %v", err)
	}

	uploadedNewObject := false
	if targetExists {
		if strings.TrimSpace(candidate.file.MD5) == "" {
			return "", errors.New("目标存储已存在同路径文件，且当前文件缺少 MD5，无法安全迁移")
		}
		matches, err := s.fileObjectMatchesMD5(ctx, targetClient, path, candidate.file.MD5)
		if err != nil {
			return "", fmt.Errorf("校验目标文件失败: %v", err)
		}
		if !matches {
			return "", errors.New("目标存储已存在同路径文件，但内容不一致")
		}
	} else {
		sourceReader, err := sourceClient.Open(ctx, path)
		if err != nil {
			return "", fmt.Errorf("读取源文件失败: %v", err)
		}

		if err := targetClient.Upload(ctx, path, sourceReader, candidate.file.Size); err != nil {
			sourceReader.Close()
			return "", fmt.Errorf("上传目标文件失败: %v", err)
		}
		if err := sourceReader.Close(); err != nil {
			return "", fmt.Errorf("关闭源文件流失败: %v", err)
		}
		uploadedNewObject = true
	}

	if strings.TrimSpace(candidate.file.MD5) != "" {
		matches, err := s.fileObjectMatchesMD5(ctx, targetClient, path, candidate.file.MD5)
		if err != nil {
			if uploadedNewObject {
				_ = targetClient.Delete(context.Background(), path)
			}
			return "", fmt.Errorf("迁移后校验失败: %v", err)
		}
		if !matches {
			if uploadedNewObject {
				_ = targetClient.Delete(context.Background(), path)
			}
			return "", errors.New("迁移后文件校验失败")
		}
	}

	updateErr := global.DB.Transaction(func(tx *gorm.DB) error {
		updateQuery := tx.Model(&model.SysFile{}).
			Where("id = ? AND status = ? AND path = ? AND url = ?", candidate.file.ID, 1, candidate.file.Path, candidate.file.URL)
		if strings.TrimSpace(candidate.file.StorageType) != "" {
			updateQuery = updateQuery.Where("storage_type = ?", candidate.file.StorageType)
		}

		updateData := map[string]interface{}{
			"storage_type":   string(candidate.targetStorage.Type),
			"storage_bucket": candidate.targetStorage.GetBucketName(),
			"url":            candidate.targetURL,
		}

		result := updateQuery.Updates(updateData)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("文件记录已变化，请刷新后重试")
		}

		// 同步更新引用了该文件的系统配置 URL
		fileIDStr := strconv.FormatUint(uint64(candidate.file.ID), 10)
		for fileIDKey, urlKey := range configsvc.ImageFileIDToURLKeyMap() {
			var cfg model.SysConfig
			if err := tx.Where("`key` = ? AND `value` = ?", fileIDKey, fileIDStr).First(&cfg).Error; err != nil {
				continue
			}
			tx.Model(&model.SysConfig{}).Where("`key` = ?", urlKey).Update("value", candidate.targetURL)
		}

		return nil
	})
	if updateErr != nil {
		if uploadedNewObject {
			_ = targetClient.Delete(context.Background(), path)
		}
		return "", updateErr
	}

	return "", nil
}

func (s *FileService) isSupportedMigrationStorageType(storageType model.StorageType) bool {
	for _, item := range storagesvc.Default.SupportedStorageTypes() {
		if item == storageType {
			return true
		}
	}
	return false
}

func (s *FileService) fileObjectMatchesMD5(ctx context.Context, client oss.Client, key, expectedMD5 string) (bool, error) {
	reader, err := client.Open(ctx, key)
	if err != nil {
		return false, err
	}
	defer reader.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return false, err
	}

	actualMD5 := fmt.Sprintf("%x", hash.Sum(nil))
	return strings.EqualFold(actualMD5, strings.TrimSpace(expectedMD5)), nil
}

func (s *FileService) toFileMigrationItem(candidate fileMigrationCandidate, action, message string) modelresponse.FileMigrationItem {
	sourceStorageType := ""
	if candidate.sourceStorage != nil {
		sourceStorageType = string(candidate.sourceStorage.Type)
	}

	return modelresponse.FileMigrationItem{
		FileID:            candidate.file.ID,
		FileName:          candidate.file.Name,
		SourceStorageType: sourceStorageType,
		TargetStorageType: string(candidate.targetStorage.Type),
		OldURL:            candidate.file.URL,
		NewURL:            candidate.targetURL,
		Action:            action,
		Message:           message,
	}
}

// GenerateFilePath 生成文件存储路径
func (s *FileService) GenerateFilePath(filename string) string {
	ext := filepath.Ext(filename)
	now := time.Now()
	return fmt.Sprintf("%d/%02d/%02d/%d%s", now.Year(), now.Month(), now.Day(), now.UnixNano(), ext)
}

// GetUploadCredential 获取上传凭证
func (s *FileService) GetUploadCredential(filename string) (*oss.UploadCredential, error) {
	storage, err := storagesvc.Default.GetDefaultStorage()
	if err != nil {
		return nil, fmt.Errorf("获取存储配置失败: %v", err)
	}

	client, err := oss.GetClient(storage)
	if err != nil {
		return nil, fmt.Errorf("创建存储客户端失败: %v", err)
	}

	key := s.GenerateFilePath(filename)
	credential, err := client.GetUploadCredential(context.Background(), key, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("获取上传凭证失败: %v", err)
	}

	return credential, nil
}

// CheckFileMD5 检查文件MD5是否存在（秒传）
func (s *FileService) CheckFileMD5(md5 string) (*model.SysFile, bool) {
	file, err := s.GetFileByMD5(md5)
	if err != nil {
		return nil, false
	}
	return file, true
}

// InitMultipartUpload 初始化分片上传
func (s *FileService) InitMultipartUpload(filename, md5 string, fileSize int64) (*oss.MultipartUpload, *model.StorageProfile, error) {
	storage, err := storagesvc.Default.GetDefaultStorage()
	if err != nil {
		return nil, nil, fmt.Errorf("获取存储配置失败: %v", err)
	}

	client, err := oss.GetClient(storage)
	if err != nil {
		return nil, nil, fmt.Errorf("创建存储客户端失败: %v", err)
	}

	key := s.GenerateFilePath(filename)
	upload, err := client.InitMultipartUpload(context.Background(), key)
	if err != nil {
		return nil, nil, fmt.Errorf("初始化分片上传失败: %v", err)
	}

	return upload, storage, nil
}

// GetMultipartUploadURLs 获取分片上传URL列表
func (s *FileService) GetMultipartUploadURLs(uploadID, key string, totalParts int, storage *model.StorageProfile) ([]string, error) {
	client, err := oss.GetClient(storage)
	if err != nil {
		return nil, fmt.Errorf("创建存储客户端失败: %v", err)
	}

	urls := make([]string, totalParts)
	for i := 1; i <= totalParts; i++ {
		url, err := client.GetMultipartUploadURL(context.Background(), uploadID, key, i, 30*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("获取分片上传URL失败: %v", err)
		}
		urls[i-1] = url
	}

	return urls, nil
}

// CompleteMultipartUpload 完成分片上传
func (s *FileService) CompleteMultipartUpload(uploadID, key, filename, md5 string, fileSize int64, uploaderID uint, parts []oss.Part) (*model.SysFile, error) {
	storage, err := storagesvc.Default.GetDefaultStorage()
	if err != nil {
		return nil, fmt.Errorf("获取存储配置失败: %v", err)
	}

	client, err := oss.GetClient(storage)
	if err != nil {
		return nil, fmt.Errorf("创建存储客户端失败: %v", err)
	}

	if err := client.CompleteMultipartUpload(context.Background(), key, uploadID, parts); err != nil {
		return nil, fmt.Errorf("完成分片上传失败: %v", err)
	}

	file := s.createFileRecord(filename, key, client.GetURL(key), md5, fileSize, uploaderID, storage)
	if err := s.CreateFile(file); err != nil {
		return nil, fmt.Errorf("保存文件记录失败: %v", err)
	}

	return file, nil
}

// GetUploadedParts 获取已上传的分片列表
func (s *FileService) GetUploadedParts(uploadID, key string) ([]oss.Part, error) {
	storage, err := storagesvc.Default.GetDefaultStorage()
	if err != nil {
		return nil, fmt.Errorf("获取存储配置失败: %v", err)
	}

	client, err := oss.GetClient(storage)
	if err != nil {
		return nil, fmt.Errorf("创建存储客户端失败: %v", err)
	}

	return client.ListParts(context.Background(), key, uploadID)
}

// AbortMultipartUpload 取消分片上传
func (s *FileService) AbortMultipartUpload(uploadID, key string) error {
	storage, err := storagesvc.Default.GetDefaultStorage()
	if err != nil {
		return fmt.Errorf("获取存储配置失败: %v", err)
	}

	client, err := oss.GetClient(storage)
	if err != nil {
		return fmt.Errorf("创建存储客户端失败: %v", err)
	}

	return client.AbortMultipartUpload(context.Background(), key, uploadID)
}

// SaveUploadedFile 保存已上传的文件记录
func (s *FileService) SaveUploadedFile(filename, key, url, md5 string, fileSize int64, uploaderID uint) (*model.SysFile, error) {
	storage, err := storagesvc.Default.GetDefaultStorage()
	if err != nil {
		return nil, fmt.Errorf("获取存储配置失败: %v", err)
	}

	file := s.createFileRecord(filename, key, url, md5, fileSize, uploaderID, storage)
	if err := s.CreateFile(file); err != nil {
		return nil, fmt.Errorf("保存文件记录失败: %v", err)
	}

	return file, nil
}

// getMimeType 根据扩展名获取MIME类型
func getMimeType(ext string) string {
	mimeTypes := map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"bmp":  "image/bmp",
		"webp": "image/webp",
		"svg":  "image/svg+xml",
		"ico":  "image/x-icon",
		"pdf":  "application/pdf",
		"doc":  "application/msword",
		"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"xls":  "application/vnd.ms-excel",
		"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"ppt":  "application/vnd.ms-powerpoint",
		"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"zip":  "application/zip",
		"rar":  "application/x-rar-compressed",
		"7z":   "application/x-7z-compressed",
		"tar":  "application/x-tar",
		"gz":   "application/gzip",
		"mp3":  "audio/mpeg",
		"wav":  "audio/wav",
		"mp4":  "video/mp4",
		"avi":  "video/x-msvideo",
		"mov":  "video/quicktime",
		"wmv":  "video/x-ms-wmv",
		"flv":  "video/x-flv",
		"txt":  "text/plain",
		"html": "text/html",
		"css":  "text/css",
		"js":   "application/javascript",
		"json": "application/json",
		"xml":  "application/xml",
	}

	if mime, ok := mimeTypes[strings.ToLower(ext)]; ok {
		return mime
	}
	return "application/octet-stream"
}
