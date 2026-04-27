package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"server/global"
	"server/model"
	"server/service/oss"
)

type FileService struct{}

var File = new(FileService)

// GetFileList 获取文件列表
func (s *FileService) GetFileList(page, pageSize int, name, ext string) ([]model.SysFile, int64, error) {
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

	db.Count(&total)
	err := db.Order("id desc").Offset((page-1)*pageSize).Limit(pageSize).Find(&files).Error
	return files, total, err
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
		return Storage.GetStorageByType(model.StorageType(file.StorageType))
	}
	return Storage.GetDefaultStorage()
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
		UploaderID:    uploaderID,
		Status:        1,
	}
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(id uint) error {
	var file model.SysFile
	if err := global.DB.First(&file, id).Error; err != nil {
		return err
	}

	var userCount int64
	global.DB.Model(&model.SysUser{}).Where("avatar_file_id = ?", id).Count(&userCount)
	if userCount > 0 {
		return errors.New("文件正在被使用，无法删除")
	}

	storage, err := s.resolveFileStorage(file)
	if err != nil {
		return err
	}

	client, err := oss.GetClient(storage)
	if err == nil {
		_ = client.Delete(context.Background(), file.Path)
	}

	return global.DB.Model(&model.SysFile{}).Where("id = ?", id).Update("status", 0).Error
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

// GenerateFilePath 生成文件存储路径
func (s *FileService) GenerateFilePath(filename string) string {
	ext := filepath.Ext(filename)
	now := time.Now()
	return fmt.Sprintf("%d/%02d/%02d/%d%s", now.Year(), now.Month(), now.Day(), now.UnixNano(), ext)
}

// GetUploadCredential 获取上传凭证
func (s *FileService) GetUploadCredential(filename string) (*oss.UploadCredential, error) {
	storage, err := Storage.GetDefaultStorage()
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
	storage, err := Storage.GetDefaultStorage()
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
	storage, err := Storage.GetDefaultStorage()
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
	storage, err := Storage.GetDefaultStorage()
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
	storage, err := Storage.GetDefaultStorage()
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
	storage, err := Storage.GetDefaultStorage()
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
