package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"go-base-server/generator"
	"go-base-server/global"
	"go-base-server/model"
)

type GeneratorService struct{}

var Generator = new(GeneratorService)

// GetTables 获取数据库表列表
func (s *GeneratorService) GetTables() ([]generator.TableInfo, error) {
	return generator.GetTables()
}

// GetTableColumns 获取表字段信息
func (s *GeneratorService) GetTableColumns(tableName string) ([]generator.ColumnConfig, error) {
	columns, err := generator.GetTableColumns(tableName)
	if err != nil {
		return nil, err
	}
	return generator.ConvertToColumnConfig(columns), nil
}

// Preview 预览生成的代码
func (s *GeneratorService) Preview(config *generator.GeneratorConfig) (*generator.PreviewResult, error) {
	serverPath, webPath := s.getPaths(config.FrontendPath)
	gen := generator.NewGenerator(config, serverPath, webPath)
	// 设置父菜单路径
	if config.MenuConfig != nil && config.MenuConfig.ParentID > 0 {
		parentPath := s.getParentMenuPath(config.MenuConfig.ParentID)
		gen.SetParentMenuPath(parentPath)
	}
	return gen.Preview()
}

// Generate 生成代码
func (s *GeneratorService) Generate(config *generator.GeneratorConfig) (*generator.PreviewResult, error) {
	serverPath, webPath := s.getPaths(config.FrontendPath)
	gen := generator.NewGenerator(config, serverPath, webPath)

	// 获取父菜单路径
	parentMenuPath := ""
	if config.MenuConfig != nil && config.MenuConfig.ParentID > 0 {
		parentMenuPath = s.getParentMenuPath(config.MenuConfig.ParentID)
		gen.SetParentMenuPath(parentMenuPath)
	}

	result, err := gen.Generate()
	if err != nil {
		return nil, err
	}

	// 入库保存生成配置（使用upsert逻辑）
	_, _ = s.SaveConfig(config)

	// 自动执行建表SQL
	if config.GenerateSQL {
		for _, file := range result.Files {
			if file.Type == "sql" && !strings.HasSuffix(file.Path, "_menu.sql") {
				// 执行建表SQL，CREATE TABLE IF NOT EXISTS 不会重复创建
				_ = s.ExecuteSQL(file.Content)
			}
		}
	}

	// 创建菜单（覆盖逻辑：先删除同路径菜单及其子菜单，再新建）
	// 仅在生成前端代码且配置了菜单时执行
	if config.GenerateFrontend && config.MenuConfig != nil && config.MenuConfig.MenuName != "" {
		menuPath := "/" + config.ModuleName

		// 删除已存在的同路径菜单及其子菜单
		var oldMenu model.SysMenu
		if err := global.DB.Where("path = ?", menuPath).First(&oldMenu).Error; err == nil {
			// 获取子菜单ID
			var childMenuIds []uint
			global.DB.Model(&model.SysMenu{}).Where("parent_id = ?", oldMenu.ID).Pluck("id", &childMenuIds)
			// 收集所有要删除的菜单ID
			menuIds := append(childMenuIds, oldMenu.ID)
			// 删除角色-菜单关联
			global.DB.Exec("DELETE FROM sys_role_menu WHERE sys_menu_id IN ?", menuIds)
			// 删除子菜单（按钮权限）
			global.DB.Where("parent_id = ?", oldMenu.ID).Delete(&model.SysMenu{})
			// 删除主菜单
			global.DB.Delete(&oldMenu)
		}

		// 构建组件路径
		componentPath := config.ModuleName + "/index"
		if parentMenuPath != "" {
			componentPath = parentMenuPath + "/" + config.ModuleName + "/index"
		}

		menu := model.SysMenu{
			ParentID:   config.MenuConfig.ParentID,
			Name:       config.MenuConfig.MenuName,
			Path:       menuPath,
			Component:  componentPath,
			Icon:       config.MenuConfig.MenuIcon,
			Sort:       config.MenuConfig.MenuSort,
			Type:       2, // 页面菜单
			Permission: config.MenuConfig.Permission + ":list",
			Status:     1,
			Hidden:     0,
		}
		if err := global.DB.Create(&menu).Error; err == nil {
			// 创建按钮权限
			btnPermissions := []struct {
				Name string
				Perm string
				Sort int
			}{
				{"查看", ":list", 1},
				{"新增", ":add", 2},
				{"编辑", ":edit", 3},
				{"删除", ":delete", 4},
			}
			// 如果启用了导入导出功能，添加导出导入按钮
			if config.EnableImportExport {
				btnPermissions = append(btnPermissions, struct {
					Name string
					Perm string
					Sort int
				}{"导出", ":export", 5})
				btnPermissions = append(btnPermissions, struct {
					Name string
					Perm string
					Sort int
				}{"导入", ":import", 6})
			}
			// 如果启用了审批功能，添加审批按钮
			if config.HasAudit {
				btnPermissions = append(btnPermissions, struct {
					Name string
					Perm string
					Sort int
				}{"审批", ":audit", 7})
			}
			for _, btn := range btnPermissions {
				global.DB.Create(&model.SysMenu{
					ParentID:   menu.ID,
					Name:       btn.Name,
					Type:       3, // 按钮
					Permission: config.MenuConfig.Permission + btn.Perm,
					Sort:       btn.Sort,
					Status:     1,
				})
			}
		}
	}

	return result, nil
}

// DeleteModule 删除模块
func (s *GeneratorService) DeleteModule(moduleName string) error {
	serverPath, webPath := s.getPaths("")

	// 删除菜单及子菜单
	menuPath := "/" + moduleName
	var menu model.SysMenu
	if err := global.DB.Where("path = ?", menuPath).First(&menu).Error; err == nil {
		// 获取子菜单ID
		var childMenuIds []uint
		global.DB.Model(&model.SysMenu{}).Where("parent_id = ?", menu.ID).Pluck("id", &childMenuIds)
		// 收集所有要删除的菜单ID
		menuIds := append(childMenuIds, menu.ID)
		// 删除角色-菜单关联
		global.DB.Exec("DELETE FROM sys_role_menu WHERE sys_menu_id IN ?", menuIds)
		// 删除子菜单（按钮权限）
		global.DB.Where("parent_id = ?", menu.ID).Delete(&model.SysMenu{})
		// 删除主菜单
		global.DB.Delete(&menu)
		// 清除所有用户缓存（菜单变更）
		Cache.ClearAllUserInfoCache()
	}

	// 获取配置信息，找到表名和父菜单路径
	var parentMenuPath string
	var genConfig model.SysGenerator
	if err := global.DB.Where("module_name = ?", moduleName).First(&genConfig).Error; err == nil {
		// 解析配置
		var config generator.GeneratorConfig
		if err := json.Unmarshal([]byte(genConfig.ConfigJSON), &config); err == nil {
			// 删除多对多中间表
			for _, rel := range config.Relations {
				if rel.RelationType == "many2many" && rel.JoinTable != "" {
					_ = global.DB.Exec("DROP TABLE IF EXISTS `" + rel.JoinTable + "`").Error
				}
			}
			// 删除主表
			if genConfig.GenTableName != "" {
				_ = global.DB.Exec("DROP TABLE IF EXISTS `" + genConfig.GenTableName + "`").Error
			}
			// 获取父菜单路径
			if config.MenuConfig != nil && config.MenuConfig.ParentID > 0 {
				parentMenuPath = s.getParentMenuPath(config.MenuConfig.ParentID)
			}
		}
	}

	// 删除代码文件
	return generator.DeleteModuleWithParentPath(moduleName, parentMenuPath, serverPath, webPath)
}

// GetGeneratedModules 获取已生成的模块列表
func (s *GeneratorService) GetGeneratedModules() ([]string, error) {
	serverPath, _ := s.getPaths("")
	return generator.GetGeneratedModules(serverPath)
}

// SaveConfig 保存配置（不生成代码）- 根据ID更新或新建
func (s *GeneratorService) SaveConfig(config *generator.GeneratorConfig) (*model.SysGenerator, error) {
	b, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	// 如果有ID，则更新
	if config.ID > 0 {
		var existing model.SysGenerator
		err = global.DB.First(&existing, config.ID).Error
		if err != nil {
			return nil, err
		}
		// 更新
		err = global.DB.Model(&existing).Updates(map[string]interface{}{
			"table_name":  config.TableName,
			"module_name": config.ModuleName,
			"description": config.Description,
			"config_json": string(b),
		}).Error
		if err != nil {
			return nil, err
		}
		return &existing, nil
	}

	// 没有ID，新建
	rec := model.SysGenerator{
		GenTableName: config.TableName,
		ModuleName:   config.ModuleName,
		Description:  config.Description,
		ConfigJSON:   string(b),
	}
	err = global.DB.Create(&rec).Error
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

// GetConfigs 获取已保存的配置列表
func (s *GeneratorService) GetConfigs() ([]model.SysGenerator, error) {
	var list []model.SysGenerator
	err := global.DB.Order("id DESC").Find(&list).Error
	return list, err
}

// GetConfig 获取单个配置详情
func (s *GeneratorService) GetConfig(id string) (*generator.GeneratorConfig, error) {
	var rec model.SysGenerator
	if err := global.DB.First(&rec, id).Error; err != nil {
		return nil, err
	}
	var config generator.GeneratorConfig
	if err := json.Unmarshal([]byte(rec.ConfigJSON), &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// DeleteConfig 删除配置
func (s *GeneratorService) DeleteConfig(id string) error {
	return global.DB.Delete(&model.SysGenerator{}, id).Error
}

// ExecuteSQL 执行建表SQL
func (s *GeneratorService) ExecuteSQL(sql string) error {
	return global.DB.Exec(sql).Error
}

// getParentMenuPath 获取父菜单路径(组件路径前缀)
func (s *GeneratorService) getParentMenuPath(parentID uint) string {
	if parentID == 0 {
		return ""
	}
	var parent model.SysMenu
	if err := global.DB.First(&parent, parentID).Error; err != nil {
		return ""
	}
	// 只有目录(type=1)和菜单(type=2)才有组件路径
	if parent.Type != 1 && parent.Type != 2 {
		return ""
	}
	// 使用父菜单的组件路径前缀
	if parent.Component != "" {
		// 去除 /index 后缀
		path := strings.TrimSuffix(parent.Component, "/index")
		return path
	}
	// 如果没有组件路径，使用路由路径
	return strings.TrimPrefix(parent.Path, "/")
}

// getPaths 获取项目路径
func (s *GeneratorService) getPaths(customFrontendPath string) (serverPath, webPath string) {
	// 获取当前文件所在目录
	_, filename, _, _ := runtime.Caller(0)
	serverPath = filepath.Dir(filepath.Dir(filename))

	// 如果指定了自定义前端路径，直接使用
	if customFrontendPath != "" {
		webPath = customFrontendPath
		return
	}

	root := filepath.Dir(serverPath)

	// 在父目录中搜索包含 "web" 的目录
	entries, err := os.ReadDir(root)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				name := entry.Name()
				// 优先匹配包含web的目录，排除server目录
				if strings.Contains(strings.ToLower(name), "web") && !strings.Contains(strings.ToLower(name), "server") {
					candidate := filepath.Join(root, name)
					// 检查是否有src目录（前端项目特征）
					if st, err := os.Stat(filepath.Join(candidate, "src")); err == nil && st.IsDir() {
						webPath = candidate
						return
					}
				}
			}
		}
	}

	// 回退: 优先使用 go-base-web，其次 web
	candidate1 := filepath.Join(root, "go-base-web")
	candidate2 := filepath.Join(root, "web")
	if st, err := os.Stat(candidate1); err == nil && st.IsDir() {
		webPath = candidate1
	} else if st, err := os.Stat(candidate2); err == nil && st.IsDir() {
		webPath = candidate2
	} else {
		webPath = candidate2
	}
	return
}
