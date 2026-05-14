package apisvc

import (
	"errors"

	"server/global"
	"server/model"
	"server/model/request"
)

type ApiService struct{}

var Default = &ApiService{}

// 获取API列表
func (s *ApiService) GetApiList(req *request.ApiListRequest) ([]model.SysApi, int64, error) {
	var apis []model.SysApi
	var total int64

	db := global.DB.Model(&model.SysApi{})

	if req.Path != "" {
		db = db.Where("path LIKE ?", "%"+req.Path+"%")
	}
	if req.Method != "" {
		db = db.Where("method = ?", req.Method)
	}
	if req.Group != "" {
		db = db.Where("`group` = ?", req.Group)
	}
	if req.NeedAuth != nil {
		db = db.Where("need_auth = ?", *req.NeedAuth)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := req.GetOffset()
	if err := db.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&apis).Error; err != nil {
		return nil, 0, err
	}

	return apis, total, nil
}

// 获取全部API
func (s *ApiService) GetAllApis() ([]model.SysApi, error) {
	var apis []model.SysApi
	if err := global.DB.Order("`group` ASC, id ASC").Find(&apis).Error; err != nil {
		return nil, err
	}
	return apis, nil
}

// 获取API详情
func (s *ApiService) GetApi(id uint) (*model.SysApi, error) {
	var api model.SysApi
	if err := global.DB.First(&api, id).Error; err != nil {
		return nil, err
	}
	return &api, nil
}

// 创建API
func (s *ApiService) CreateApi(req *request.CreateApiRequest) error {
	// 检查API是否存在
	var count int64
	global.DB.Model(&model.SysApi{}).Where("path = ? AND method = ?", req.Path, req.Method).Count(&count)
	if count > 0 {
		return errors.New("API已存在")
	}

	api := model.SysApi{
		Path:        req.Path,
		Method:      req.Method,
		Group:       req.Group,
		Description: req.Description,
	}

	return global.DB.Create(&api).Error
}

// 更新API
func (s *ApiService) UpdateApi(id uint, req *request.UpdateApiRequest) error {
	var api model.SysApi
	if err := global.DB.First(&api, id).Error; err != nil {
		return errors.New("API不存在")
	}

	// 检查是否与其他API重复
	if req.Path != "" && req.Method != "" {
		var count int64
		global.DB.Model(&model.SysApi{}).Where("path = ? AND method = ? AND id != ?", req.Path, req.Method, id).Count(&count)
		if count > 0 {
			return errors.New("API已存在")
		}
	}

	updates := map[string]interface{}{
		"path":        req.Path,
		"method":      req.Method,
		"group":       req.Group,
		"description": req.Description,
	}

	return global.DB.Model(&api).Updates(updates).Error
}

// 删除API
func (s *ApiService) DeleteApi(id uint) error {
	return global.DB.Delete(&model.SysApi{}, id).Error
}

// ApiGroupStats API分组统计结果
type ApiGroupStats struct {
	Group    string `json:"group"`
	ApiCount int64  `json:"api_count"`
}

// 获取API分组列表(含数量统计)
func (s *ApiService) GetApiGroups() ([]ApiGroupStats, error) {
	var results []ApiGroupStats
	if err := global.DB.Model(&model.SysApi{}).
		Select("`group`, COUNT(*) as api_count").
		Group("`group`").
		Order("`group` ASC").
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// 同步API路由到数据库
func (s *ApiService) SyncApis(routes []model.SysApi) (int, int, int, error) {
	var added, updated, deleted int

	// 构建当前路由的map，用于快速查找
	routeMap := make(map[string]model.SysApi)
	for _, route := range routes {
		key := route.Method + ":" + route.Path
		routeMap[key] = route
	}

	// 获取数据库中所有API
	var existingApis []model.SysApi
	if err := global.DB.Find(&existingApis).Error; err != nil {
		return 0, 0, 0, err
	}

	// 检查数据库中的API是否还存在于路由中
	for _, existing := range existingApis {
		key := existing.Method + ":" + existing.Path
		if _, ok := routeMap[key]; !ok {
			// 路由中不存在，删除数据库记录
			if err := global.DB.Delete(&existing).Error; err != nil {
				return added, updated, deleted, err
			}
			deleted++
		}
	}

	// 添加或更新路由
	for _, route := range routes {
		var existing model.SysApi
		err := global.DB.Where("path = ? AND method = ?", route.Path, route.Method).First(&existing).Error
		if err != nil {
			// 不存在，创建新记录
			if err := global.DB.Create(&route).Error; err != nil {
				return added, updated, deleted, err
			}
			added++
		} else {
			// 已存在，检查是否需要更新
			needUpdate := existing.Group != route.Group ||
				existing.Description != route.Description ||
				existing.RequestParams != route.RequestParams ||
				existing.ResponseParams != route.ResponseParams ||
				existing.NeedAuth != route.NeedAuth

			if needUpdate {
				global.DB.Model(&existing).Updates(map[string]interface{}{
					"group":           route.Group,
					"description":     route.Description,
					"request_params":  route.RequestParams,
					"response_params": route.ResponseParams,
					"need_auth":       route.NeedAuth,
				})
				updated++
			}
		}
	}

	return added, updated, deleted, nil
}
