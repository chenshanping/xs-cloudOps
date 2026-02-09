package initialize

import (
	"go-base-server/global"
	"go-base-server/model"
	"go-base-server/utils"
)

func InitDBTables() {
	err := global.DB.AutoMigrate(
		&model.SysUser{},
		&model.SysRole{},
		&model.SysMenu{},
		&model.SysApi{},
		&model.SysOperationLog{},
		&model.SysLoginLog{},
		&model.SysSlowLog{},
		&model.SysConfig{},
		&model.SysStorage{},
		&model.SysFile{},
		&model.SysFileChunk{},
		&model.SysGenerator{},
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
		Name:   "超级管理员",
		Code:   "admin",
		Sort:   1,
		Status: 1,
		Remark: "拥有所有权限",
	}
	global.DB.Create(&adminRole)

	// 创建默认用户
	hashedPassword, _ := utils.HashPassword("123456")
	adminUser := model.SysUser{
		Username: "admin",
		Password: hashedPassword,
		Nickname: "管理员",
		Status:   1,
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
		{ParentID: sysMgmt.ID, Name: "菜单管理", Path: "/system/menu", Component: "system/menu/index", Icon: "menu", Sort: 3, Type: 2, Permission: "system:menu:list", Status: 1},
		{ParentID: sysMgmt.ID, Name: "API管理", Path: "/system/api", Component: "system/api/index", Icon: "api", Sort: 4, Type: 2, Permission: "system:api:list", Status: 1},
		{ParentID: sysMgmt.ID, Name: "参数配置", Path: "/system/config", Component: "system/config/index", Icon: "setting", Sort: 5, Type: 2, Permission: "system:config:list", Status: 1},
		{ParentID: sysMgmt.ID, Name: "存储管理", Path: "/system/storage", Component: "system/storage/index", Icon: "cloud-server", Sort: 6, Type: 2, Permission: "system:storage:list", Status: 1},
		{ParentID: sysMgmt.ID, Name: "文件管理", Path: "/system/file", Component: "system/file/index", Icon: "folder", Sort: 7, Type: 2, Permission: "system:file:list", Status: 1},
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
	}
	global.DB.Create(&configs)
	global.Log.Info("系统配置初始化成功")
}
