package modules

import (
	"github.com/gin-gonic/gin"

	v1 "server/api/v1"
	"server/router/registry"
)

func init() {
	RegisterModule(&FileModule{})
}

type FileModule struct{}

func (m *FileModule) Name() string {
	return "文件管理"
}

func (m *FileModule) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 无公开路由
}

func (m *FileModule) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	// 文件列表和详情
	R(rg, "GET", "/files", m.Name(), "文件列表", v1.File.GetFileList, registry.WithAuth())
	R(rg, "GET", "/files/:id", m.Name(), "文件详情", v1.File.GetFile, registry.WithAuth())
	R(rg, "DELETE", "/files/:id", m.Name(), "删除文件", v1.File.DeleteFile, registry.WithAuth())
	R(rg, "DELETE", "/files/batch", m.Name(), "批量删除文件", v1.File.BatchDeleteFiles, registry.WithAuth())
	R(rg, "POST", "/files/migrate/preview", m.Name(), "预览文件迁移", v1.File.PreviewFileMigration, registry.WithAuth())
	R(rg, "POST", "/files/migrate/execute", m.Name(), "执行文件迁移", v1.File.ExecuteFileMigration, registry.WithAuth())
	R(rg, "GET", "/files/migrate/task/current", m.Name(), "获取当前文件迁移任务", v1.File.GetCurrentFileMigrationTask, registry.WithAuth())

	// 上传相关
	R(rg, "POST", "/files/credential", m.Name(), "获取上传凭证", v1.File.GetUploadCredential, registry.WithAuth())
	R(rg, "POST", "/files/check-md5", m.Name(), "MD5秒传检查", v1.File.CheckFileMD5, registry.WithAuth())
	R(rg, "POST", "/files/save", m.Name(), "保存上传文件", v1.File.SaveUploadedFile, registry.WithAuth())

	// 分片上传
	R(rg, "POST", "/files/multipart/init", m.Name(), "初始化分片上传", v1.File.InitMultipartUpload, registry.WithAuth())
	R(rg, "GET", "/files/multipart/parts", m.Name(), "获取已上传分片", v1.File.GetUploadedParts, registry.WithAuth())
	R(rg, "POST", "/files/multipart/complete", m.Name(), "完成分片上传", v1.File.CompleteMultipartUpload, registry.WithAuth())
	R(rg, "POST", "/files/multipart/abort", m.Name(), "取消分片上传", v1.File.AbortMultipartUpload, registry.WithAuth())

	// 本地上传（代理）
	R(rg, "POST", "/files/upload/local", m.Name(), "本地文件上传", v1.File.UploadLocalFile, registry.WithAuth())
	R(rg, "POST", "/files/upload/chunk", m.Name(), "上传分片", v1.File.UploadChunk, registry.WithAuth())
}
