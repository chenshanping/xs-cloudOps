package initialize

import (
	"encoding/json"
	"errors"
	"fmt"
	appconfig "server/config"
	"server/global"
	"server/model"
	"server/service"
	"server/utils"
	"strings"

	"gorm.io/gorm"
)

func InitDBTables() {
	err := global.DB.AutoMigrate(
		&model.SysUser{},
		&model.SysRole{},
		&model.SysRoleDataScope{},
		&model.SysDept{},
		&model.SysMenu{},
		&model.SysApi{},
		&model.SysOperationLog{},
		&model.SysLoginLog{},
		&model.SysConfig{},
		&model.AIProviderConfig{},
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
		Name:         "超级管理员",
		Code:         "admin",
		Sort:         1,
		Status:       1,
		IsSuperAdmin: true,
		DataScope:    model.DataScopeAll,
		Remark:       "拥有所有权限",
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
		{ParentID: sysMgmt.ID, Name: "文件管理", Path: "/system/file", Component: "system/file/index", Icon: "folder", Sort: 7, Type: 2, Permission: "system:file:list", Status: 1},
		{ParentID: monitor.ID, Name: "操作日志", Path: "/monitor/operation-log", Component: "monitor/operation-log/index", Icon: "file-text", Sort: 1, Type: 2, Permission: "monitor:operation-log:list", Status: 1},
		{ParentID: monitor.ID, Name: "登录日志", Path: "/monitor/login-log", Component: "monitor/login-log/index", Icon: "login", Sort: 2, Type: 2, Permission: "monitor:login-log:list", Status: 1},
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
		{Path: "/api/v1/roles/data-scope-resources", Method: "GET", Group: "角色管理", Description: "数据权限资源列表"},
		{Path: "/api/v1/roles/:id/menus", Method: "PUT", Group: "角色管理", Description: "分配菜单"},
		{Path: "/api/v1/roles/:id/apis", Method: "PUT", Group: "角色管理", Description: "分配API"},
		{Path: "/api/v1/roles/:id/data-scopes", Method: "PUT", Group: "角色管理", Description: "分配数据权限"},
		{Path: "/api/v1/roles/:id/permissions", Method: "PUT", Group: "角色管理", Description: "统一保存角色权限"},
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
		// 配置管理
		{Path: "/api/v1/configs", Method: "GET", Group: "配置管理", Description: "配置列表"},
		{Path: "/api/v1/configs/key/:key", Method: "GET", Group: "配置管理", Description: "根据key获取配置"},
		{Path: "/api/v1/configs/keys", Method: "POST", Group: "配置管理", Description: "批量获取配置"},
		{Path: "/api/v1/configs", Method: "POST", Group: "配置管理", Description: "创建配置"},
		{Path: "/api/v1/configs/:id", Method: "PUT", Group: "配置管理", Description: "更新配置"},
		{Path: "/api/v1/configs/batch", Method: "PUT", Group: "配置管理", Description: "批量更新配置"},
		{Path: "/api/v1/configs/:id", Method: "DELETE", Group: "配置管理", Description: "删除配置"},
		{Path: "/api/v1/configs/storage/test", Method: "POST", Group: "配置管理", Description: "测试存储配置"},
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
		{Name: "系统Logo文件ID", Key: service.SysLogoFileIDConfigKey, Value: "", ValueType: "string", Remark: "系统Logo关联文件ID"},
		{Name: "注册默认头像文件ID", Key: service.RegisterLogoFileIDConfigKey, Value: "", ValueType: "string", Remark: "注册默认头像关联文件ID"},
		{Name: "登录页背景图文件ID", Key: service.LoginBGImageFileIDConfigKey, Value: "", ValueType: "string", Remark: "登录页背景图关联文件ID"},
		{Name: "公开配置键", Key: service.PublicConfigKeysConfigKey, Value: service.DefaultPublicConfigKeysValue(), ValueType: "json", Remark: "允许匿名批量读取的配置键(JSON数组)，敏感键即使写入也不会公开"},
		{Name: "前台模式", Key: "front_mode", Value: "full", ValueType: "string", Remark: "前台模式: full=完整前台, profile=仅个人中心(用于身份认证)"},
		{Name: "用户身份按钮显示", Key: "user_profile_button_visible", Value: "false", ValueType: "string", Remark: "后台用户管理列表是否显示身份按钮"},
		{Name: "用户默认密码", Key: "user_default_password", Value: "123456", ValueType: "string", Remark: "后台用户管理单条/批量重置密码默认值"},
		{Name: "文件删除方式", Key: service.FileDeleteModeConfigKey, Value: service.FileDeleteModeLogical, ValueType: "string", Remark: "文件删除方式: logical=逻辑删除, physical=物理删除"},
		{Name: "存储类型", Key: service.StorageTypeConfigKey, Value: string(service.Storage.DefaultStorageType()), ValueType: "string", Remark: "当前文件上传使用的存储类型"},
	}
	for _, storageType := range service.Storage.SupportedStorageTypes() {
		configs = append(configs, model.SysConfig{
			Name:      storageConfigName(storageType),
			Key:       service.StorageConfigKey(storageType),
			Value:     service.Storage.DefaultStorageConfig(storageType),
			ValueType: "json",
			Remark:    fmt.Sprintf("%s的已保存配置(JSON)", storageConfigLabel(storageType)),
		})
	}
	global.DB.Create(&configs)
	global.Log.Info("系统配置初始化成功")
}

func ensureBuiltInData() {
	ensureSystemStorageConfigs()
	ensureFileStorageSnapshots()
	ensureGenderDictData()
	ensureAIProvidersExist()
	ensureConfigExists(model.SysConfig{
		Name:      "公开配置键",
		Key:       service.PublicConfigKeysConfigKey,
		Value:     service.DefaultPublicConfigKeysValue(),
		ValueType: "json",
		Remark:    "允许匿名批量读取的配置键(JSON数组)，敏感键即使写入也不会公开",
	})
	ensureConfigExists(model.SysConfig{
		Name:      "系统Logo文件ID",
		Key:       service.SysLogoFileIDConfigKey,
		Value:     "",
		ValueType: "string",
		Remark:    "系统Logo关联文件ID",
	})
	ensureConfigExists(model.SysConfig{
		Name:      "注册默认头像文件ID",
		Key:       service.RegisterLogoFileIDConfigKey,
		Value:     "",
		ValueType: "string",
		Remark:    "注册默认头像关联文件ID",
	})
	ensureConfigExists(model.SysConfig{
		Name:      "登录页背景图文件ID",
		Key:       service.LoginBGImageFileIDConfigKey,
		Value:     "",
		ValueType: "string",
		Remark:    "登录页背景图关联文件ID",
	})
	ensureConfigExists(model.SysConfig{
		Name:      "用户身份按钮显示",
		Key:       "user_profile_button_visible",
		Value:     "false",
		ValueType: "string",
		Remark:    "后台用户管理列表是否显示身份按钮",
	})
	ensureConfigExists(model.SysConfig{
		Name:      "用户默认密码",
		Key:       "user_default_password",
		Value:     "123456",
		ValueType: "string",
		Remark:    "后台用户管理单条/批量重置密码默认值",
	})
	ensureConfigExists(model.SysConfig{
		Name:      "文件删除方式",
		Key:       service.FileDeleteModeConfigKey,
		Value:     service.FileDeleteModeLogical,
		ValueType: "string",
		Remark:    "文件删除方式: logical=逻辑删除, physical=物理删除",
	})

	rootDept := ensureRootDeptExists()
	backfillDepartmentFoundation(rootDept.ID)
	backfillExplicitSuperAdminRoles()

	ensureDeptApiAccess()
	ensureAIApiAccess()
	ensureAIToolMenus()
	ensureDeptMenus()
	ensureLogAuditMenus()
	cleanupSlowLogBuiltInData()
	cleanupStorageBuiltInData()

	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/users/batch-status",
		Method:      "PUT",
		Group:       "用户管理",
		Description: "批量修改用户状态",
		NeedAuth:    true,
	}, "/api/v1/users/:id/status", "PUT")
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/users/default-password",
		Method:      "GET",
		Group:       "用户管理",
		Description: "默认密码",
		NeedAuth:    true,
	}, "/api/v1/users", "POST")
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/users/batch-password",
		Method:      "PUT",
		Group:       "用户管理",
		Description: "批量重置密码",
		NeedAuth:    true,
	}, "/api/v1/users/:id/password", "PUT")
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/roles/data-scope-resources",
		Method:      "GET",
		Group:       "角色管理",
		Description: "数据权限资源列表",
		NeedAuth:    true,
	}, "/api/v1/roles/:id/apis", "PUT")
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/roles/:id/data-scopes",
		Method:      "PUT",
		Group:       "角色管理",
		Description: "分配数据权限",
		NeedAuth:    true,
	}, "/api/v1/roles/:id/apis", "PUT")
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/roles/:id/permissions",
		Method:      "PUT",
		Group:       "角色管理",
		Description: "统一保存角色权限",
		NeedAuth:    true,
	}, "/api/v1/roles/:id/apis", "PUT")
	ensureUserOperationMenus()
	ensureUserImportExportMenus()
}

func ensureGenderDictData() {
	dictType := model.SysDictType{
		Name:   "性别",
		Type:   "sys_gender",
		Status: 1,
		Remark: "用户性别字典",
	}

	if err := global.DB.
		Where("type = ?", dictType.Type).
		Attrs(model.SysDictType{
			Name:   dictType.Name,
			Status: dictType.Status,
			Remark: dictType.Remark,
		}).
		FirstOrCreate(&dictType).Error; err != nil {
		global.Log.Errorf("补齐性别字典类型失败: %v", err)
		return
	}

	items := []model.SysDictData{
		{DictType: "sys_gender", Label: "男", Value: "0", Sort: 1, Status: 1, TagType: "processing", IsDefault: 0, Remark: ""},
		{DictType: "sys_gender", Label: "女", Value: "1", Sort: 2, Status: 1, TagType: "pink", IsDefault: 0, Remark: ""},
	}

	for _, definition := range items {
		item := definition
		if err := global.DB.
			Where("dict_type = ? AND value = ?", item.DictType, item.Value).
			Attrs(model.SysDictData{
				Label:     item.Label,
				Sort:      item.Sort,
				Status:    item.Status,
				TagType:   item.TagType,
				IsDefault: item.IsDefault,
				Remark:    item.Remark,
			}).
			FirstOrCreate(&item).Error; err != nil {
			global.Log.Errorf("补齐性别字典数据失败(%s): %v", definition.Value, err)
		}
	}

	if err := service.Cache.ClearDictCache(dictType.Type); err != nil {
		global.Log.Errorf("清理性别字典缓存失败: %v", err)
	}
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

func ensureAIProvidersExist() {
	var count int64
	if err := global.DB.Model(&model.AIProviderConfig{}).Count(&count).Error; err != nil {
		global.Log.Errorf("统计 AI 平台配置失败: %v", err)
		return
	}
	if count > 0 {
		return
	}

	if cfg, found, err := loadAIConfigFromLegacySysConfig(); err != nil {
		global.Log.Errorf("读取历史 AI 配置失败: %v", err)
		return
	} else if found {
		if err := createAIProvidersFromConfig(cfg); err != nil {
			global.Log.Errorf("迁移历史 sys_config.ai_config 到 ai_providers 失败: %v", err)
		}
		return
	}

	if err := createAIProvidersFromConfig(defaultAIConfig()); err != nil {
		global.Log.Errorf("补齐 AI 平台配置失败: %v", err)
	}
}

func loadAIConfigFromLegacySysConfig() (appconfig.AI, bool, error) {
	var cfg appconfig.AI
	var legacy model.SysConfig
	err := global.DB.Where("`key` = ?", "ai_config").First(&legacy).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return cfg, false, nil
	}
	if err != nil {
		return cfg, false, err
	}
	if strings.TrimSpace(legacy.Value) == "" {
		return cfg, false, nil
	}

	if err := json.Unmarshal([]byte(legacy.Value), &cfg); err != nil {
		return cfg, false, err
	}
	return cfg, true, nil
}

func createAIProvidersFromConfig(cfg appconfig.AI) error {
	rows, err := buildAIProviderRows(cfg)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	return global.DB.Create(&rows).Error
}

func buildAIProviderRows(cfg appconfig.AI) ([]model.AIProviderConfig, error) {
	if len(cfg.Providers) == 0 {
		return nil, nil
	}

	defaultName := strings.TrimSpace(cfg.DefaultProvider)
	if defaultName == "" {
		defaultName = cfg.Providers[0].Name
	}

	hasDefault := false
	for _, provider := range cfg.Providers {
		if provider.Name == defaultName {
			hasDefault = true
			break
		}
	}
	if !hasDefault {
		defaultName = cfg.Providers[0].Name
	}

	rows := make([]model.AIProviderConfig, 0, len(cfg.Providers))
	for index, provider := range cfg.Providers {
		modelsJSON, err := json.Marshal(provider.Models)
		if err != nil {
			return nil, err
		}
		rows = append(rows, model.AIProviderConfig{
			Name:       provider.Name,
			APIKey:     provider.APIKey,
			BaseURL:    provider.BaseURL,
			ModelsJSON: string(modelsJSON),
			IsDefault:  provider.Name == defaultName,
			Sort:       index,
		})
	}

	return rows, nil
}

func defaultAIConfig() appconfig.AI {
	return appconfig.AI{
		DefaultProvider: "阿里云百炼",
		Providers: []appconfig.AIProvider{
			{
				Name:    "阿里云百炼",
				APIKey:  "",
				BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
				Models: []appconfig.AIModel{
					{ID: "deepseek-v3.2", Name: "DeepSeek-V3.2", Description: "DeepSeek最新模型,支持联网和思考"},
					{ID: "qwen3-max", Name: "通义千问3-Max", Description: "通义千问3系列Max模型"},
				},
			},
		},
	}
}

func ensureSystemStorageConfigs() {
	keys := []string{service.StorageTypeConfigKey, service.LegacyStorageConfigConfigKey}
	for _, storageType := range service.Storage.SupportedStorageTypes() {
		keys = append(keys, service.StorageConfigKey(storageType))
	}

	configs, err := service.Config.GetConfigsByKeys(keys)
	if err != nil {
		global.Log.Errorf("查询系统存储配置失败: %v", err)
		return
	}

	storageType := service.Storage.DefaultStorageType()
	typeConfigured := false
	if typeConfig, ok := configs[service.StorageTypeConfigKey]; ok && strings.TrimSpace(typeConfig.Value) != "" {
		storageType = model.StorageType(typeConfig.Value)
		typeConfigured = true
	}

	typeConfigs := make(map[model.StorageType]string)
	for _, itemType := range service.Storage.SupportedStorageTypes() {
		typeConfigs[itemType] = service.Storage.DefaultStorageConfig(itemType)
		if config, ok := configs[service.StorageConfigKey(itemType)]; ok && strings.TrimSpace(config.Value) != "" {
			typeConfigs[itemType] = config.Value
		}
	}

	if global.DB.Migrator().HasTable((&model.LegacyStorageRecord{}).TableName()) {
		var legacyStorages []model.LegacyStorageRecord
		if err := global.DB.Where("status = ?", 1).Order("is_default DESC, id ASC").Find(&legacyStorages).Error; err == nil {
			for _, legacy := range legacyStorages {
				if current, ok := configs[service.StorageConfigKey(legacy.Type)]; !ok || strings.TrimSpace(current.Value) == "" {
					typeConfigs[legacy.Type] = legacy.Config
				}
				if !typeConfigured {
					storageType = legacy.Type
					typeConfigured = true
				}
			}
		}
	}

	if config, ok := configs[service.LegacyStorageConfigConfigKey]; ok && strings.TrimSpace(config.Value) != "" {
		if current, ok := configs[service.StorageConfigKey(storageType)]; !ok || strings.TrimSpace(current.Value) == "" {
			typeConfigs[storageType] = config.Value
		}
	}

	upsertConfigValue(model.SysConfig{
		Name:      "存储类型",
		Key:       service.StorageTypeConfigKey,
		Value:     string(storageType),
		ValueType: "string",
		Remark:    "当前文件上传使用的存储类型",
	})
	for _, itemType := range service.Storage.SupportedStorageTypes() {
		upsertConfigValue(model.SysConfig{
			Name:      storageConfigName(itemType),
			Key:       service.StorageConfigKey(itemType),
			Value:     typeConfigs[itemType],
			ValueType: "json",
			Remark:    fmt.Sprintf("%s的已保存配置(JSON)", storageConfigLabel(itemType)),
		})
	}
}

func ensureFileStorageSnapshots() {
	hasLegacyTable := global.DB.Migrator().HasTable((&model.LegacyStorageRecord{}).TableName())
	hasFileStorageID := global.DB.Migrator().HasColumn(&model.SysFile{}, "storage_id")
	hasChunkStorageID := global.DB.Migrator().HasColumn(&model.SysFileChunk{}, "storage_id")
	if hasLegacyTable && (hasFileStorageID || hasChunkStorageID) {
		var legacyStorages []model.LegacyStorageRecord
		if err := global.DB.Where("status = ?", 1).Find(&legacyStorages).Error; err != nil {
			global.Log.Errorf("查询历史存储配置失败: %v", err)
			return
		}

		for _, legacy := range legacyStorages {
			if hasFileStorageID {
				if err := global.DB.Model(&model.SysFile{}).
					Where("storage_id = ? AND (storage_type IS NULL OR storage_type = '')", legacy.ID).
					Update("storage_type", string(legacy.Type)).Error; err != nil {
					global.Log.Errorf("回填文件存储快照失败(storage_id=%d): %v", legacy.ID, err)
				}
			}
			if hasChunkStorageID {
				if err := global.DB.Model(&model.SysFileChunk{}).
					Where("storage_id = ? AND (storage_type IS NULL OR storage_type = '')", legacy.ID).
					Update("storage_type", string(legacy.Type)).Error; err != nil {
					global.Log.Errorf("回填分片存储快照失败(storage_id=%d): %v", legacy.ID, err)
				}
			}
		}
	}

	systemStorage, err := service.Storage.GetDefaultStorage()
	if err != nil {
		global.Log.Errorf("查询系统存储配置失败: %v", err)
		return
	}

	if err := global.DB.Model(&model.SysFile{}).
		Where("(storage_type IS NULL OR storage_type = '')").
		Update("storage_type", string(systemStorage.Type)).Error; err != nil {
		global.Log.Errorf("回填文件存储快照失败(默认配置): %v", err)
	}

	if err := global.DB.Model(&model.SysFileChunk{}).
		Where("(storage_type IS NULL OR storage_type = '')").
		Update("storage_type", string(systemStorage.Type)).Error; err != nil {
		global.Log.Errorf("回填分片存储快照失败(默认配置): %v", err)
	}
}

func storageConfigLabel(storageType model.StorageType) string {
	switch storageType {
	case model.StorageTypeLocal:
		return "本地存储"
	case model.StorageTypeAliyun:
		return "阿里云 OSS"
	case model.StorageTypeTencent:
		return "腾讯云 COS"
	case model.StorageTypeMinio:
		return "MinIO"
	default:
		return string(storageType)
	}
}

func storageConfigName(storageType model.StorageType) string {
	return fmt.Sprintf("%s配置", storageConfigLabel(storageType))
}

func upsertConfigValue(config model.SysConfig) {
	var existing model.SysConfig
	err := global.DB.Where("`key` = ?", config.Key).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := global.DB.Create(&config).Error; err != nil {
			global.Log.Errorf("补齐系统配置失败(%s): %v", config.Key, err)
		}
		return
	}
	if err != nil {
		global.Log.Errorf("查询系统配置失败(%s): %v", config.Key, err)
		return
	}
	if strings.TrimSpace(existing.Value) != "" {
		return
	}
	if err := global.DB.Model(&model.SysConfig{}).
		Where("id = ?", existing.ID).
		Updates(map[string]interface{}{
			"name":       config.Name,
			"value":      config.Value,
			"value_type": config.ValueType,
			"remark":     config.Remark,
		}).Error; err != nil {
		global.Log.Errorf("回填系统配置失败(%s): %v", config.Key, err)
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
		Attrs(model.SysDept{
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

func backfillExplicitSuperAdminRoles() {
	if !global.DB.Migrator().HasColumn(&model.SysRole{}, "is_super_admin") {
		return
	}

	if err := global.DB.Model(&model.SysRole{}).
		Where("code IN ?", []string{"admin", "system_admin"}).
		Where("is_super_admin = ? OR is_super_admin IS NULL", false).
		Update("is_super_admin", true).Error; err != nil {
		global.Log.Errorf("回填显式超管角色失败: %v", err)
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

func ensureAIToolMenus() {
	aiToolsMenu := model.SysMenu{
		ParentID:   0,
		Name:       "AI工具",
		Path:       "/ai-tools",
		Component:  "Layout",
		Icon:       "robot",
		Sort:       3,
		Type:       1,
		Permission: "ai:tools",
		Status:     1,
		Hidden:     0,
	}

	if err := global.DB.
		Where("permission = ?", aiToolsMenu.Permission).
		Attrs(model.SysMenu{
			ParentID:  aiToolsMenu.ParentID,
			Name:      aiToolsMenu.Name,
			Path:      aiToolsMenu.Path,
			Component: aiToolsMenu.Component,
			Icon:      aiToolsMenu.Icon,
			Sort:      aiToolsMenu.Sort,
			Type:      aiToolsMenu.Type,
			Status:    aiToolsMenu.Status,
			Hidden:    aiToolsMenu.Hidden,
		}).
		FirstOrCreate(&aiToolsMenu).Error; err != nil {
		global.Log.Errorf("补齐AI工具目录失败: %v", err)
		return
	}

	menuDefinitions := []model.SysMenu{
		{
			ParentID:   aiToolsMenu.ID,
			Name:       "AI对话",
			Path:       "/ai/chat",
			Component:  "ai",
			Icon:       "robot",
			Sort:       1,
			Type:       2,
			Permission: "ai:chat:list",
			Status:     1,
			Hidden:     0,
		},
		{
			ParentID:   aiToolsMenu.ID,
			Name:       "AI配置",
			Path:       "/ai/config",
			Component:  "ai/config/index",
			Icon:       "setting",
			Sort:       2,
			Type:       2,
			Permission: "ai:config:list",
			Status:     1,
			Hidden:     0,
		},
	}

	menuIDs := []uint{aiToolsMenu.ID}
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
			global.Log.Errorf("补齐AI工具菜单失败(%s): %v", definition.Permission, err)
			continue
		}
		menuIDs = append(menuIDs, menu.ID)
	}

	grantMenusToRoleCodes(menuIDs, []string{"admin", "system_admin"})
}

func ensureLogAuditMenus() {
	var systemMenu model.SysMenu
	if err := global.DB.Where("path = ? AND type = ?", "/system", 1).First(&systemMenu).Error; err != nil {
		global.Log.Errorf("查询系统管理目录失败: %v", err)
		return
	}

	changed := false
	auditMenu := model.SysMenu{
		ParentID:  systemMenu.ID,
		Name:      "操作审计",
		Path:      "/system/operation-audit",
		Component: "Layout",
		Icon:      "audit",
		Sort:      8,
		Type:      1,
		Status:    1,
		Hidden:    0,
	}
	result := global.DB.
		Where("path = ? AND type = ?", auditMenu.Path, auditMenu.Type).
		Attrs(model.SysMenu{
			ParentID:  auditMenu.ParentID,
			Name:      auditMenu.Name,
			Component: auditMenu.Component,
			Icon:      auditMenu.Icon,
			Sort:      auditMenu.Sort,
			Status:    auditMenu.Status,
			Hidden:    auditMenu.Hidden,
		}).
		FirstOrCreate(&auditMenu)
	if err := result.Error; err != nil {
		global.Log.Errorf("补齐操作审计目录失败: %v", err)
		return
	}
	if result.RowsAffected > 0 {
		changed = true
	}

	menuDefinitions := []model.SysMenu{
		{
			ParentID:   auditMenu.ID,
			Name:       "操作日志",
			Path:       "/monitor/operation-log",
			Component:  "monitor/operation-log/index",
			Icon:       "file-text",
			Sort:       1,
			Type:       2,
			Permission: "monitor:operation-log:list",
			Status:     1,
			Hidden:     0,
		},
		{
			ParentID:   auditMenu.ID,
			Name:       "登录日志",
			Path:       "/monitor/login-log",
			Component:  "monitor/login-log/index",
			Icon:       "login",
			Sort:       2,
			Type:       2,
			Permission: "monitor:login-log:list",
			Status:     1,
			Hidden:     0,
		},
	}

	menuIDs := []uint{auditMenu.ID}
	for _, definition := range menuDefinitions {
		menu := definition
		result := global.DB.
			Where("permission = ? AND type = ?", menu.Permission, menu.Type).
			Attrs(model.SysMenu{
				ParentID:  menu.ParentID,
				Name:      menu.Name,
				Path:      menu.Path,
				Component: menu.Component,
				Icon:      menu.Icon,
				Sort:      menu.Sort,
				Status:    menu.Status,
				Hidden:    menu.Hidden,
			}).
			FirstOrCreate(&menu)
		if err := result.Error; err != nil {
			global.Log.Errorf("补齐日志审计菜单失败(%s): %v", definition.Permission, err)
			continue
		}
		if result.RowsAffected > 0 {
			changed = true
		}

		if menu.ParentID != auditMenu.ID {
			if err := global.DB.Model(&menu).Update("parent_id", auditMenu.ID).Error; err != nil {
				global.Log.Errorf("迁移日志审计菜单父级失败(%s): %v", definition.Permission, err)
				continue
			}
			menu.ParentID = auditMenu.ID
			changed = true
		}
		menuIDs = append(menuIDs, menu.ID)
	}

	grantMenusToRoleCodes(menuIDs, []string{"admin", "system_admin"})
	if changed {
		if err := service.Cache.ClearAllUserInfoCache(); err != nil {
			global.Log.Errorf("清理用户菜单缓存失败: %v", err)
		}
	}
}

func ensureAIApiAccess() {
	apiDefinitions := []model.SysApi{
		{Path: "/api/v1/ai/conversations", Method: "GET", Group: "AI对话", Description: "获取对话列表", NeedAuth: true},
		{Path: "/api/v1/ai/conversations", Method: "POST", Group: "AI对话", Description: "创建对话", NeedAuth: true},
		{Path: "/api/v1/ai/conversations/:id", Method: "GET", Group: "AI对话", Description: "获取对话详情", NeedAuth: true},
		{Path: "/api/v1/ai/conversations/:id", Method: "PUT", Group: "AI对话", Description: "更新对话标题", NeedAuth: true},
		{Path: "/api/v1/ai/conversations/:id", Method: "DELETE", Group: "AI对话", Description: "删除对话", NeedAuth: true},
		{Path: "/api/v1/ai/conversations/:id/messages", Method: "GET", Group: "AI对话", Description: "获取对话消息", NeedAuth: true},
		{Path: "/api/v1/ai/conversations/:id/messages", Method: "DELETE", Group: "AI对话", Description: "清空对话消息", NeedAuth: true},
		{Path: "/api/v1/ai/conversations/:id/clear-context", Method: "POST", Group: "AI对话", Description: "清空上下文", NeedAuth: true},
		{Path: "/api/v1/ai/messages/:id", Method: "DELETE", Group: "AI对话", Description: "删除单条消息", NeedAuth: true},
		{Path: "/api/v1/ai/chat", Method: "POST", Group: "AI对话", Description: "AI对话", NeedAuth: true},
		{Path: "/api/v1/ai/chat/stream", Method: "POST", Group: "AI对话", Description: "AI流式对话", NeedAuth: true},
		{Path: "/api/v1/ai/config", Method: "GET", Group: "AI配置", Description: "获取AI配置", NeedAuth: true},
		{Path: "/api/v1/ai/config", Method: "PUT", Group: "AI配置", Description: "保存AI配置", NeedAuth: true},
		{Path: "/api/v1/ai/test", Method: "POST", Group: "AI配置", Description: "测试AI配置", NeedAuth: true},
		{Path: "/api/v1/ai/providers/models/fetch", Method: "POST", Group: "AI配置", Description: "拉取平台模型列表", NeedAuth: true},
	}

	for _, definition := range apiDefinitions {
		ensureApiAccessForRoleCodes(definition, []string{"admin", "system_admin"})
	}
}

func ensureApiAccessInheritedFrom(api model.SysApi, sourcePath, sourceMethod string) {
	if err := global.DB.
		Where("path = ? AND method = ?", api.Path, api.Method).
		Attrs(model.SysApi{
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

func ensureApiAccessForRoleCodes(api model.SysApi, roleCodes []string) {
	if err := global.DB.
		Where("path = ? AND method = ?", api.Path, api.Method).
		Attrs(model.SysApi{
			Group:       api.Group,
			Description: api.Description,
			NeedAuth:    api.NeedAuth,
		}).
		FirstOrCreate(&api).Error; err != nil {
		global.Log.Errorf("补齐系统API失败(%s %s): %v", api.Method, api.Path, err)
		return
	}

	var roles []model.SysRole
	if err := global.DB.Where("code IN ?", roleCodes).Find(&roles).Error; err != nil {
		global.Log.Errorf("查询内置角色失败: %v", err)
		return
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

func ensureUserOperationMenus() {
	var userMenu model.SysMenu
	if err := global.DB.Where("permission = ? AND type = ?", "system:user:list", 2).First(&userMenu).Error; err != nil {
		global.Log.Errorf("查询用户管理菜单失败: %v", err)
		return
	}

	menuDefinitions := []model.SysMenu{
		{
			ParentID:   userMenu.ID,
			Name:       "用户重置密码",
			Path:       "",
			Component:  "",
			Icon:       "",
			Sort:       4,
			Type:       3,
			Permission: "system:user:resetPwd",
			Status:     1,
			Hidden:     0,
		},
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
		{
			ParentID:   userMenu.ID,
			Name:       "批量重置密码",
			Path:       "",
			Component:  "",
			Icon:       "",
			Sort:       7,
			Type:       3,
			Permission: "system:user:batchResetPwd",
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
			global.Log.Errorf("补齐用户操作按钮失败(%s): %v", definition.Permission, err)
			continue
		}

		sourcePermission := "system:user:edit"
		if definition.Permission == "system:user:batchResetPwd" {
			sourcePermission = "system:user:resetPwd"
		}
		grantMenuToRolesWithPermission(menu.ID, sourcePermission)
	}
}

func ensureUserImportExportMenus() {
	var userMenu model.SysMenu
	if err := global.DB.Where("permission = ? AND type = ?", "system:user:list", 2).First(&userMenu).Error; err != nil {
		global.Log.Errorf("查询用户管理菜单失败(导入导出): %v", err)
		return
	}

	menuDefinitions := []model.SysMenu{
		{
			ParentID:   userMenu.ID,
			Name:       "导入用户",
			Sort:       8,
			Type:       3,
			Permission: "system:user:import",
			Status:     1,
			Hidden:     0,
		},
		{
			ParentID:   userMenu.ID,
			Name:       "导出用户",
			Sort:       9,
			Type:       3,
			Permission: "system:user:export",
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
				Sort:      menu.Sort,
				Type:      menu.Type,
				Status:    menu.Status,
				Hidden:    menu.Hidden,
			}).
			FirstOrCreate(&menu).Error; err != nil {
			global.Log.Errorf("补齐用户导入导出按钮失败(%s): %v", definition.Permission, err)
			continue
		}

		grantMenuToRolesWithPermission(menu.ID, "system:user:list")
	}

	// 补齐导入导出API访问权限
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/users/import-template",
		Method:      "GET",
		Group:       "用户管理",
		Description: "下载导入模板",
		NeedAuth:    true,
	}, "/api/v1/users", "GET")
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/users/import",
		Method:      "POST",
		Group:       "用户管理",
		Description: "导入用户",
		NeedAuth:    true,
	}, "/api/v1/users", "POST")
	ensureApiAccessInheritedFrom(model.SysApi{
		Path:        "/api/v1/users/export",
		Method:      "GET",
		Group:       "用户管理",
		Description: "导出用户",
		NeedAuth:    true,
	}, "/api/v1/users", "GET")
}

func cleanupStorageBuiltInData() {
	tx := global.DB.Begin()
	if tx.Error != nil {
		global.Log.Errorf("初始化存储内置数据清理事务失败: %v", tx.Error)
		return
	}

	var menuIDs []uint
	if err := tx.Model(&model.SysMenu{}).
		Where("path = ? OR permission = ?", "/system/storage", "system:storage:list").
		Pluck("id", &menuIDs).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("查询旧存储菜单失败: %v", err)
		return
	}
	if len(menuIDs) > 0 {
		if err := tx.Exec("DELETE FROM sys_role_menu WHERE sys_menu_id IN ?", menuIDs).Error; err != nil {
			tx.Rollback()
			global.Log.Errorf("清理旧存储菜单角色关联失败: %v", err)
			return
		}
		if err := tx.Where("id IN ?", menuIDs).Delete(&model.SysMenu{}).Error; err != nil {
			tx.Rollback()
			global.Log.Errorf("删除旧存储菜单失败: %v", err)
			return
		}
	}

	type apiPolicy struct {
		Path   string
		Method string
	}

	storageApis := []apiPolicy{
		{Path: "/api/v1/storages", Method: "GET"},
		{Path: "/api/v1/storages/:id", Method: "GET"},
		{Path: "/api/v1/storages", Method: "POST"},
		{Path: "/api/v1/storages/:id", Method: "PUT"},
		{Path: "/api/v1/storages/:id", Method: "DELETE"},
		{Path: "/api/v1/storages/:id/default", Method: "PUT"},
		{Path: "/api/v1/storages/test", Method: "POST"},
	}
	apiPaths := make([]string, 0, len(storageApis))
	for _, api := range storageApis {
		apiPaths = append(apiPaths, api.Path)
	}

	var apiIDs []uint
	if err := tx.Model(&model.SysApi{}).Where("path IN ?", apiPaths).Pluck("id", &apiIDs).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("查询旧存储 API 失败: %v", err)
		return
	}
	if len(apiIDs) > 0 {
		if err := tx.Exec("DELETE FROM sys_role_api WHERE sys_api_id IN ?", apiIDs).Error; err != nil {
			tx.Rollback()
			global.Log.Errorf("清理旧存储 API 角色关联失败: %v", err)
			return
		}
		if err := tx.Where("id IN ?", apiIDs).Delete(&model.SysApi{}).Error; err != nil {
			tx.Rollback()
			global.Log.Errorf("删除旧存储 API 失败: %v", err)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		global.Log.Errorf("提交旧存储内置数据清理失败: %v", err)
		return
	}

	if global.Enforcer != nil {
		policyChanged := false
		for _, api := range storageApis {
			if ok, err := global.Enforcer.RemoveFilteredPolicy(1, api.Path, api.Method); err != nil {
				global.Log.Errorf("清理旧存储 Casbin 策略失败(%s %s): %v", api.Method, api.Path, err)
			} else if ok {
				policyChanged = true
			}
		}
		if policyChanged {
			_ = global.Enforcer.SavePolicy()
		}
	}
}

func cleanupSlowLogBuiltInData() {
	tx := global.DB.Begin()
	if tx.Error != nil {
		global.Log.Errorf("初始化慢查询日志内置数据清理事务失败: %v", tx.Error)
		return
	}

	var menuIDs []uint
	if err := tx.Model(&model.SysMenu{}).
		Where("path = ? OR permission = ?", "/monitor/show-log", "monitor:show-log:list").
		Pluck("id", &menuIDs).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("查询慢查询日志菜单失败: %v", err)
		return
	}
	menuChanged := len(menuIDs) > 0
	if menuChanged {
		if err := tx.Exec("DELETE FROM sys_role_menu WHERE sys_menu_id IN ?", menuIDs).Error; err != nil {
			tx.Rollback()
			global.Log.Errorf("清理慢查询日志菜单角色关联失败: %v", err)
			return
		}
		if err := tx.Where("id IN ?", menuIDs).Delete(&model.SysMenu{}).Error; err != nil {
			tx.Rollback()
			global.Log.Errorf("删除慢查询日志菜单失败: %v", err)
			return
		}
	}

	type apiPolicy struct {
		Path   string
		Method string
	}

	slowLogApis := []apiPolicy{
		{Path: "/api/v1/logs/slow", Method: "GET"},
		{Path: "/api/v1/logs/slow/:id", Method: "DELETE"},
	}
	apiPaths := make([]string, 0, len(slowLogApis))
	for _, api := range slowLogApis {
		apiPaths = append(apiPaths, api.Path)
	}

	var apiIDs []uint
	if err := tx.Model(&model.SysApi{}).Where("path IN ?", apiPaths).Pluck("id", &apiIDs).Error; err != nil {
		tx.Rollback()
		global.Log.Errorf("查询慢查询日志 API 失败: %v", err)
		return
	}
	if len(apiIDs) > 0 {
		if err := tx.Exec("DELETE FROM sys_role_api WHERE sys_api_id IN ?", apiIDs).Error; err != nil {
			tx.Rollback()
			global.Log.Errorf("清理慢查询日志 API 角色关联失败: %v", err)
			return
		}
		if err := tx.Where("id IN ?", apiIDs).Delete(&model.SysApi{}).Error; err != nil {
			tx.Rollback()
			global.Log.Errorf("删除慢查询日志 API 失败: %v", err)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		global.Log.Errorf("提交慢查询日志内置数据清理失败: %v", err)
		return
	}

	if menuChanged {
		if err := service.Cache.ClearAllUserInfoCache(); err != nil {
			global.Log.Errorf("清理慢查询日志菜单缓存失败: %v", err)
		}
	}

	if global.Enforcer != nil {
		policyChanged := false
		for _, api := range slowLogApis {
			if ok, err := global.Enforcer.RemoveFilteredPolicy(1, api.Path, api.Method); err != nil {
				global.Log.Errorf("清理慢查询日志 Casbin 策略失败(%s %s): %v", api.Method, api.Path, err)
			} else if ok {
				policyChanged = true
			}
		}
		if policyChanged {
			_ = global.Enforcer.SavePolicy()
		}
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
