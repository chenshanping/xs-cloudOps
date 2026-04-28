package v1

import (
	"bytes"
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/service/oss"
)

type FileApi struct{}

var File = new(FileApi)

func fileDeleteFailureMessage(err error) string {
	if err == nil {
		return "删除失败"
	}
	message := strings.TrimSpace(err.Error())
	if message == "" {
		return "删除失败"
	}
	return message
}

func batchFileDeleteFailureMessage(failedMsgs []string) string {
	if len(failedMsgs) == 0 {
		return "删除失败"
	}

	normalized := make([]string, 0, len(failedMsgs))
	for _, failedMsg := range failedMsgs {
		message := strings.TrimSpace(failedMsg)
		if message != "" {
			normalized = append(normalized, message)
		}
	}

	if len(normalized) == 0 {
		return "删除失败"
	}
	if len(normalized) == 1 {
		return normalized[0]
	}
	return strings.Join(normalized, "；")
}

// GetFileList 获取文件列表
func (a *FileApi) GetFileList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	name := c.Query("name")
	ext := c.Query("ext")

	files, total, err := service.File.GetFileList(page, pageSize, name, ext)
	if err != nil {
		response.Fail(c, "获取文件列表失败")
		return
	}
	response.OkWithPage(c, files, total, page, pageSize)
}

// GetFile 获取文件详情
func (a *FileApi) GetFile(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	file, err := service.File.GetFileByID(uint(id))
	if err != nil {
		response.NotFound(c, "文件不存在")
		return
	}
	response.OkWithData(c, file)
}

// DeleteFile 删除文件
func (a *FileApi) DeleteFile(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := service.File.DeleteFile(uint(id)); err != nil {
		global.Log.Warnf("删除文件失败: file_id=%d, path=%s, err=%v", id, c.FullPath(), err)
		response.Fail(c, fileDeleteFailureMessage(err))
		return
	}
	response.OkWithMessage(c, "删除成功")
}

// BatchDeleteFiles 批量删除文件
func (a *FileApi) BatchDeleteFiles(c *gin.Context) {
	var req struct {
		Ids []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if len(req.Ids) == 0 {
		response.BadRequest(c, "请选择要删除的文件")
		return
	}

	successCount, failedMsgs := service.File.BatchDeleteFiles(req.Ids)
	if len(failedMsgs) == 0 {
		response.OkWithMessage(c, "batch_delete_success")
		return
	}
	if successCount > 0 {
		global.Log.Warnf("批量删除文件部分失败: ids=%v, success_count=%d, failed=%v", req.Ids, successCount, failedMsgs)
		response.OkWithData(c, gin.H{
			"success_count": successCount,
			"failed_count":  len(failedMsgs),
			"failed_msgs":   failedMsgs,
		})
		return
	}
	global.Log.Warnf("批量删除文件全部失败: ids=%v, failed=%v", req.Ids, failedMsgs)
	response.Fail(c, batchFileDeleteFailureMessage(failedMsgs))
}

// PreviewFileMigration 预览文件迁移
func (a *FileApi) PreviewFileMigration(c *gin.Context) {
	var req request.FileMigrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		if strings.TrimSpace(req.Scope) == "selected" {
			response.BadRequest(c, "请选择要迁移的文件")
			return
		}
	}
	if strings.TrimSpace(req.SourceStorageType) == "" {
		response.BadRequest(c, "请选择源存储")
		return
	}
	if strings.TrimSpace(req.TargetStorageType) == "" {
		response.BadRequest(c, "请选择目标存储")
		return
	}

	result, err := service.File.PreviewFileMigration(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, result)
}

// ExecuteFileMigration 执行文件迁移
func (a *FileApi) ExecuteFileMigration(c *gin.Context) {
	var req request.FileMigrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		if strings.TrimSpace(req.Scope) == "selected" {
			response.BadRequest(c, "请选择要迁移的文件")
			return
		}
	}
	if strings.TrimSpace(req.SourceStorageType) == "" {
		response.BadRequest(c, "请选择源存储")
		return
	}
	if strings.TrimSpace(req.TargetStorageType) == "" {
		response.BadRequest(c, "请选择目标存储")
		return
	}

	result, err := service.File.StartFileMigrationTask(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, result)
}

// GetCurrentFileMigrationTask 获取当前文件迁移任务状态
func (a *FileApi) GetCurrentFileMigrationTask(c *gin.Context) {
	response.OkWithData(c, service.File.GetCurrentFileMigrationTask())
}

// GetUploadCredential 获取上传凭证
func (a *FileApi) GetUploadCredential(c *gin.Context) {
	var req struct {
		Filename string `json:"filename" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	credential, err := service.File.GetUploadCredential(req.Filename)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, credential)
}

// CheckFileMD5 检查文件MD5（秒传）
func (a *FileApi) CheckFileMD5(c *gin.Context) {
	var req struct {
		MD5 string `json:"md5" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	file, exists := service.File.CheckFileMD5(req.MD5)
	if exists {
		response.OkWithData(c, gin.H{"exists": true, "file": file})
		return
	}
	response.OkWithData(c, gin.H{"exists": false})
}

// InitMultipartUpload 初始化分片上传
func (a *FileApi) InitMultipartUpload(c *gin.Context) {
	var req struct {
		Filename string `json:"filename" binding:"required"`
		FileSize int64  `json:"file_size" binding:"required"`
		MD5      string `json:"md5" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if file, exists := service.File.CheckFileMD5(req.MD5); exists {
		response.OkWithData(c, gin.H{
			"instant_upload": true,
			"file":           file,
		})
		return
	}

	upload, storage, err := service.File.InitMultipartUpload(req.Filename, req.MD5, req.FileSize)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	chunkSize := upload.ChunkSize
	totalParts := int(req.FileSize / chunkSize)
	if req.FileSize%chunkSize > 0 {
		totalParts++
	}

	urls, err := service.File.GetMultipartUploadURLs(upload.UploadID, upload.Key, totalParts, storage)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithData(c, gin.H{
		"instant_upload": false,
		"upload_id":      upload.UploadID,
		"key":            upload.Key,
		"chunk_size":     chunkSize,
		"total_parts":    totalParts,
		"upload_urls":    urls,
	})
}

// GetUploadedParts 获取已上传的分片列表
func (a *FileApi) GetUploadedParts(c *gin.Context) {
	uploadID := c.Query("upload_id")
	key := c.Query("key")

	parts, err := service.File.GetUploadedParts(uploadID, key)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, parts)
}

// CompleteMultipartUpload 完成分片上传
func (a *FileApi) CompleteMultipartUpload(c *gin.Context) {
	var req struct {
		UploadID string     `json:"upload_id" binding:"required"`
		Key      string     `json:"key" binding:"required"`
		Filename string     `json:"filename" binding:"required"`
		FileSize int64      `json:"file_size" binding:"required"`
		MD5      string     `json:"md5" binding:"required"`
		Parts    []oss.Part `json:"parts" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	uploaderID := userID.(uint)

	file, err := service.File.CompleteMultipartUpload(
		req.UploadID, req.Key, req.Filename, req.MD5, req.FileSize, uploaderID, req.Parts,
	)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, file)
}

// AbortMultipartUpload 取消分片上传
func (a *FileApi) AbortMultipartUpload(c *gin.Context) {
	var req struct {
		UploadID string `json:"upload_id" binding:"required"`
		Key      string `json:"key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.File.AbortMultipartUpload(req.UploadID, req.Key); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "取消成功")
}

// SaveUploadedFile 保存已上传的文件记录
func (a *FileApi) SaveUploadedFile(c *gin.Context) {
	var req struct {
		Filename string `json:"filename" binding:"required"`
		Key      string `json:"key" binding:"required"`
		URL      string `json:"url" binding:"required"`
		FileSize int64  `json:"file_size" binding:"required"`
		MD5      string `json:"md5" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	uploaderID := userID.(uint)

	file, err := service.File.SaveUploadedFile(req.Filename, req.Key, req.URL, req.MD5, req.FileSize, uploaderID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, file)
}

// UploadLocalFile 本地存储文件上传（代理上传）
func (a *FileApi) UploadLocalFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择文件")
		return
	}

	key := c.PostForm("key")
	if key == "" {
		key = service.File.GenerateFilePath(file.Filename)
	}

	storage, err := service.Storage.GetDefaultStorage()
	if err != nil {
		response.Fail(c, "获取存储配置失败")
		return
	}

	client, err := oss.GetClient(storage)
	if err != nil {
		response.Fail(c, "创建存储客户端失败")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Fail(c, "打开文件失败")
		return
	}
	defer src.Close()

	if err := client.Upload(c.Request.Context(), key, src, file.Size); err != nil {
		response.Fail(c, "上传文件失败: "+err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	uploaderID := userID.(uint)

	md5 := c.PostForm("md5")
	savedFile, err := service.File.SaveUploadedFile(file.Filename, key, client.GetURL(key), md5, file.Size, uploaderID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithData(c, savedFile)
}

// UploadChunk 上传分片（本地存储和MinIO代理）
func (a *FileApi) UploadChunk(c *gin.Context) {
	uploadID := c.Query("upload_id")
	key := c.Query("key")
	partNumber, _ := strconv.Atoi(c.Query("part_number"))

	storage, err := service.Storage.GetDefaultStorage()
	if err != nil {
		response.Fail(c, "获取存储配置失败")
		return
	}

	client, err := oss.GetClient(storage)
	if err != nil {
		response.Fail(c, "创建存储客户端失败")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		body := c.Request.Body
		defer body.Close()

		if localClient, ok := client.(*oss.LocalClient); ok {
			if err := localClient.UploadChunk(uploadID, partNumber, body); err != nil {
				response.Fail(c, "上传分片失败: "+err.Error())
				return
			}
			response.OkWithData(c, gin.H{"part_number": partNumber})
			return
		}

		if minioClient, ok := client.(*oss.MinioClient); ok {
			data, _ := io.ReadAll(body)
			etag, err := minioClient.UploadPart(c.Request.Context(), key, uploadID, partNumber, io.NopCloser(bytes.NewReader(data)), int64(len(data)))
			if err != nil {
				response.Fail(c, "上传分片失败: "+err.Error())
				return
			}
			response.OkWithData(c, gin.H{"part_number": partNumber, "etag": etag})
			return
		}

		response.Fail(c, "不支持的存储类型")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Fail(c, "打开文件失败")
		return
	}
	defer src.Close()

	if localClient, ok := client.(*oss.LocalClient); ok {
		if err := localClient.UploadChunk(uploadID, partNumber, src); err != nil {
			response.Fail(c, "上传分片失败: "+err.Error())
			return
		}
		response.OkWithData(c, gin.H{"part_number": partNumber})
		return
	}

	if minioClient, ok := client.(*oss.MinioClient); ok {
		etag, err := minioClient.UploadPart(c.Request.Context(), key, uploadID, partNumber, src, file.Size)
		if err != nil {
			response.Fail(c, "上传分片失败: "+err.Error())
			return
		}
		response.OkWithData(c, gin.H{"part_number": partNumber, "etag": etag})
		return
	}

	response.Fail(c, "不支持的存储类型")
}
