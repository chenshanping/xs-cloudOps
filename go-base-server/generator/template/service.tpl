package service

import (
	"errors"
{{- if .HasMultiFiles}}
	"strings"
	"strconv"
{{- end}}
{{- if or .HasAudit .UniqueColumns}}
	"time"
{{- end}}
{{- if or .HasMany2Many .UniqueColumns .LinkToUser}}

	"gorm.io/gorm"
{{- end}}

	"go-base-server/global"
	"go-base-server/model"
	"go-base-server/model/request"
)

type {{.ModelName}}Service struct{}

var {{.ModelName}} = new({{.ModelName}}Service)

// Get{{.ModelName}}List 获取{{.Description}}列表
{{- if .DataIsolation}}
// userID: 当前用户ID, isAdmin: 是否为管理员
func (s *{{.ModelName}}Service) Get{{.ModelName}}List(req *request.{{.ModelName}}ListRequest, userID uint, isAdmin bool) ([]model.{{.ModelName}}, int64, error) {
{{- else}}
func (s *{{.ModelName}}Service) Get{{.ModelName}}List(req *request.{{.ModelName}}ListRequest) ([]model.{{.ModelName}}, int64, error) {
{{- end}}
	var list []model.{{.ModelName}}
	var total int64

	db := global.DB.Model(&model.{{.ModelName}}{})
{{- if .DataIsolation}}

	// 数据隔离：非管理员只能看到自己创建的数据
	if !isAdmin {
		db = db.Where("created_by = ?", userID)
	}
{{- end}}

{{- range .SearchColumns}}
{{- if eq .SearchType "eq"}}
	if req.{{.FieldName}} != nil {
		db = db.Where("{{.ColumnName}} = ?", *req.{{.FieldName}})
	}
{{- else if eq .SearchType "like"}}
	if req.{{.FieldName}} != "" {
		db = db.Where("{{.ColumnName}} LIKE ?", "%"+req.{{.FieldName}}+"%")
	}
{{- else if eq .SearchType "gt"}}
	if req.{{.FieldName}} != nil {
		db = db.Where("{{.ColumnName}} > ?", *req.{{.FieldName}})
	}
{{- else if eq .SearchType "gte"}}
	if req.{{.FieldName}} != nil {
		db = db.Where("{{.ColumnName}} >= ?", *req.{{.FieldName}})
	}
{{- else if eq .SearchType "lt"}}
	if req.{{.FieldName}} != nil {
		db = db.Where("{{.ColumnName}} < ?", *req.{{.FieldName}})
	}
{{- else if eq .SearchType "lte"}}
	if req.{{.FieldName}} != nil {
		db = db.Where("{{.ColumnName}} <= ?", *req.{{.FieldName}})
	}
{{- else}}
{{- if eq .FieldType "string"}}
	if req.{{.FieldName}} != "" {
		db = db.Where("{{.ColumnName}} LIKE ?", "%"+req.{{.FieldName}}+"%")
	}
{{- else}}
	if req.{{.FieldName}} != nil {
		db = db.Where("{{.ColumnName}} = ?", *req.{{.FieldName}})
	}
{{- end}}
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
	if req.{{.ForeignKey | ToPascalCase}} != nil {
		db = db.Where("{{.ForeignKeyJson}} = ?", *req.{{.ForeignKey | ToPascalCase}})
	}
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
	if req.CreatedBy != nil {
		db = db.Where("created_by = ?", *req.CreatedBy)
	}
{{- end}}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := req.GetOffset()
	// 排序处理
	orderBy := "id DESC"
	if req.SortField != "" && req.SortOrder != "" {
		// 前端传入排序
		allowedFields := map[string]bool{
			"id": true,
{{- range .SortColumns}}
			"{{.ColumnName}}": true,
{{- end}}
			"created_at": true,
		}
		if allowedFields[req.SortField] {
			order := "ASC"
			if req.SortOrder == "descend" || req.SortOrder == "desc" {
				order = "DESC"
			}
			orderBy = req.SortField + " " + order
		}
	}
	query := db.Offset(offset).Limit(req.PageSize).Order(orderBy)
{{- if .HasPreloads}}
{{- range .Preloads}}
	query = query.Preload("{{.}}")
{{- end}}
{{- end}}
{{- range .FormColumns}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
	query = query.Preload("{{.FieldName}}File")
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
	query = query.Preload("Creator.AvatarFile")
{{- end}}
{{- if .HasCreatedByProfile}}
	// 创建者身份信息需要在循环中单独查询（因为关联是通过 user_id 而不是 created_by）
{{- end}}
{{- if .HasAudit}}
	query = query.Preload("Auditor.AvatarFile")
{{- end}}
	if err := query.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	// 填充文件URL
	for i := range list {
		list[i].FillFileURLs()
{{- if .HasMultiFiles}}
		s.fillMultiFileURLs(&list[i])
{{- end}}
	}
{{- if .HasCreatedByProfile}}

	// 填充创建者身份信息
	s.fillCreatorProfiles(list)
{{- end}}

	return list, total, nil
}

// Get{{.ModelName}} 获取{{.Description}}详情
func (s *{{.ModelName}}Service) Get{{.ModelName}}(id uint) (*model.{{.ModelName}}, error) {
	var data model.{{.ModelName}}
	query := global.DB
{{- if .HasPreloads}}
{{- range .Preloads}}
	query = query.Preload("{{.}}")
{{- end}}
{{- end}}
{{- range .FormColumns}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
	query = query.Preload("{{.FieldName}}File")
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
	query = query.Preload("Creator.AvatarFile")
{{- end}}
{{- if .HasAudit}}
	query = query.Preload("Auditor.AvatarFile")
{{- end}}
	if err := query.First(&data, id).Error; err != nil {
		return nil, err
	}
	data.FillFileURLs()
{{- if .HasMultiFiles}}
	s.fillMultiFileURLs(&data)
{{- end}}
	return &data, nil
}

// Create{{.ModelName}} 创建{{.Description}}
func (s *{{.ModelName}}Service) Create{{.ModelName}}(req *request.Create{{.ModelName}}Request, userID uint) error {
{{- if .UniqueColumns}}
	// 唯一性校验
	var count int64
{{- range .UniqueColumns}}
	global.DB.Model(&model.{{$.ModelName}}{}).Where("{{.ColumnName}} = ?", req.{{.FieldName}}).Count(&count)
	if count > 0 {
		return errors.New("{{.Comment}}已存在")
	}
{{- end}}
{{- end}}

	data := model.{{.ModelName}}{
{{- if .LinkToUser}}
		UserID: req.UserID,
{{- end}}
{{- if .HasCreatedBy}}
		CreatedBy: userID,
{{- end}}
{{- range .FormColumns}}
    {{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
		{{.FieldName}}FileID: req.{{.FieldName}}FileID,
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
		{{.FieldName}}FileIDs: req.{{.FieldName}}FileIDs,
{{- else}}
		{{.FieldName}}: {{if and .IsRequired (or (eq .FieldType "int") (eq .FieldType "int64") (eq .FieldType "uint") (eq .FieldType "float64"))}}*{{end}}req.{{.FieldName}},
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
		{{.ForeignKey | ToPascalCase}}: req.{{.ForeignKey | ToPascalCase}},
{{- end}}
{{- end}}
	}

{{- if .HasMany2Many}}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&data).Error; err != nil {
			return err
		}
{{- range .Relations}}
{{- if eq .RelationType "many2many"}}
		if len(req.{{.FieldName}}Ids) > 0 {
			var {{.JsonName}} []model.{{.RelatedModel}}
			if err := tx.Where("id IN ?", req.{{.FieldName}}Ids).Find(&{{.JsonName}}).Error; err != nil {
				return err
			}
			if err := tx.Model(&data).Association("{{.FieldName}}").Replace({{.JsonName}}); err != nil {
				return err
			}
		}
{{- end}}
{{- end}}
		return nil
	})
{{- else}}
	return global.DB.Create(&data).Error
{{- end}}
}

// Update{{.ModelName}} 更新{{.Description}}
func (s *{{.ModelName}}Service) Update{{.ModelName}}(id uint, req *request.Update{{.ModelName}}Request) error {
	var data model.{{.ModelName}}
	if err := global.DB.First(&data, id).Error; err != nil {
		return errors.New("数据不存在")
	}

{{- if .UniqueColumns}}
	// 唯一性校验（排除自己）
	var count int64
{{- range .UniqueColumns}}
	global.DB.Model(&model.{{$.ModelName}}{}).Where("{{.ColumnName}} = ? AND id != ?", req.{{.FieldName}}, id).Count(&count)
	if count > 0 {
		return errors.New("{{.Comment}}已存在")
	}
{{- end}}
{{- end}}

	updates := map[string]interface{}{
{{- if .LinkToUser}}
		"user_id": req.UserID,
{{- end}}
{{- range .FormColumns}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
		"{{.ColumnName}}_file_id": req.{{.FieldName}}FileID,
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
		"{{.ColumnName}}_file_ids": req.{{.FieldName}}FileIDs,
{{- else}}
		"{{.ColumnName}}": req.{{.FieldName}},
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
		"{{.ForeignKeyJson}}": req.{{.ForeignKey | ToPascalCase}},
{{- end}}
{{- end}}
	}

{{- if .HasMany2Many}}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&data).Updates(updates).Error; err != nil {
			return err
		}
{{- range .Relations}}
{{- if eq .RelationType "many2many"}}
		var {{.JsonName}} []model.{{.RelatedModel}}
		if len(req.{{.FieldName}}Ids) > 0 {
			if err := tx.Where("id IN ?", req.{{.FieldName}}Ids).Find(&{{.JsonName}}).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&data).Association("{{.FieldName}}").Replace({{.JsonName}}); err != nil {
			return err
		}
{{- end}}
{{- end}}
		return nil
	})
{{- else}}
	return global.DB.Model(&data).Updates(updates).Error
{{- end}}
}

{{- if .HasMultiFiles}}
// fillMultiFileURLs 填充多文件URL和名称
func (s *{{.ModelName}}Service) fillMultiFileURLs(data *model.{{.ModelName}}) {
{{- range .FormColumns}}
{{- if or (eq .FormType "images") (eq .FormType "files")}}
	if data.{{.FieldName}}FileIDs != "" {
		ids := strings.Split(data.{{.FieldName}}FileIDs, ",")
		uintIds := make([]uint, 0, len(ids))
		for _, idStr := range ids {
			if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64); err == nil {
				uintIds = append(uintIds, uint(id))
			}
		}
		if len(uintIds) > 0 {
			var files []model.SysFile
			global.DB.Where("id IN ?", uintIds).Find(&files)
			// 保持原有顺序
			urlMap := make(map[uint]string)
			nameMap := make(map[uint]string)
			for _, f := range files {
				urlMap[f.ID] = f.URL
				nameMap[f.ID] = f.Name
			}
			data.{{.FieldName}}URLs = make([]string, 0, len(uintIds))
			data.{{.FieldName}}Names = make([]string, 0, len(uintIds))
			for _, id := range uintIds {
				if url, ok := urlMap[id]; ok {
					data.{{.FieldName}}URLs = append(data.{{.FieldName}}URLs, url)
					data.{{.FieldName}}Names = append(data.{{.FieldName}}Names, nameMap[id])
				}
			}
		}
	}
{{- end}}
{{- end}}
}
{{- end}}
{{- if .HasCreatedByProfile}}

// fillCreatorProfiles 填充创建者身份信息
func (s *{{.ModelName}}Service) fillCreatorProfiles(list []model.{{.ModelName}}) {
	if len(list) == 0 {
		return
	}
	
	// 收集所有创建者 ID
	userIDs := make([]uint, 0, len(list))
	for _, item := range list {
		if item.CreatedBy > 0 {
			userIDs = append(userIDs, item.CreatedBy)
		}
	}
	if len(userIDs) == 0 {
		return
	}
	
	// 根据 user_id 查询身份信息
	var profiles []model.{{.CreatedByProfileModel}}
	global.DB.Where("user_id IN ?", userIDs).Find(&profiles)
	
	// 构建 user_id -> profile 的映射
	profileMap := make(map[uint]*model.{{.CreatedByProfileModel}})
	for i := range profiles {
		profiles[i].FillFileURLs()
		profileMap[profiles[i].UserID] = &profiles[i]
	}
	
	// 填充到列表
	for i := range list {
		if profile, ok := profileMap[list[i].CreatedBy]; ok {
			list[i].CreatorProfile = profile
		}
	}
}
{{- end}}

// Delete{{.ModelName}} 删除{{.Description}}
func (s *{{.ModelName}}Service) Delete{{.ModelName}}(id uint) error {
	var data model.{{.ModelName}}
	if err := global.DB.First(&data, id).Error; err != nil {
		return errors.New("数据不存在")
	}

{{- if or .HasMany2Many .UniqueColumns}}
	return global.DB.Transaction(func(tx *gorm.DB) error {
{{- if .HasMany2Many}}
{{- range .Relations}}
{{- if eq .RelationType "many2many"}}
		if err := tx.Model(&data).Association("{{.FieldName}}").Clear(); err != nil {
			return err
		}
{{- end}}
{{- end}}
{{- end}}
{{- range .UniqueColumns}}
		// 软删除前修改 {{.ColumnName}}，避免唯一索引冲突
{{- if eq .FieldType "string"}}
		deleted{{.FieldName}} := data.{{.FieldName}} + "_deleted_" + time.Now().Format("20060102150405")
		if err := tx.Model(&data).Update("{{.ColumnName}}", deleted{{.FieldName}}).Error; err != nil {
			return err
		}
{{- end}}
{{- end}}
		return tx.Delete(&data).Error
	})
{{- else}}
	return global.DB.Delete(&data).Error
{{- end}}
}

// BatchDelete{{.ModelName}} 批量删除{{.Description}}
func (s *{{.ModelName}}Service) BatchDelete{{.ModelName}}(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
{{- if .UniqueColumns}}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			var data model.{{.ModelName}}
			if err := tx.First(&data, id).Error; err != nil {
				continue // 跳过不存在的记录
			}
{{- range .UniqueColumns}}
			// 软删除前修改 {{.ColumnName}}，避免唯一索引冲突
{{- if eq .FieldType "string"}}
			deleted{{.FieldName}} := data.{{.FieldName}} + "_deleted_" + time.Now().Format("20060102150405")
			if err := tx.Model(&data).Update("{{.ColumnName}}", deleted{{.FieldName}}).Error; err != nil {
				return err
			}
{{- end}}
{{- end}}
			if err := tx.Delete(&data).Error; err != nil {
				return err
			}
		}
		return nil
	})
{{- else}}
	return global.DB.Where("id IN ?", ids).Delete(&model.{{.ModelName}}{}).Error
{{- end}}
}

// Get{{.ModelName}}Options 获取{{.Description}}选项列表（带可选关联统计）
// excludeDeleted: 是否排除软删除数据（统计表有deleted_at字段时传true）
// countCreatedBy: 统计时按创建人过滤（数据隔离用，传当前用户ID，0表示不过滤）
{{- if .DataIsolation}}
func (s *{{.ModelName}}Service) Get{{.ModelName}}Options(displayField, countTable, countForeignKey string, excludeDeleted bool, countCreatedBy uint, userID uint, isAdmin bool) ([]map[string]interface{}, error) {
{{- else}}
func (s *{{.ModelName}}Service) Get{{.ModelName}}Options(displayField, countTable, countForeignKey string, excludeDeleted bool, countCreatedBy uint) ([]map[string]interface{}, error) {
{{- end}}
	var results []map[string]interface{}

	if displayField == "" {
		displayField = "name"
	}

	// 无统计关联时，简单查询
	if countTable == "" || countForeignKey == "" {
{{- if and .DataIsolation .HasCreatedBy}}
		db := global.DB.Model(&model.{{.ModelName}}{})
		if !isAdmin {
			db = db.Where("created_by = ?", userID)
		}
		err := db.Select("id, `"+displayField+"` as name").Order("id ASC").Find(&results).Error
{{- else}}
		err := global.DB.Model(&model.{{.ModelName}}{}).
			Select("id, `" + displayField + "` as name").
			Order("id ASC").
			Find(&results).Error
{{- end}}
		// 转换[]byte为string（GORM返回map时字符串字段会是[]byte类型）
		for i := range results {
			if nameBytes, ok := results[i]["name"].([]byte); ok {
				results[i]["name"] = string(nameBytes)
			}
		}
		return results, err
	}

	// 有统计关联时，使用子查询
	subQuery := global.DB.Table(countTable).
		Select(countForeignKey + " as fk, COUNT(*) as cnt")

	// 排除软删除数据
	if excludeDeleted {
		subQuery = subQuery.Where("deleted_at IS NULL")
	}
	// 数据隔离：统计时按创建人过滤
	if countCreatedBy > 0 {
		subQuery = subQuery.Where("created_by = ?", countCreatedBy)
	}
	subQuery = subQuery.Group(countForeignKey)

{{- if and .DataIsolation .HasCreatedBy}}
	query := global.DB.Table("{{.TableName}}")
	if !isAdmin {
		query = query.Where("{{.TableName}}.created_by = ?", userID)
	}
	err := query.
		Select("{{.TableName}}.id, {{.TableName}}.`" + displayField + "` as name, COALESCE(sub.cnt, 0) as count").
		Joins("LEFT JOIN (?) as sub ON {{.TableName}}.id = sub.fk", subQuery).
		Order("{{.TableName}}.id ASC").
		Find(&results).Error
{{- else}}
	err := global.DB.Table("{{.TableName}}").
		Select("{{.TableName}}.id, {{.TableName}}.`" + displayField + "` as name, COALESCE(sub.cnt, 0) as count").
		Joins("LEFT JOIN (?) as sub ON {{.TableName}}.id = sub.fk", subQuery).
		Order("{{.TableName}}.id ASC").
		Find(&results).Error
{{- end}}
	// 转换[]byte为string
	for i := range results {
		if nameBytes, ok := results[i]["name"].([]byte); ok {
			results[i]["name"] = string(nameBytes)
		}
	}

	return results, err
}
{{- if .HasCreatedBy}}

// Get{{.ModelName}}CreatorOptions 获取创建人选项列表（从当前表 GROUP BY created_by 关联用户表）
func (s *{{.ModelName}}Service) Get{{.ModelName}}CreatorOptions() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 从当前表获取有数据的创建人，关联用户表获取名称
	err := global.DB.Table("{{.TableName}}").
		Select("sys_user.id, COALESCE(NULLIF(sys_user.nickname, ''), sys_user.username) as name, COUNT({{.TableName}}.id) as count").
		Joins("LEFT JOIN sys_user ON {{.TableName}}.created_by = sys_user.id").
{{- if .HasDeletedAt}}
		Where("{{.TableName}}.deleted_at IS NULL").
{{- end}}
		Group("{{.TableName}}.created_by").
		Order("count DESC").
		Find(&results).Error
	// 转换[]byte为string
	for i := range results {
		if nameBytes, ok := results[i]["name"].([]byte); ok {
			results[i]["name"] = string(nameBytes)
		}
	}

	return results, err
}
{{- end}}
{{- if .HasAudit}}

// Audit{{.ModelName}} 审批{{.Description}}
func (s *{{.ModelName}}Service) Audit{{.ModelName}}(id uint, req *request.Audit{{.ModelName}}Request, userID uint) error {
	var data model.{{.ModelName}}
	if err := global.DB.First(&data, id).Error; err != nil {
		return errors.New("数据不存在")
	}

	// 检查是否已审批
	if data.AuditStatus != 0 {
		return errors.New("该数据已审批，不可重复审批")
	}

	updates := map[string]interface{}{
		"audit_status": req.AuditStatus,
		"audit_remark": req.AuditRemark,
		"audit_time":   time.Now(),
		"audit_by":     userID,
	}

	return global.DB.Model(&data).Updates(updates).Error
}
{{- end}}
{{- if .GenerateFrontendApi}}

// GetFrontend{{.ModelName}}List 获取前台{{.Description}}列表（前台用户使用，不做 created_by 过滤{{if .HasAudit}}，仅返回已启用且审批通过的数据{{end}}）
func (s *{{.ModelName}}Service) GetFrontend{{.ModelName}}List(req *request.Frontend{{.ModelName}}ListRequest) ([]model.{{.ModelName}}, int64, error) {
	var list []model.{{.ModelName}}
	var total int64

	db := global.DB.Model(&model.{{.ModelName}}{})
{{- if .HasAudit}}

	// 只查询已启用且审批通过的数据
	db = db.Where("status = ? AND audit_status = ?", 1, 1)
{{- end}}

{{- range .SearchColumns}}
{{- if ne .ColumnName "status"}}
{{- if eq .SearchType "eq"}}
	if req.{{.FieldName}} != nil {
		db = db.Where("{{.ColumnName}} = ?", *req.{{.FieldName}})
	}
{{- else if eq .SearchType "like"}}
	if req.{{.FieldName}} != "" {
		db = db.Where("{{.ColumnName}} LIKE ?", "%"+req.{{.FieldName}}+"%")
	}
{{- else}}
{{- if eq .FieldType "string"}}
	if req.{{.FieldName}} != "" {
		db = db.Where("{{.ColumnName}} LIKE ?", "%"+req.{{.FieldName}}+"%")
	}
{{- else}}
	if req.{{.FieldName}} != nil {
		db = db.Where("{{.ColumnName}} = ?", *req.{{.FieldName}})
	}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
	if req.{{.ForeignKey | ToPascalCase}} != nil {
		db = db.Where("{{.ForeignKeyJson}} = ?", *req.{{.ForeignKey | ToPascalCase}})
	}
{{- end}}
{{- end}}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := req.GetOffset()
	query := db.Offset(offset).Limit(req.PageSize).Order("sort ASC, id DESC")
{{- if .HasPreloads}}
{{- range .Preloads}}
	query = query.Preload("{{.}}")
{{- end}}
{{- end}}
{{- range .FormColumns}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
	query = query.Preload("{{.FieldName}}File")
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
	query = query.Preload("Creator.AvatarFile")
{{- end}}
	if err := query.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	// 填充文件URL
	for i := range list {
		list[i].FillFileURLs()
{{- if .HasMultiFiles}}
		s.fillMultiFileURLs(&list[i])
{{- end}}
	}

	return list, total, nil
}

// GetFrontend{{.ModelName}} 获取前台{{.Description}}详情（前台用户使用{{if .HasAudit}}，仅返回已启用且审批通过的数据{{end}}）
func (s *{{.ModelName}}Service) GetFrontend{{.ModelName}}(id uint) (*model.{{.ModelName}}, error) {
	var data model.{{.ModelName}}
{{- if .HasAudit}}
	query := global.DB.Where("status = ? AND audit_status = ?", 1, 1)
{{- else}}
	query := global.DB
{{- end}}
{{- if .HasPreloads}}
{{- range .Preloads}}
	query = query.Preload("{{.}}")
{{- end}}
{{- end}}
{{- range .FormColumns}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
	query = query.Preload("{{.FieldName}}File")
{{- end}}
{{- end}}
{{- if .HasCreatedBy}}
	query = query.Preload("Creator.AvatarFile")
{{- end}}
	if err := query.First(&data, id).Error; err != nil {
		return nil, err
	}
	data.FillFileURLs()
{{- if .HasMultiFiles}}
	s.fillMultiFileURLs(&data)
{{- end}}
	return &data, nil
}
{{- end}}
{{- if .HasStats}}
{{- range $chart := .StatsCharts}}

// Get{{$.ModelName}}Stats{{$chart.Field}} 获取{{$.Description}}按{{$chart.Title}}分组统计
func (s *{{$.ModelName}}Service) Get{{$.ModelName}}Stats{{$chart.Field}}() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := global.DB.Table("{{$.TableName}}").
		Select("{{$chart.Column}} as group_key, COUNT(*) as value").
{{- if $.HasDeletedAt}}
		Where("deleted_at IS NULL").
{{- end}}
		Group("{{$chart.Column}}").
		Order("value DESC").
		Find(&results).Error
	return results, err
}
{{- end}}
{{- if .HasStatsTrend}}

// Get{{.ModelName}}TrendStats 获取{{.Description}}趋势统计
func (s *{{.ModelName}}Service) Get{{.ModelName}}TrendStats(days int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	if days <= 0 {
		days = 30
	}
	err := global.DB.Table("{{.TableName}}").
		Select("DATE({{.StatsTimeColumn}}) as date, COUNT(*) as value").
		Where("{{.StatsTimeColumn}} >= DATE_SUB(CURDATE(), INTERVAL ? DAY)", days).
{{- if .HasDeletedAt}}
		Where("deleted_at IS NULL").
{{- end}}
		Group("DATE({{.StatsTimeColumn}})").
		Order("date ASC").
		Find(&results).Error
	return results, err
}
{{- end}}
{{- end}}
{{- if .LinkToUser}}

// GetMy{{.ModelName}} 获取当前用户的{{.Description}}信息
func (s *{{.ModelName}}Service) GetMy{{.ModelName}}(userID uint) (*model.{{.ModelName}}, error) {
	var data model.{{.ModelName}}
	query := global.DB.Where("user_id = ?", userID)
{{- if .HasPreloads}}
{{- range .Preloads}}
	query = query.Preload("{{.}}")
{{- end}}
{{- end}}
{{- range .FormColumns}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
	query = query.Preload("{{.FieldName}}File")
{{- end}}
{{- end}}
	query = query.Preload("User.AvatarFile")
	if err := query.First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户还没有创建记录，返回 nil
		}
		return nil, err
	}
	data.FillFileURLs()
{{- if .HasMultiFiles}}
	s.fillMultiFileURLs(&data)
{{- end}}
	return &data, nil
}

// SaveMy{{.ModelName}} 保存当前用户的{{.Description}}信息（创建或更新）
func (s *{{.ModelName}}Service) SaveMy{{.ModelName}}(userID uint, req *request.SaveMy{{.ModelName}}Request) error {
	var data model.{{.ModelName}}
	err := global.DB.Where("user_id = ?", userID).First(&data).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新记录
		data = model.{{.ModelName}}{
			UserID: userID,
{{- if .HasCreatedBy}}
			CreatedBy: userID,
{{- end}}
{{- range .FormColumns}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
			{{.FieldName}}FileID: req.{{.FieldName}}FileID,
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
			{{.FieldName}}FileIDs: req.{{.FieldName}}FileIDs,
{{- else}}
			{{.FieldName}}: req.{{.FieldName}},
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
			{{.ForeignKey | ToPascalCase}}: req.{{.ForeignKey | ToPascalCase}},
{{- end}}
{{- end}}
		}
{{- if .HasMany2Many}}
		return global.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&data).Error; err != nil {
				return err
			}
{{- range .Relations}}
{{- if eq .RelationType "many2many"}}
			if len(req.{{.FieldName}}Ids) > 0 {
				var {{.JsonName}} []model.{{.RelatedModel}}
				if err := tx.Where("id IN ?", req.{{.FieldName}}Ids).Find(&{{.JsonName}}).Error; err != nil {
					return err
				}
				if err := tx.Model(&data).Association("{{.FieldName}}").Replace({{.JsonName}}); err != nil {
					return err
				}
			}
{{- end}}
{{- end}}
			return nil
		})
{{- else}}
		return global.DB.Create(&data).Error
{{- end}}
	} else if err != nil {
		return err
	}

	// 更新现有记录
	updates := map[string]interface{}{
{{- range .FormColumns}}
{{- if or (eq .FormType "image") (eq .FormType "file") (eq .FormType "upload")}}
		"{{.ColumnName}}_file_id": req.{{.FieldName}}FileID,
{{- else if or (eq .FormType "images") (eq .FormType "files")}}
		"{{.ColumnName}}_file_ids": req.{{.FieldName}}FileIDs,
{{- else}}
		"{{.ColumnName}}": req.{{.FieldName}},
{{- end}}
{{- end}}
{{- range .Relations}}
{{- if eq .RelationType "belongsTo"}}
		"{{.ForeignKeyJson}}": req.{{.ForeignKey | ToPascalCase}},
{{- end}}
{{- end}}
	}
{{- if .HasAudit}}
	// 如果之前审批被拒绝，重新保存时自动重置为待审批
	if data.AuditStatus == 2 {
		updates["audit_status"] = 0
		updates["audit_remark"] = ""
		updates["audit_time"] = nil
		updates["audit_by"] = 0
	}
{{- end}}

{{- if .HasMany2Many}}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&data).Updates(updates).Error; err != nil {
			return err
		}
{{- range .Relations}}
{{- if eq .RelationType "many2many"}}
		var {{.JsonName}} []model.{{.RelatedModel}}
		if len(req.{{.FieldName}}Ids) > 0 {
			if err := tx.Where("id IN ?", req.{{.FieldName}}Ids).Find(&{{.JsonName}}).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&data).Association("{{.FieldName}}").Replace({{.JsonName}}); err != nil {
			return err
		}
{{- end}}
{{- end}}
		return nil
	})
{{- else}}
	return global.DB.Model(&data).Updates(updates).Error
{{- end}}
}

// GetByUserID 根据用户ID获取{{.Description}}信息（管理员使用）
func (s *{{.ModelName}}Service) GetByUserID(userID uint) (*model.{{.ModelName}}, error) {
	var data model.{{.ModelName}}
	if err := global.DB.Where("user_id = ?", userID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// init 注册用户身份处理器
func init() {
	global.Profiles.Register(&global.ProfileHandler{
		Key:         "{{.ModuleName}}",
		Name:        "{{if .ProfileName}}{{.ProfileName}}{{else}}{{.Description}}{{end}}",
		Description: "{{.Description}}信息",
		Icon:        "{{if .ProfileIcon}}{{.ProfileIcon}}{{else}}UserOutlined{{end}}",
		RoleCode:    "{{.ProfileRoleCode}}",
		Fields: []global.FieldConfig{
{{- range .FormColumns}}
{{- if eq .FormType "image"}}
			{Key: "{{.JsonName}}_url", Label: "{{.Comment}}", Required: {{.IsRequired}}, Type: "image"},
{{- else if eq .FormType "images"}}
			{Key: "{{.JsonName}}_urls", Label: "{{.Comment}}", Required: {{.IsRequired}}, Type: "images"},
{{- else if or (eq .FormType "file") (eq .FormType "upload")}}
			{Key: "{{.JsonName}}_url", Label: "{{.Comment}}", Required: {{.IsRequired}}, Type: "file"},
{{- else if eq .FormType "files"}}
			{Key: "{{.JsonName}}_urls", Label: "{{.Comment}}", Required: {{.IsRequired}}, Type: "files"},
{{- else if eq .FormType "select"}}
			{Key: "{{.JsonName}}", Label: "{{.Comment}}", Required: {{.IsRequired}}, Type: "select"},
{{- else}}
			{Key: "{{.JsonName}}", Label: "{{.Comment}}", Required: {{.IsRequired}}, Type: "text"},
{{- end}}
{{- end}}
{{- if .HasAudit}}
			{Key: "audit_status", Label: "审批状态", Required: false, Type: "select"},
			{Key: "audit_remark", Label: "审批备注", Required: false, Type: "text"},
{{- end}}
		},
		GetProfile: func(userID uint) (interface{}, error) {
			return {{.ModelName}}.GetMy{{.ModelName}}(userID)
		},
	})
}
{{- end}}
