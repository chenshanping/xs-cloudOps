package initialize

import (
	"errors"
	"server/global"
	"server/model"
	"server/utils"

	"gorm.io/gorm"
)

func InitDBTables() {
	err := global.DB.AutoMigrate(
		&model.SysUser{},
		&model.SysRole{},
		&model.SysDept{},
		&model.SysMenu{},
		&model.SysApi{},
		&model.SysOperationLog{},
		&model.SysLoginLog{},
		&model.SysSlowLog{},
		&model.SysConfig{},
		&model.SysStorage{},
		&model.SysFile{},
		&model.SysFileChunk{},
		// 数据字典
		&model.SysDictType{},
		&model.SysDictData{},
		// AI对话相关
		&model.AIConversation{},
		&model.AIMessage{},
	)
	if err != nil {
		panic("数据库表迁移失败: " + err.Error())
	}

	global.Log.Info("数据库表迁移成功")

	// 初始化默认数据
	initDefaultData()

	// 单独初始化系统配置（即使已有用户数据也会执行）
	initDefaultConfigs()

	// 补齐升级场景下缺失的内置配置和 API 权限元数据
	ensureBuiltInData()
}

func initDefaultData() {
	// 若已有用户数据则跳过
	var count int64
	global.DB.Model(&model.SysUser{}).Count(&count)
	if count > 0 {
		return
	}

	// 创建默认角色
	adminRole := model.SysRole{
		Name:      "超级管理员",
		Code:      "admin",
		Sort:      1,
		Status:    1,
		DataScope: model.DataScopeAll,
		Remark:    "拥有所有权限",
	}
	global.DB.Create(&adminRole)

	rootDept := model.SysDept{
		ParentID:  0,
		Ancestors: "0",
		Name:      "平台",
		Sort:      1,
		Status:    1,
		Remark:    "系统根部门",
	}
	global.DB.Create(&rootDept)

	// 创建默认用户
	hashedPassword, _ := utils.HashPassword("123456")
	adminUser := model.SysUser{
		Username: "admin",
		Password: hashedPassword,
		Nickname: "管理员",
		Status:   1,
		DeptID:   rootDept.ID,
		Roles:    []model.SysRole{adminRole},
	}
	global.DB.Create(&adminUser)

	// 根菜单
	sysMgmt := model.SysMenu{ParentID: 0, Name: "系统管理", Path: "/system", Component: "Layout", Icon: "setting", Sort: 1, Type: 1, Status: 1}
	global.DB.Create(&sysMgmt)
	monitor := model.SysMenu{ParentID: 0, Name: "系统监控", Path: "/monitor", Component: "Layout", Icon: "monitor", Sort: 2, Type: 1, Status: 1}
	global.DB.Create(&monitor)

	// 子菜单
	menus := []model.SysMenu{
		{ParentID: sysMgmt.ID, Name: "用户管理", Path: "/system/user", Component: "system/user/index", Icon: "user", Sort: 1, Type: 2, Permission: "system:user:list", Status: 1},
		{ParentID: sysMgmt.ID, Name: "角色管理", Path: "/system/role", Component: "system/role/index", Icon: "team", Sort: 2, Type: 2, Permission: "system:role:list", Status: 1},
		{ParentID: sysMgmt.ID, Name: "部门管理", Path: "/system/dept", Component: "system/dept/index", Icon: "apartment", Sort: 3, Type: 2, Permission: "system:dept:list", Status: 1},
		{ParentID: sysMgmt.ID, Name: "菜单管理", Path: "/system/menu", Component: "system/menu/index", Icon: "menu", Sort: 4, Type: 2, Permission: "system:menu:list", Status: 1},
		{ParentID: sysMgmt.ID, Name: "API管理", Path: "/system/api", Component: "system/api/index", Icon: "api", Sort: 5, Type: 2, Permission: "system:api:list", Status: 1},
		{ParentID: sysMgmt.ID, Name: "参数配置", Path: "/system/config", Component: "system/config/index", Icon: "setting", Sort: 6, Type: 2, Permission: "system:config:list", Status: 1},
		{ParentID: sysMgmt.ID, Name: "存储管理", Path: "/system/storage", Component: "system/storage/index", Icon: "cloud-server", Sort: 7, Type: 2, Permission: "system:storage:list", Status: 1},
		{ParentID: sysMgmt.ID, Name: "文件管理", Path: "/system/file", Component: "system/file/index", Icon: "folder", Sort: 8, Type: 2, Permission: "system:file:list", Status: 1},
		{ParentID: monitor.ID, Name: "操作日志", Path: "/monitor/operation-log", Component: "monitor/operation-log/index", Icon: "file-text", Sort: 1, Type: 2, Permission: "monitor:operation-log:list", Status: 1},
		{ParentID: monitor.ID, Name: "登录日志", Path: "/monitor/login-log", Component: "monitor/login-log/index", Icon: "login", Sort: 2, Type: 2, Permission: "monitor:login-log:list", Status: 1},
		{ParentID: monitor.ID, Name: "慢查询日志", Path: "/monitor/show-log", Component: "monitor/show-log/index", Icon: "login", Sort: 2, Type: 2, Permission: "monitor:show-log:list", Status: 1},
	}
	global.DB.Create(&menus)

	// 分配菜单
	global.DB.Model(&adminRole).Association("Menus").Replace(append(menus, sysMgmt, monitor))

	// 默认API
	apis := []model.SysApi{
		// 用户管理
		{Path: "/api/v1/users", Method: "GET", Group: "用户管理", Description: "用户列表"},
		{Path: "/api/v1/users/:id", Method: "GET", Group: "用户管理", Description: "用户详情"},
		{Path: "/api/v1/users", Method: "POST", Group: "用户管理", Description: "创建用户"},
		{Path: "/api/v1/users/:id", Method: "PUT", Group: "用户管理", Description: "更新用户"},
		{Path: "/api/v1/users/:id", Method: "DELETE", Group: "用户管理", Description: "删除用户"},
		{Path: "/api/v1/users/:id/status", Method: "PUT", Group: "用户管理", Description: "修改用户状态"},
		{Path: "/api/v1/users/:id/password", Method: "PUT", Group: "用户管理", Description: "重置密码"},
		// 角色管理
		{Path: "/api/v1/roles", Method: "GET", Group: "角色管理", Description: "角色列表"},
		{Path: "/api/v1/roles/:id", Method: "GET", Group: "角色管理", Description: "角色详情"},
		{Path: "/api/v1/roles", Method: "POST", Group: "角色管理", Description: "创建角色"},
		{Path: "/api/v1/roles/:id", Method: "PUT", Group: "角色管理", Description: "更新角色"},
		{Path: "/api/v1/roles/:id", Method: "DELETE", Group: "角色管理", Description: "删除角色"},
		{Path: "/api/v1/roles/:id/menus", Method: "PUT", Group: "角色管理", Description: "分配菜单"},
		{Path: "/api/v1/roles/:id/apis", Method: "PUT", Group: "角色管理", Description: "分配API"},
		// 部门管理
		{Path: "/api/v1/depts/tree", Method: "GET", Group: "部门管理", Description: "部门树"},
		{Path: "/api/v1/depts/:id", Method: "GET", Group: "部门管理", Description: "部门详情"},
		{Path: "/api/v1/depts", Method: "POST", Group: "部门管理", Description: "创建部门"},
		{Path: "/api/v1/depts/:id", Method: "PUT", Group: "部门管理", Description: "更新部门"},
		{Path: "/api/v1/depts/:id", Method: "DELETE", Group: "部门管理", Description: "删除部门"},
		// 菜单管理
		{Path: "/api/v1/menus", Method: "GET", Group: "菜单管理", Description: "菜单列表"},
		{Path: "/api/v1/menus/tree-with-apis", Method: "GET", Group: "菜单管理", Description: "菜单树(带API)"},
		{Path: "/api/v1/menus/:id", Method: "GET", Group: "菜单管理", Description: "菜单详情"},
		{Path: "/api/v1/menus/:id/apis", Method: "GET", Group: "菜单管理", Description: "菜单API列表"},
		{Path: "/api/v1/menus/:id/apis", Method: "PUT", Group: "菜单管理", Description: "更新菜单API"},
		{Path: "/api/v1/menus", Method: "POST", Group: "菜单管理", Description: "创建菜单"},
		{Path: "/api/v1/menus/:id", Method: "PUT", Group: "菜单管理", Description: "更新菜单"},
		{Path: "/api/v1/menus/:id", Method: "DELETE", Group: "菜单管理", Description: "删除菜单"},
		// API管理
		{Path: "/api/v1/apis", Method: "GET", Group: "API管理", Description: "API列表"},
		{Path: "/api/v1/apis/:id", Method: "GET", Group: "API管理", Description: "API详情"},
		{Path: "/api/v1/apis", Method: "POST", Group: "API管理", Description: "创建API"},
		{Path: "/api/v1/apis/:id", Method: "PUT", Group: "API管理", Description: "更新API"},
		{Path: "/api/v1/apis/:id", Method: "DELETE", Group: "API管理", Description: "删除API"},
		// 日志管理
		{Path: "/api/v1/logs/operation", Method: "GET", Group: "日志管理", Description: "操作日志列表"},
		{Path: "/api/v1/logs/login", Method: "GET", Group: "日志管理", Description: "登录日志列表"},
		{Path: "/api/v1/logs/slow", Method: "GET", Group: "日志管理", Description: "慢查询日志列表"},
		// 配置管理
		{Path: "/api/v1/configs", Method: "GET", Group: "配置管理", Description: "配置列表"},
		{Path: "/api/v1/configs/key/:key", Method: "GET", Group: "配置管理", Description: "根据key获取配置"},
		{Path: "/api/v1/configs/keys", Method: "POST", Group: "配置管理", Description: "批量获取配置"},
		{Path: "/api/v1/configs", Method: "POST", Group: "配置管理", Description: "创建配置"},
		{Path: "/api/v1/configs/:id", Method: "PUT", Group: "配置管理", Description: "更新配置"},
		{Path: "/api/v1/configs/batch", Method: "PUT", Group: "配置管理", Description: "批量更新配置"},
		{Path: "/api/v1/configs/:id", Method: "DELETE", Group: "配置管理", Description: "删除配置"},
		// 存储管理
		{Path: "/api/v1/storages", Method: "GET", Group: "存储管理", Description: "存储配置列表"},
		{Path: "/api/v1/storages/:id", Method: "GET", Group: "存储管理", Description: "存储配置详情"},
		{Path: "/api/v1/storages", Method: "POST", Group: "存储管理", Description: "创建存储配置"},
		{Path: "/api/v1/storages/:id", Method: "PUT", Group: "存储管理", Description: "更新存储配置"},
		{Path: "/api/v1/storages/:id", Method: "DELETE", Group: "存储管理", Description: "删除存储配置"},
		{Path: "/api/v1/storages/:id/default", Method: "PUT", Group: "存储管理", Description: "设置默认存储"},
		{Path: "/api/v1/storages/test", Method: "POST", Group: "存储管理", Description: "测试存储配置"},
		// 文件管理
		{Path: "/api/v1/files", Method: "GET", Group: "文件管理", Description: "文件列表"},
		{Path: "/api/v1/files/:id", Method: "GET", Group: "文件管理", Description: "文件详情"},
		{Path: "/api/v1/files/:id", Method: "DELETE", Group: "文件管理", Description: "删除文件"},
		{Path: "/api/v1/files/credential", Method: "POST", Group: "文件管理", Description: "获取上传凭证"},
		{Path: "/api/v1/files/check-md5", Method: "POST", Group: "文件管理", Description: "MD5秒传检查"},
		{Path: "/api/v1/files/save", Method: "POST", Group: "文件管理", Description: "保存上传文件"},
		{Path: "/api/v1/files/multipart/init", Method: "POST", Group: "文件管理", Description: "初始化分片上传"},
		{Path: "/api/v1/files/multipart/parts", Method: "GET", Group: "文件管理", Description: "获取已上传分片"},
		{Path: "/api/v1/files/multipart/complete", Method: "POST", Group: "文件管理", Description: "完成分片上传"},
		{Path: "/api/v1/files/multipart/abort", Method: "POST", Group: "文件管理", Description: "取消分片上传"},
		{Path: "/api/v1/files/upload/local", Method: "POST", Group: "文件管理", Description: "本地文件上传"},
		{Path: "/api/v1/files/upload/chunk", Method: "POST", Group: "文件管理", Description: "上传分片"},
	}
	global.DB.Create(&apis)

	// 角色关联API
	global.DB.Model(&adminRole).Association("Apis").Replace(apis)

	// 同步Casbin策略：为admin角色授予所有API权限
	if global.Enforcer != nil {
		_, _ = global.Enforcer.RemoveFilteredPolicy(0, "admin")
		policies := make([][]string, 0, len(apis))
		for _, api := range apis {
			policies = append(policies, []string{"admin", api.Path, api.Method})
		}
		if len(policies) > 0 {
			_, _ = global.Enforcer.AddPolicies(policies)
			_ = global.Enforcer.SavePolicy()
		}
	}

	global.Log.Info("默认数据初始化成功")
}

func initDefaultConfigs() {
	// 若已有配置数据则跳过
	var count int64
	global.DB.Model(&model.SysConfig{}).Count(&count)
	if count > 0 {
		return
	}

	// 默认系统配置
	configs := []model.SysConfig{
		{Name: "系统名称", Key: "sys_name", Value: "Go RBAC Admin", ValueType: "string", Remark: "显示在侧边栏顶部"},
		{Name: "系统Logo", Key: "sys_logo", Value: "/src/assets/logo.svg", ValueType: "string", Remark: "系统Logo图片地址"},
		{Name: "菜单背景色", Key: "menu_bg_color", Value: "#001529", ValueType: "string", Remark: "侧边栏菜单背景色"},
		{Name: "菜单文字颜色", Key: "menu_text_color", Value: "rgba(255, 255, 255, 0.65)", ValueType: "string", Remark: "菜单文字颜色"},
		{Name: "菜单激活文字颜色", Key: "menu_active_text_color", Value: "#ffffff", ValueType: "string", Remark: "菜单选中时文字颜色"},
		{Name: "菜单激活背景色", Key: "menu_active_bg_color", Value: "#1890ff", ValueType: "string", Remark: "菜单选中时背景色"},
		{Name: "头部背景色", Key: "header_bg_color", Value: "#ffffff", ValueType: "string", Remark: "顶部栏背景色"},
		{Name: "头部文字颜色", Key: "header_text_color", Value: "#333333", ValueType: "string", Remark: "顶部栏文字颜色"},
		{Name: "AI配置", Key: "ai_config", Value: `{"default_provider":"阿里云百炼","providers":[{"name":"阿里云百炼","api_key":"","base_url":"https://dashscope.aliyuncs.com/compatible-mode/v1","models":[{"id":"deepseek-v3.2","name":"DeepSeek-V3.2","description":"DeepSeek最新模型,支持联网和思考"},{"id":"qwen3-max","name":"通义千问3-Max","description":"通义千问3系列Max模型"}]}]}`, ValueType: "json", Remark: "AI平台配置，包含平台名称、API Key、基础URL和模型列表"},
		{Name: "前台模式", Key: "front_mode", Value: "full", ValueType: "string", Remark: "前台模式: full=完整前台, profile=仅个人中心(用于身份认证)"},
		{Name: "用户身份按钮显示", Key: "user_profile_button_visible", Value: "false", ValueType: "string", Remark: "后台用户管理列表是否显示身份按钮"},
		{Name: "部门模块显示", Key: "dept_module_enabled", Value: "true", ValueType: "string", Remark: "后台菜单中是否显示部门管理模块"},
	}
	global.DB.Create(&configs)
	global.Log.Info("系统配置初始化成功")
}

func ensureBuiltInData() {
	ensureConfigExists(model.SysConfig{
		Name:      "用户身份按钮显示",
		Key:       "user_profile_button_visible",
		Value:     "false",
		ValueType: "string",
		Remark:    "后台用户管理列表是否显示身份按钮",
	})
	ensureConfigExists(model.SysConfig{
		Name:      "部门模块显示",
		Key:       "dept_module_enabled",
		Value:     "true",
		ValueType: "string",
		Remark:    "后台菜单中是否显示部门管理模块",
	})

	rootDept := ensureRootDeptExists()
	backfillDepartmentFoundation(rootDept.ID)

	ensureDeptApiAccess()
	ensureDeptMenus()

	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/users/batch-status",
		Method:      "PUT",
		Group:       "用户管理",
		Description: "批量修改用户状态",
		NeedAuth:    true,
	}, "/api/v1/users/:id/status", "PUT")

	ensureUserBatchStatusMenus()
}

func ensureConfigExists(config model.SysConfig) {
	var existing model.SysConfig
	err := global.DB.Where("`key` = ?", config.Key).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := global.DB.Create(&config).Error; err != nil {
			global.Log.Errorf("补齐系统配置失败(%s): %v", config.Key, err)
		}
	}
}

func ensureRootDeptExists() model.SysDept {
	rootDept := model.SysDept{
		ParentID:  0,
		Ancestors: "0",
		Name:      "平台",
		Sort:      1,
		Status:    1,
		Remark:    "系统根部门",
	}

	if err := global.DB.
		Where("parent_id = ? AND name = ?", 0, rootDept.Name).
		Assign(model.SysDept{
			Ancestors: rootDept.Ancestors,
			Sort:      rootDept.Sort,
			Status:    rootDept.Status,
			Remark:    rootDept.Remark,
		}).
		FirstOrCreate(&rootDept).Error; err != nil {
		global.Log.Errorf("补齐根部门失败: %v", err)
	}

	return rootDept
}

func backfillDepartmentFoundation(rootDeptID uint) {
	if rootDeptID == 0 {
		return
	}

	if err := global.DB.Model(&model.SysUser{}).
		Where("dept_id = 0 OR dept_id IS NULL").
		Update("dept_id", rootDeptID).Error; err != nil {
		global.Log.Errorf("回填用户部门失败: %v", err)
	}

	if err := global.DB.Model(&model.SysRole{}).
		Where("data_scope = 0 OR data_scope IS NULL").
		Update("data_scope", model.DataScopeAll).Error; err != nil {
		global.Log.Errorf("回填角色数据范围失败: %v", err)
	}
}

func ensureDeptApiAccess() {
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/depts/tree",
		Method:      "GET",
		Group:       "部门管理",
		Description: "部门树",
		NeedAuth:    true,
	}, "/api/v1/menus", "GET")
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/depts/manageable-tree",
		Method:      "GET",
		Group:       "部门管理",
		Description: "可管理部门树",
		NeedAuth:    true,
	}, "/api/v1/depts/tree", "GET")
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/depts/:id",
		Method:      "GET",
		Group:       "部门管理",
		Description: "部门详情",
		NeedAuth:    true,
	}, "/api/v1/menus/:id", "GET")
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/depts",
		Method:      "POST",
		Group:       "部门管理",
		Description: "创建部门",
		NeedAuth:    true,
	}, "/api/v1/menus", "POST")
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/depts/:id",
		Method:      "PUT",
		Group:       "部门管理",
		Description: "更新部门",
		NeedAuth:    true,
	}, "/api/v1/menus/:id", "PUT")
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/depts/:id",
		Method:      "DELETE",
		Group:       "部门管理",
		Description: "删除部门",
		NeedAuth:    true,
	}, "/api/v1/menus/:id", "DELETE")
}

func ensureDeptMenus() {
	var systemMenu model.SysMenu
	if err := global.DB.Where("path = ? AND type = ?", "/system", 1).First(&systemMenu).Error; err != nil {
		global.Log.Errorf("查询系统管理目录失败: %v", err)
		return
	}

	deptMenu := model.SysMenu{
		ParentID:   systemMenu.ID,
		Name:       "部门管理",
		Path:       "/system/dept",
		Component:  "system/dept/index",
		Icon:       "apartment",
		Sort:       3,
		Type:       2,
		Permission: "system:dept:list",
		Status:     1,
		Hidden:     0,
	}

	if err := global.DB.
		Where("permission = ?", deptMenu.Permission).
		Attrs(model.SysMenu{
			ParentID:  deptMenu.ParentID,
			Name:      deptMenu.Name,
			Path:      deptMenu.Path,
			Component: deptMenu.Component,
			Icon:      deptMenu.Icon,
			Sort:      deptMenu.Sort,
			Type:      deptMenu.Type,
			Status:    deptMenu.Status,
			Hidden:    deptMenu.Hidden,
		}).
		FirstOrCreate(&deptMenu).Error; err != nil {
		global.Log.Errorf("补齐部门管理菜单失败: %v", err)
		return
	}

	buttonDefinitions := []model.SysMenu{
		{ParentID: deptMenu.ID, Name: "新增", Sort: 1, Type: 3, Permission: "system:dept:add", Status: 1},
		{ParentID: deptMenu.ID, Name: "编辑", Sort: 2, Type: 3, Permission: "system:dept:edit", Status: 1},
		{ParentID: deptMenu.ID, Name: "删除", Sort: 3, Type: 3, Permission: "system:dept:delete", Status: 1},
	}

	menuIDs := []uint{deptMenu.ID}
	for _, definition := range buttonDefinitions {
		menu := definition
		if err := global.DB.
			Where("permission = ?", menu.Permission).
			Attrs(model.SysMenu{
				ParentID: menu.ParentID,
				Name:     menu.Name,
				Sort:     menu.Sort,
				Type:     menu.Type,
				Status:   menu.Status,
				Hidden:   0,
			}).
			FirstOrCreate(&menu).Error; err != nil {
			global.Log.Errorf("补齐部门按钮权限失败(%s): %v", definition.Permission, err)
			continue
		}
		menuIDs = append(menuIDs, menu.ID)
	}

	grantMenusToRoleCodes(menuIDs, []string{"admin", "system_admin"})
}

func ensureApiAccessInheritedFrom(api model.SysApi, sourcePath, sourceMethod string) {
	if err := global.DB.
		Where("path = ? AND method = ?", api.Path, api.Method).
		Assign(model.SysApi{
			Group:       api.Group,
			Description: api.Description,
			NeedAuth:    api.NeedAuth,
		}).
		FirstOrCreate(&api).Error; err != nil {
		global.Log.Errorf("补齐系统API失败(%s %s): %v", api.Method, api.Path, err)
		return
	}

	var roles []model.SysRole
	if err := global.DB.Table("sys_role AS sr").
		Select("sr.*").
		Joins("JOIN sys_role_api AS sra ON sra.sys_role_id = sr.id").
		Joins("JOIN sys_api AS sa ON sa.id = sra.sys_api_id").
		Where("sa.path = ? AND sa.method = ?", sourcePath, sourceMethod).
		Group("sr.id").
		Find(&roles).Error; err != nil {
		global.Log.Errorf("查询源API授权角色失败(%s %s): %v", sourceMethod, sourcePath, err)
		return
	}
	if len(roles) == 0 {
		if err := global.DB.Where("code IN ?", []string{"admin", "system_admin"}).Find(&roles).Error; err != nil {
			global.Log.Errorf("查询内置角色失败: %v", err)
			return
		}
	}

	policyChanged := false
	for _, role := range roles {
		var count int64
		if err := global.DB.Table("sys_role_api").
			Where("sys_role_id = ? AND sys_api_id = ?", role.ID, api.ID).
			Count(&count).Error; err == nil && count == 0 {
			if err := global.DB.Model(&role).Association("Apis").Append(&api); err != nil {
				global.Log.Errorf("关联角色API失败(%s): %v", role.Code, err)
			}
		}

		if global.Enforcer != nil {
			if ok, err := global.Enforcer.AddPolicy(role.Code, api.Path, api.Method); err != nil {
				global.Log.Errorf("补齐Casbin策略失败(%s %s %s): %v", role.Code, api.Method, api.Path, err)
			} else if ok {
				policyChanged = true
			}
		}
	}

	if policyChanged && global.Enforcer != nil {
		_ = global.Enforcer.SavePolicy()
	}
}

func ensureUserBatchStatusMenus() {
	var userMenu model.SysMenu
	if err := global.DB.Where("permission = ? AND type = ?", "system:user:list", 2).First(&userMenu).Error; err != nil {
		global.Log.Errorf("查询用户管理菜单失败: %v", err)
		return
	}

	menuDefinitions := []model.SysMenu{
		{
			ParentID:   userMenu.ID,
			Name:       "批量启用",
			Path:       "",
			Component:  "",
			Icon:       "",
			Sort:       5,
			Type:       3,
			Permission: "system:user:batchEnable",
			Status:     1,
			Hidden:     0,
		},
		{
			ParentID:   userMenu.ID,
			Name:       "批量禁用",
			Path:       "",
			Component:  "",
			Icon:       "",
			Sort:       6,
			Type:       3,
			Permission: "system:user:batchDisable",
			Status:     1,
			Hidden:     0,
		},
	}

	for _, definition := range menuDefinitions {
		menu := definition
		if err := global.DB.
			Where("permission = ?", menu.Permission).
			Attrs(model.SysMenu{
				ParentID:  menu.ParentID,
				Name:      menu.Name,
				Path:      menu.Path,
				Component: menu.Component,
				Icon:      menu.Icon,
				Sort:      menu.Sort,
				Type:      menu.Type,
				Status:    menu.Status,
				Hidden:    menu.Hidden,
			}).
			FirstOrCreate(&menu).Error; err != nil {
			global.Log.Errorf("补齐用户批量状态按钮失败(%s): %v", definition.Permission, err)
			continue
		}

		grantMenuToRolesWithPermission(menu.ID, "system:user:edit")
	}
}

func grantMenuToRolesWithPermission(menuID uint, sourcePermission string) {
	var sourceMenus []model.SysMenu
	if err := global.DB.Where("permission = ?", sourcePermission).Find(&sourceMenus).Error; err != nil {
		global.Log.Errorf("查询源权限菜单失败(%s): %v", sourcePermission, err)
		return
	}
	if len(sourceMenus) == 0 {
		return
	}

	sourceMenuIDs := make([]uint, 0, len(sourceMenus))
	for _, menu := range sourceMenus {
		sourceMenuIDs = append(sourceMenuIDs, menu.ID)
	}

	var roleIDs []uint
	if err := global.DB.Table("sys_role_menu").
		Distinct("sys_role_id").
		Where("sys_menu_id IN ?", sourceMenuIDs).
		Pluck("sys_role_id", &roleIDs).Error; err != nil {
		global.Log.Errorf("查询源权限角色失败(%s): %v", sourcePermission, err)
		return
	}

	for _, roleID := range roleIDs {
		var count int64
		if err := global.DB.Table("sys_role_menu").
			Where("sys_role_id = ? AND sys_menu_id = ?", roleID, menuID).
			Count(&count).Error; err == nil && count == 0 {
			if err := global.DB.Exec(
				"INSERT INTO sys_role_menu (sys_role_id, sys_menu_id) VALUES (?, ?)",
				roleID, menuID,
			).Error; err != nil {
				global.Log.Errorf("补齐角色菜单权限失败(role=%d, menu=%d): %v", roleID, menuID, err)
			}
		}
	}
}

func grantMenusToRoleCodes(menuIDs []uint, roleCodes []string) {
	if len(menuIDs) == 0 || len(roleCodes) == 0 {
		return
	}

	var roles []model.SysRole
	if err := global.DB.Where("code IN ?", roleCodes).Find(&roles).Error; err != nil {
		global.Log.Errorf("查询内置角色失败: %v", err)
		return
	}

	for _, role := range roles {
		for _, menuID := range menuIDs {
			var count int64
			if err := global.DB.Table("sys_role_menu").
				Where("sys_role_id = ? AND sys_menu_id = ?", role.ID, menuID).
				Count(&count).Error; err == nil && count == 0 {
				if err := global.DB.Exec(
					"INSERT INTO sys_role_menu (sys_role_id, sys_menu_id) VALUES (?, ?)",
					role.ID, menuID,
				).Error; err != nil {
					global.Log.Errorf("补齐角色菜单权限失败(role=%d, menu=%d): %v", role.ID, menuID, err)
				}
			}
		}
	}
}
