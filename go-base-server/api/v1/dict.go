package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"go-base-server/model/request"
	"go-base-server/model/response"
	"go-base-server/service"
)

type DictApi struct{}

var Dict = new(DictApi)

// ==================== 字典类型 ====================

// GetDictTypeList 获取字典类型列表
func (a *DictApi) GetDictTypeList(c *gin.Context) {
	var req request.DictTypeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	list, total, err := service.Dict.GetDictTypeList(&req)
	if err != nil {
		response.Fail(c, "获取字典类型列表失败")
		return
	}

	response.OkWithData(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// GetAllDictTypes 获取所有字典类型（不分页）
func (a *DictApi) GetAllDictTypes(c *gin.Context) {
	list, err := service.Dict.GetAllDictTypes()
	if err != nil {
		response.Fail(c, "获取字典类型列表失败")
		return
	}
	response.OkWithData(c, list)
}

// GetDictType 获取字典类型详情
func (a *DictApi) GetDictType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	dictType, err := service.Dict.GetDictType(uint(id))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithData(c, dictType)
}

// CreateDictType 创建字典类型
func (a *DictApi) CreateDictType(c *gin.Context) {
	var req request.CreateDictTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Dict.CreateDictType(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "创建成功")
}

// UpdateDictType 更新字典类型
func (a *DictApi) UpdateDictType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.UpdateDictTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Dict.UpdateDictType(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "更新成功")
}

// DeleteDictType 删除字典类型
func (a *DictApi) DeleteDictType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Dict.DeleteDictType(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "删除成功")
}

// ==================== 字典数据 ====================

// GetDictDataList 获取字典数据列表
func (a *DictApi) GetDictDataList(c *gin.Context) {
	var req request.DictDataListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: " + err.Error())
		return
	}

	list, total, err := service.Dict.GetDictDataList(&req)
	if err != nil {
		response.Fail(c, "获取字典数据列表失败")
		return
	}

	response.OkWithData(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// GetDictDataByType 根据字典类型获取字典数据
func (a *DictApi) GetDictDataByType(c *gin.Context) {
	dictType := c.Param("type")
	if dictType == "" {
		response.BadRequest(c, "字典类型不能为空")
		return
	}

	list, err := service.Dict.GetDictDataByType(dictType)
	if err != nil {
		response.Fail(c, "获取字典数据失败")
		return
	}

	response.OkWithData(c, list)
}

// GetDictData 获取字典数据详情
func (a *DictApi) GetDictData(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	dictData, err := service.Dict.GetDictData(uint(id))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithData(c, dictData)
}

// CreateDictData 创建字典数据
func (a *DictApi) CreateDictData(c *gin.Context) {
	var req request.CreateDictDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Dict.CreateDictData(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "创建成功")
}

// UpdateDictData 更新字典数据
func (a *DictApi) UpdateDictData(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req request.UpdateDictDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Dict.UpdateDictData(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "更新成功")
}

// DeleteDictData 删除字典数据
func (a *DictApi) DeleteDictData(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := service.Dict.DeleteDictData(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.OkWithMessage(c, "删除成功")
}
