package v1

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"

	"server/model/request"
	"server/model/response"
	"server/service"
)

type CmdbApi struct{}

var Cmdb = new(CmdbApi)

func (a *CmdbApi) GetHostGroups(c *gin.Context) {
	var req request.CmdbHostGroupListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	list, err := service.CmdbHost.ListGroups(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, list)
}

func (a *CmdbApi) CreateHostGroup(c *gin.Context) {
	var req request.CreateCmdbHostGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	if err := service.CmdbHost.CreateGroup(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "创建成功")
}

func (a *CmdbApi) UpdateHostGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	var req request.UpdateCmdbHostGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	if err := service.CmdbHost.UpdateGroup(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "更新成功")
}

func (a *CmdbApi) DeleteHostGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := service.CmdbHost.DeleteGroup(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "删除成功")
}

func (a *CmdbApi) GetHostTags(c *gin.Context) {
	var req request.CmdbHostTagListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	list, err := service.CmdbHost.ListTags(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, list)
}

func (a *CmdbApi) CreateHostTag(c *gin.Context) {
	var req request.CreateCmdbHostTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	if err := service.CmdbHost.CreateTag(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "创建成功")
}

func (a *CmdbApi) UpdateHostTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	var req request.UpdateCmdbHostTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	if err := service.CmdbHost.UpdateTag(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "更新成功")
}

func (a *CmdbApi) DeleteHostTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := service.CmdbHost.DeleteTag(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "删除成功")
}

func (a *CmdbApi) GetCredentials(c *gin.Context) {
	var req request.CmdbCredentialListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	list, total, err := service.CmdbCredential.List(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithPage(c, list, total, req.Page, req.PageSize)
}

func (a *CmdbApi) GetCredentialOptions(c *gin.Context) {
	list, err := service.CmdbCredential.ListOptions()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, list)
}

func (a *CmdbApi) GetCredential(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	item, err := service.CmdbCredential.Get(uint(id))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, item)
}

func (a *CmdbApi) CreateCredential(c *gin.Context) {
	var req request.CreateCmdbCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	if err := service.CmdbCredential.Create(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "创建成功")
}

func (a *CmdbApi) UpdateCredential(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	var req request.UpdateCmdbCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	if err := service.CmdbCredential.Update(uint(id), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "更新成功")
}

func (a *CmdbApi) DeleteCredential(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := service.CmdbCredential.Delete(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "删除成功")
}

func (a *CmdbApi) GetHosts(c *gin.Context) {
	var req request.CmdbHostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	list, total, err := service.CmdbHost.ListHosts(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithPage(c, list, total, req.Page, req.PageSize)
}

func (a *CmdbApi) GetHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	item, err := service.CmdbHost.GetHost(uint(id))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, item)
}

func (a *CmdbApi) CreateHost(c *gin.Context) {
	var req request.CreateCmdbHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	item, err := service.CmdbHost.CreateHost(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, item)
}

func (a *CmdbApi) UpdateHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	var req request.UpdateCmdbHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.BindErrorMessage(err, req))
		return
	}
	item, err := service.CmdbHost.UpdateHost(uint(id), &req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, item)
}

func (a *CmdbApi) DeleteHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := service.CmdbHost.DeleteHost(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithMessage(c, "删除成功")
}

func (a *CmdbApi) VerifyHost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := service.CmdbHost.VerifyHost(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}
	item, err := service.CmdbHost.GetHost(uint(id))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, item)
}

func (a *CmdbApi) ImportHosts(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "导入文件不能为空")
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		response.Fail(c, "打开导入文件失败")
		return
	}
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		response.Fail(c, "读取导入文件失败")
		return
	}
	result, err := service.CmdbHost.ImportHosts(fileData)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OkWithData(c, result)
}

func (a *CmdbApi) GetHostImportTemplate(c *gin.Context) {
	buf, filename, err := service.CmdbHost.GetImportTemplate()
	if err != nil {
		response.Fail(c, "生成导入模板失败")
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(200, "application/octet-stream", buf)
}
