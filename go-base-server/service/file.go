package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"go-base-server/global"
	"go-base-server/model"
	"go-base-server/service/oss"
)

type FileService struct{}

var File = new(FileService)

// GetFileList 获取文件列表
func (s *FileService) GetFileList(page, pageSize int, name, ext string) ([]model.SysFile, int64, error) {
	var files []model.SysFile
	var total int64

	db := global.DB.Model(&model.SysFile{}).
		Where("status = ?", 1).
		Preload("Storage")
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
	err := db.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&files).Error
	return files, total, err
}

// GetFileByID 根据ID获取文件
func (s *FileService) GetFileByID(id uint) (*model.SysFile, error) {
	var file model.SysFile
	err := global.DB.First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// GetFileByMD5 根据MD5获取文件（用于秒传）
func (s *FileService) GetFileByMD5(md5 string) (*model.SysFile, error) {
	var file model.SysFile
	err := global.DB.Where("md5 = ? AND status = ?", md5, 1).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// CreateFile 创建文件记录
func (s *FileService) CreateFile(file *model.SysFile) error {
	return global.DB.Create(file).Error
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(id uint) error {
	var file model.SysFile

	if err := global.DB.First(&file, id).Error; err != nil {
		return err
	}
	// DELETE FROM sys_user WHERE deleted_at IS  not  NULL
	// DELETE FROM sys_file WHERE status=0
	// 检查是否有用户正在使用该文件作为头像
	var userCount int64
	global.DB.Model(&model.SysUser{}).Where("avatar_file_id = ?", id).Count(&userCount)
	if userCount > 0 {

		return errors.New("文件正在被使用，无法删除")
	}

	// 获取存储配置
	storage, err := Storage.GetStorageByID(file.StorageID)
	if err != nil {
		return err
	}

	// 获取客户端并删除文件
	client, err := oss.GetClient(storage)
	if err == nil {
		_ = client.Delete(context.Background(), file.Path)
	}

	// 软删除数据库记录
	return global.DB.Model(&model.SysFile{}).Where("id = ?", id).Update("status", 0).Error
}

// GenerateFilePath 生成文件存储路径
// 注意：这里返回的是相对于存储根目录的 key，不再带 "uploads/" 前缀，
// 方便通过 base_path = "uploads" + Gin 静态路由 /uploads 直接访问。
func (s *FileService) GenerateFilePath(filename string) string {
	ext := filepath.Ext(filename)
	now := time.Now()
	// 按日期组织目录，例如：2026/01/25/xxx.png
	return fmt.Sprintf("%d/%02d/%02d/%d%s", now.Year(), now.Month(), now.Day(), now.UnixNano(), ext)
}

// GetUploadCredential 获取上传凭证
func (s *FileService) GetUploadCredential(filename string, storageID uint) (*oss.UploadCredential, error) {
	var storage *model.SysStorage
	var err error

	if storageID > 0 {
		storage, err = Storage.GetStorageByID(storageID)
	} else {
		storage, err = Storage.GetDefaultStorage()
	}
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
func (s *FileService) InitMultipartUpload(filename, md5 string, fileSize int64, storageID uint) (*oss.MultipartUpload, *model.SysStorage, error) {
	var storage *model.SysStorage
	var err error

	if storageID > 0 {
		storage, err = Storage.GetStorageByID(storageID)
	} else {
		storage, err = Storage.GetDefaultStorage()
	}
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
func (s *FileService) GetMultipartUploadURLs(uploadID, key string, totalParts int, storageID uint) ([]string, error) {
	storage, err := Storage.GetStorageByID(storageID)
	if err != nil {
		return nil, fmt.Errorf("获取存储配置失败: %v", err)
	}

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

		// 本地存储和 MinIO 的分片上传走后端代理，需要在 URL 上携带 storage_id，
		// 远程 OSS/COS 使用预签名直传，URL 不能再追加参数，否则会导致签名失效。
		if storage.Type == model.StorageTypeLocal || storage.Type == model.StorageTypeMinio {
			if strings.Contains(url, "?") {
				url = fmt.Sprintf("%s&storage_id=%d", url, storageID)
			} else {
				url = fmt.Sprintf("%s?storage_id=%d", url, storageID)
			}
		}

		urls[i-1] = url
	}

	return urls, nil
}

// CompleteMultipartUpload 完成分片上传
func (s *FileService) CompleteMultipartUpload(uploadID, key, filename, md5 string, fileSize int64, storageID, uploaderID uint, parts []oss.Part) (*model.SysFile, error) {
	storage, err := Storage.GetStorageByID(storageID)
	if err != nil {
		return nil, fmt.Errorf("获取存储配置失败: %v", err)
	}

	client, err := oss.GetClient(storage)
	if err != nil {
		return nil, fmt.Errorf("创建存储客户端失败: %v", err)
	}

	// 完成分片上传
	if err := client.CompleteMultipartUpload(context.Background(), key, uploadID, parts); err != nil {
		return nil, fmt.Errorf("完成分片上传失败: %v", err)
	}

	// 创建文件记录
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	file := &model.SysFile{
		Name:       filename,
		Path:       key,
		URL:        client.GetURL(key),
		Size:       fileSize,
		Ext:        ext,
		MimeType:   getMimeType(ext),
		MD5:        md5,
		StorageID:  storageID,
		UploaderID: uploaderID,
		Status:     1,
	}

	if err := s.CreateFile(file); err != nil {
		return nil, fmt.Errorf("保存文件记录失败: %v", err)
	}

	return file, nil
}

// GetUploadedParts 获取已上传的分片列表
func (s *FileService) GetUploadedParts(uploadID, key string, storageID uint) ([]oss.Part, error) {
	storage, err := Storage.GetStorageByID(storageID)
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
func (s *FileService) AbortMultipartUpload(uploadID, key string, storageID uint) error {
	storage, err := Storage.GetStorageByID(storageID)
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
func (s *FileService) SaveUploadedFile(filename, key, url, md5 string, fileSize int64, storageID, uploaderID uint) (*model.SysFile, error) {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	file := &model.SysFile{
		Name:       filename,
		Path:       key,
		URL:        url,
		Size:       fileSize,
		Ext:        ext,
		MimeType:   getMimeType(ext),
		MD5:        md5,
		StorageID:  storageID,
		UploaderID: uploaderID,
		Status:     1,
	}

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
