package v1

import (
	"strconv"

	"go-base-server/model"
	"go-base-server/model/response"
	"go-base-server/service"

	"github.com/gin-gonic/gin"
)

type StorageApi struct{}

var Storage = new(StorageApi)

// GetStorageList 获取存储配置列表
func (a *StorageApi) GetStorageList(c *gin.Context) {
	storages, err := service.Storage.GetStorageList()
	if err != nil {
		response.Fail(c, "获取存储配置列表失败")
		return
	}
	response.OkWithData(c, storages)
}

// GetStorage 获取单个存储配置
func (a *StorageApi) GetStorage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	storage, err := service.Storage.GetStorageByID(uint(id))
	if err != nil {
		response.NotFound(c, "存储配置不存在")
		return
	}
	response.OkWithData(c, storage)
}

// CreateStorage 创建存储配置
func (a *StorageApi) CreateStorage(c *gin.Context) {
	var storage model.SysStorage
	if err := c.ShouldBindJSON(&storage); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 测试配置是否可用
	if err := service.Storage.TestStorage(&storage); err != nil {
		response.Fail(c, "存储配置测试失败: "+err.Error())
		return
	}

	if err := service.Storage.CreateStorage(&storage); err != nil {
		response.Fail(c, "创建存储配置失败")
		return
	}
	response.OkWithData(c, storage)
}

// UpdateStorage 更新存储配置
func (a *StorageApi) UpdateStorage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 如果更新了配置，测试是否可用
	if config, ok := data["config"]; ok {
		storage, _ := service.Storage.GetStorageByID(uint(id))
		if storage != nil {
			storage.Config = config.(string)
			if err := service.Storage.TestStorage(storage); err != nil {
				response.Fail(c, "存储配置测试失败: "+err.Error())
				return
			}
		}
	}

	if err := service.Storage.UpdateStorage(uint(id), data); err != nil {
		response.Fail(c, "更新存储配置失败")
		return
	}
	response.OkWithMessage(c, "更新成功")
}

// DeleteStorage 删除存储配置
func (a *StorageApi) DeleteStorage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := service.Storage.DeleteStorage(uint(id)); err != nil {
		response.Fail(c, "删除存储配置失败")
		return
	}
	response.OkWithMessage(c, "删除成功")
}

// SetDefaultStorage 设置默认存储
func (a *StorageApi) SetDefaultStorage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := service.Storage.SetDefaultStorage(uint(id)); err != nil {
		response.Fail(c, "设置默认存储失败")
		return
	}
	response.OkWithMessage(c, "设置成功")
}

// TestStorage 测试存储配置
func (a *StorageApi) TestStorage(c *gin.Context) {
	var storage model.SysStorage
	if err := c.ShouldBindJSON(&storage); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Storage.TestStorage(&storage); err != nil {
		response.Fail(c, "测试失败: "+err.Error())
		return
	}
	response.OkWithMessage(c, "测试成功")
}
