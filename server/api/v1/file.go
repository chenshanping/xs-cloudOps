package v1

import (
	"bytes"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"server/model/response"
	"server/service"
	"server/service/oss"
)

type FileApi struct{}

var File = new(FileApi)

// GetFileList 获取文件列表
func (a *FileApi) GetFileList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	name := c.Query("name")
	ext := c.Query("ext")
	storageID, _ := strconv.ParseUint(c.Query("storage_id"), 10, 32)

	files, total, err := service.File.GetFileList(page, pageSize, name, ext, uint(storageID))
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
		response.Fail(c, "删除文件失败")
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
	} else if successCount > 0 {
		response.OkWithData(c, gin.H{
			"success_count": successCount,
			"failed_count":  len(failedMsgs),
			"failed_msgs":   failedMsgs,
		})
	} else {
		response.Fail(c, "删除失败")
	}
}

// GetUploadCredential 获取上传凭证
func (a *FileApi) GetUploadCredential(c *gin.Context) {
	var req struct {
		Filename  string `json:"filename" binding:"required"`
		StorageID uint   `json:"storage_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	credential, err := service.File.GetUploadCredential(req.Filename, req.StorageID)
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
		response.OkWithData(c, gin.H{
			"exists": true,
			"file":   file,
		})
	} else {
		response.OkWithData(c, gin.H{
			"exists": false,
		})
	}
}

// InitMultipartUpload 初始化分片上传
func (a *FileApi) InitMultipartUpload(c *gin.Context) {
	var req struct {
		Filename  string `json:"filename" binding:"required"`
		FileSize  int64  `json:"file_size" binding:"required"`
		MD5       string `json:"md5" binding:"required"`
		StorageID uint   `json:"storage_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 先检查是否可以秒传
	if file, exists := service.File.CheckFileMD5(req.MD5); exists {
		response.OkWithData(c, gin.H{
			"instant_upload": true,
			"file":           file,
		})
		return
	}

	upload, storage, err := service.File.InitMultipartUpload(req.Filename, req.MD5, req.FileSize, req.StorageID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	// 计算分片数量
	chunkSize := upload.ChunkSize
	totalParts := int(req.FileSize / chunkSize)
	if req.FileSize%chunkSize > 0 {
		totalParts++
	}

	// 获取所有分片的上传URL
	urls, err := service.File.GetMultipartUploadURLs(upload.UploadID, upload.Key, totalParts, storage.ID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithData(c, gin.H{
		"instant_upload": false,
		"upload_id":      upload.UploadID,
		"key":            upload.Key,
		"storage_id":     storage.ID,
		"chunk_size":     chunkSize,
		"total_parts":    totalParts,
		"upload_urls":    urls,
	})
}

// GetUploadedParts 获取已上传的分片列表
func (a *FileApi) GetUploadedParts(c *gin.Context) {
	uploadID := c.Query("upload_id")
	key := c.Query("key")
	storageID, _ := strconv.ParseUint(c.Query("storage_id"), 10, 32)

	parts, err := service.File.GetUploadedParts(uploadID, key, uint(storageID))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, parts)
}

// CompleteMultipartUpload 完成分片上传
func (a *FileApi) CompleteMultipartUpload(c *gin.Context) {
	var req struct {
		UploadID  string     `json:"upload_id" binding:"required"`
		Key       string     `json:"key" binding:"required"`
		Filename  string     `json:"filename" binding:"required"`
		FileSize  int64      `json:"file_size" binding:"required"`
		MD5       string     `json:"md5" binding:"required"`
		StorageID uint       `json:"storage_id" binding:"required"`
		Parts     []oss.Part `json:"parts" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("user_id")
	uploaderID := userID.(uint)

	file, err := service.File.CompleteMultipartUpload(
		req.UploadID, req.Key, req.Filename, req.MD5, req.FileSize,
		req.StorageID, uploaderID, req.Parts,
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
		UploadID  string `json:"upload_id" binding:"required"`
		Key       string `json:"key" binding:"required"`
		StorageID uint   `json:"storage_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.File.AbortMultipartUpload(req.UploadID, req.Key, req.StorageID); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "取消成功")
}

// SaveUploadedFile 保存已上传的文件记录
func (a *FileApi) SaveUploadedFile(c *gin.Context) {
	var req struct {
		Filename  string `json:"filename" binding:"required"`
		Key       string `json:"key" binding:"required"`
		URL       string `json:"url" binding:"required"`
		FileSize  int64  `json:"file_size" binding:"required"`
		MD5       string `json:"md5" binding:"required"`
		StorageID uint   `json:"storage_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("user_id")
	uploaderID := userID.(uint)

	file, err := service.File.SaveUploadedFile(
		req.Filename, req.Key, req.URL, req.MD5, req.FileSize,
		req.StorageID, uploaderID,
	)
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

	// 获取默认存储
	storage, err := service.Storage.GetDefaultStorage()
	if err != nil {
		response.Fail(c, "获取存储配置失败")
		return
	}

	// 获取客户端
	client, err := oss.GetClient(storage)
	if err != nil {
		response.Fail(c, "创建存储客户端失败")
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		response.Fail(c, "打开文件失败")
		return
	}
	defer src.Close()

	// 上传文件
	if err := client.Upload(c.Request.Context(), key, src, file.Size); err != nil {
		response.Fail(c, "上传文件失败: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("user_id")
	uploaderID := userID.(uint)

	// 保存文件记录
	md5 := c.PostForm("md5")
	savedFile, err := service.File.SaveUploadedFile(
		file.Filename, key, client.GetURL(key), md5, file.Size,
		storage.ID, uploaderID,
	)
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
	storageID, _ := strconv.ParseUint(c.Query("storage_id"), 10, 32)

	// 获取文件内容
	file, err := c.FormFile("file")
	if err != nil {
		// 尝试从body读取
		body := c.Request.Body
		defer body.Close()

		storage, err := service.Storage.GetStorageByID(uint(storageID))
		if err != nil {
			response.Fail(c, "获取存储配置失败")
			return
		}

		client, err := oss.GetClient(storage)
		if err != nil {
			response.Fail(c, "创建存储客户端失败")
			return
		}

		// 本地存储
		if localClient, ok := client.(*oss.LocalClient); ok {
			if err := localClient.UploadChunk(uploadID, partNumber, body); err != nil {
				response.Fail(c, "上传分片失败: "+err.Error())
				return
			}
			response.OkWithData(c, gin.H{"part_number": partNumber})
			return
		}

		// MinIO
		if minioClient, ok := client.(*oss.MinioClient); ok {
			data, _ := io.ReadAll(body)
			etag, err := minioClient.UploadPart(c.Request.Context(), key, uploadID, partNumber,
				io.NopCloser(bytes.NewReader(data)), int64(len(data)))
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

	// 使用表单文件上传
	src, err := file.Open()
	if err != nil {
		response.Fail(c, "打开文件失败")
		return
	}
	defer src.Close()

	storage, err := service.Storage.GetStorageByID(uint(storageID))
	if err != nil {
		response.Fail(c, "获取存储配置失败")
		return
	}

	client, err := oss.GetClient(storage)
	if err != nil {
		response.Fail(c, "创建存储客户端失败")
		return
	}

	// 本地存储
	if localClient, ok := client.(*oss.LocalClient); ok {
		if err := localClient.UploadChunk(uploadID, partNumber, src); err != nil {
			response.Fail(c, "上传分片失败: "+err.Error())
			return
		}
		response.OkWithData(c, gin.H{"part_number": partNumber})
		return
	}

	// MinIO
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
