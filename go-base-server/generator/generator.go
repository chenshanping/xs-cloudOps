package generator

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

//go:embed template/*.tpl
var templateFS embed.FS

// Generator 代码生成器
type Generator struct {
	Config         *GeneratorConfig
	ServerPath     string
	WebPath        string
	parentMenuPath string // 父菜单路径
}

// SetParentMenuPath 设置父菜单路径
func (g *Generator) SetParentMenuPath(path string) {
	g.parentMenuPath = path
}

// TemplateData 模板数据
type TemplateData struct {
	TableName           string
	ModelName           string
	ModuleName          string
	Description         string
	RoutePath           string
	Columns             []ColumnConfig
	SearchColumns       []ColumnConfig
	ListColumns         []ColumnConfig
	FormColumns         []ColumnConfig
	SortColumns         []ColumnConfig // 可排序字段
	UniqueColumns       []ColumnConfig // 唯一字段
	Relations           []RelationTemplateData
	HasTimeField        bool
	HasPreloads         bool
	HasMany2Many        bool
	HasMultiFiles       bool // 是否有多图片/多文件字段 (images/files)
	HasFiles            bool // 是否有文件字段 (file/files)
	HasEditor           bool // 是否有富文本编辑器字段
	HasRelations        bool // 是否有 belongsTo/many2many 关联
	HasTreeLayout       bool // 是否有左树右表布局
	HasSearchFields     bool // 是否有搜索字段（搜索列 + 关联下拉 + 创建人）
	OnlyCreatedBySearch bool // 搜索字段仅有 created_by（无其他搜索列和 belongsTo 关联）
	HasDictSelect       bool // 是否有使用数据字典的下拉框
	Preloads            []string
	// 时间字段开关（用于模型模板）
	HasCreatedAt bool
	HasUpdatedAt bool
	HasDeletedAt bool
	HasCreatedBy bool // 是否包含创建人字段
	// 创建人身份关联
	CreatedByProfileTable string // 创建者身份关联表名
	CreatedByProfileModel string // 创建者身份关联模型名（大驼峰）
	CreatedByProfileField string // 身份表中要显示的字段
	HasCreatedByProfile   bool   // 是否有创建者身份关联
	// 数据隔离配置
	DataIsolation bool   // 是否启用数据隔离
	AdminRoleIds  string // 管理员角色ID列表
	// 审批功能
	HasAudit bool // 是否启用审批功能
	// 前台接口
	GenerateFrontendApi bool // 是否生成前台用户使用的接口
	// SQL生成相关
	GenerateTime string
	// 菜单配置
	MenuConfig     *MenuConfig
	ParentMenuPath string // 父菜单路径
	ComponentPath  string // 组件路径
	HasMenu        bool   // 是否配置了菜单
	// 用户关联
	LinkToUser      bool   // 是否关联用户表（一对一）
	ProfileName     string // 身份显示名称
	ProfileIcon     string // 身份图标
	ProfileRoleCode string // 身份限定角色编码
	// 统计配置
	HasStats         bool   // 是否启用统计功能
	StatsGroupField  string // 分组字段
	StatsGroupColumn string // 分组字段（数据库列名）
	StatsGroupDisplay string // 分组显示字段
	StatsSumField    string // 求和字段
	StatsSumColumn   string // 求和字段（数据库列名）
	StatsTimeField   string // 时间字段
	StatsTimeColumn  string // 时间字段（数据库列名）
	StatsChartTypes  []string // 图表类型
}

// RelationTemplateData 关联关系模板数据
type RelationTemplateData struct {
	RelationType   string
	RelatedTable   string
	RelatedModule  string // 关联模块名（用于 API import 路径）
	RelatedModel   string
	FieldName      string
	JsonName       string
	ForeignKey     string
	ForeignKeyJson string
	ReferenceKey   string
	JoinTable      string
	DisplayField   string // 显示字段（如name）
	Comment        string // 关联注释（如"产品类型"）
	IsRequired     bool   // 是否必填
	UseOptionsApi  bool   // 使用轻量options接口
	UseTreeLayout  bool   // 使用左树右表布局
}

// NewGenerator 创建生成器
func NewGenerator(config *GeneratorConfig, serverPath, webPath string) *Generator {
	return &Generator{
		Config:     config,
		ServerPath: serverPath,
		WebPath:    webPath,
	}
}

// Preview 预览生成的代码
func (g *Generator) Preview() (*PreviewResult, error) {
	result := &PreviewResult{
		Files: make([]GeneratedFile, 0),
	}

	data := g.buildTemplateData()

	if g.Config.GenerateBackend {
		// 生成后端代码
		files, err := g.generateBackendCode(data)
		if err != nil {
			return nil, err
		}
		result.Files = append(result.Files, files...)
	}

	if g.Config.GenerateFrontend {
		// 生成前端代码
		files, err := g.generateFrontendCode(data)
		if err != nil {
			return nil, err
		}
		result.Files = append(result.Files, files...)
	}

	return result, nil
}

// Generate 生成代码并写入文件
func (g *Generator) Generate() (*PreviewResult, error) {
	result, err := g.Preview()
	if err != nil {
		return nil, err
	}

	for _, file := range result.Files {
		// 确保目录存在
		dir := filepath.Dir(file.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建目录失败: %v", err)
		}

		// 写入文件
		if err := os.WriteFile(file.Path, []byte(file.Content), 0644); err != nil {
			return nil, fmt.Errorf("写入文件失败: %v", err)
		}
	}

	return result, nil
}

// buildTemplateData 构建模板数据
func (g *Generator) buildTemplateData() *TemplateData {
	modelName := ToPascalCase(g.Config.TableName)

	// 过滤各类字段
	searchColumns := make([]ColumnConfig, 0)
	listColumns := make([]ColumnConfig, 0)
	formColumns := make([]ColumnConfig, 0)
	sortColumns := make([]ColumnConfig, 0)
	uniqueColumns := make([]ColumnConfig, 0)
	hasTimeField := false
	hasMultiFiles := false
	hasFiles := false
	hasEditor := false
	hasDictSelect := false

	// 收集 belongsTo 关联的外键字段，用于排除
	belongsToForeignKeys := make(map[string]bool)
	for _, rel := range g.Config.Relations {
		if rel.RelationType == "belongsTo" && rel.ForeignKey != "" {
			belongsToForeignKeys[rel.ForeignKey] = true
		}
	}

	// belongsTo 外键字段类型强制为 uint
	for i := range g.Config.Columns {
		if belongsToForeignKeys[g.Config.Columns[i].ColumnName] {
			g.Config.Columns[i].FieldType = "uint"
			g.Config.Columns[i].TsType = "number"
		}
	}

	for _, col := range g.Config.Columns {
		// 搜索时排除 belongsTo 外键字段（会用关联下拉框搜索）
		if col.IsSearchable && !belongsToForeignKeys[col.ColumnName] {
			// 默认搜索类型为 like
			if col.SearchType == "" {
				col.SearchType = "like"
			}
			searchColumns = append(searchColumns, col)
		}
		// 列表显示时排除 belongsTo 外键字段（会用关联对象显示）
		if col.IsListVisible && !belongsToForeignKeys[col.ColumnName] {
			listColumns = append(listColumns, col)
		}
		// 表单显示时排除 belongsTo 外键字段（会用关联下拉框）
		if col.IsFormVisible && !belongsToForeignKeys[col.ColumnName] {
			formColumns = append(formColumns, col)
		}
		if col.IsSortable {
			sortColumns = append(sortColumns, col)
		}
		if col.IsUnique {
			uniqueColumns = append(uniqueColumns, col)
		}
		if col.FieldType == "time.Time" {
			hasTimeField = true
		}
		if col.FormType == "images" || col.FormType == "files" {
			hasMultiFiles = true
		}
		if col.FormType == "file" || col.FormType == "files" {
			hasFiles = true
		}
		if col.FormType == "editor" {
			hasEditor = true
		}
		if col.DictType != "" {
			hasDictSelect = true
		}
	}

	// 如果有审批功能，也需要time包（AuditTime字段）
	if g.Config.HasAudit {
		hasTimeField = true
	}

	// 处理关联关系
	relations := make([]RelationTemplateData, 0)
	preloads := make([]string, 0)
	hasMany2Many := false

	for _, rel := range g.Config.Relations {
		fieldName := ToPascalCase(rel.RelatedTable)
		if rel.RelationType == "hasMany" || rel.RelationType == "many2many" {
			fieldName = fieldName + "s" // 复数形式
		}

		// 显示字段默认为 name
		displayField := rel.DisplayField
		if displayField == "" {
			displayField = "name"
		}
		// 注释默认为关联表名
		comment := rel.Comment
		if comment == "" {
			comment = rel.RelatedTable
		}

		// 关联模块名：优先使用配置的 RelatedModule，否则使用表名
		relatedModule := rel.RelatedModule
		if relatedModule == "" {
			relatedModule = rel.RelatedTable
		}

	relData := RelationTemplateData{
			RelationType:   rel.RelationType,
			RelatedTable:   rel.RelatedTable,
			RelatedModule:  relatedModule,
			RelatedModel:   ToPascalCase(rel.RelatedTable),
			FieldName:      fieldName,
			JsonName:       ToSnakeCase(fieldName),
			ForeignKey:     rel.ForeignKey,
			ForeignKeyJson: ToSnakeCase(rel.ForeignKey),
			ReferenceKey:   rel.ReferenceKey,
			JoinTable:      rel.JoinTable,
			DisplayField:   displayField,
			Comment:        comment,
			IsRequired:     rel.IsRequired,
			UseOptionsApi:  rel.UseOptionsApi,
			UseTreeLayout:  rel.UseTreeLayout,
		}
		relations = append(relations, relData)
		preloads = append(preloads, fieldName)

		if rel.RelationType == "many2many" {
			hasMany2Many = true
		}
	}

	// 检查是否有 belongsTo/many2many 关联
	hasRelations := false
	hasTreeLayout := false
	hasBelongsToSearch := false // belongsTo 不使用左树布局的关联（会生成搜索下拉框）
	for _, rel := range relations {
		if rel.RelationType == "belongsTo" || rel.RelationType == "many2many" {
			hasRelations = true
		}
		if rel.RelationType == "belongsTo" && rel.UseTreeLayout {
			hasTreeLayout = true
		}
		if rel.RelationType == "belongsTo" && !rel.UseTreeLayout {
			hasBelongsToSearch = true
		}
	}

	// 检查是否有搜索字段：搜索列 + belongsTo下拉(非左树) + 创建人(非数据隔离或管理员可见)
	hasSearchFields := len(searchColumns) > 0 || hasBelongsToSearch || g.Config.HasCreatedBy
	// 搜索字段仅有 created_by（无其他搜索列和 belongsTo 关联）
	onlyCreatedBySearch := len(searchColumns) == 0 && !hasBelongsToSearch && g.Config.HasCreatedBy

	// 菜单配置处理
	hasMenu := g.Config.MenuConfig != nil && g.Config.MenuConfig.MenuName != ""
	parentMenuPath := ""
	componentPath := g.Config.ModuleName + "/index"
	if hasMenu && g.Config.MenuConfig.ParentID > 0 {
		// 有父菜单时，组件路径需要包含父菜单路径
		// 父菜单路径需要在调用时传入（通过 SetParentMenuPath 方法）
		parentMenuPath = g.parentMenuPath
		if parentMenuPath != "" {
			componentPath = parentMenuPath + "/" + g.Config.ModuleName + "/index"
		}
	}

	// 角色限定 RoleCode：如果启用了用户关联且没有指定 RoleCode，默认使用模块名
	profileRoleCode := g.Config.ProfileRoleCode
	if g.Config.LinkToUser && profileRoleCode == "" {
		profileRoleCode = g.Config.ModuleName
	}

	// ProfileName：如果为空，默认使用 Description
	profileName := g.Config.ProfileName
	if g.Config.LinkToUser && profileName == "" {
		profileName = g.Config.Description
	}

	return &TemplateData{
		TableName:           g.Config.TableName,
		ModelName:           modelName,
		ModuleName:          g.Config.ModuleName,
		Description:         g.Config.Description,
		RoutePath:           g.Config.ModuleName,
		Columns:             g.Config.Columns,
		SearchColumns:       searchColumns,
		ListColumns:         listColumns,
		FormColumns:         formColumns,
		SortColumns:         sortColumns,
		UniqueColumns:       uniqueColumns,
		Relations:           relations,
		HasTimeField:        hasTimeField,
		HasPreloads:         len(preloads) > 0,
		HasMany2Many:        hasMany2Many,
		HasMultiFiles:       hasMultiFiles,
		HasFiles:            hasFiles,
		HasEditor:           hasEditor,
		HasRelations:        hasRelations,
		HasTreeLayout:       hasTreeLayout,
		HasSearchFields:     hasSearchFields,
		OnlyCreatedBySearch: onlyCreatedBySearch,
		HasDictSelect:       hasDictSelect,
		Preloads:            preloads,
		HasCreatedAt:        g.Config.HasCreatedAt,
		HasUpdatedAt:        g.Config.HasUpdatedAt,
		HasDeletedAt:        g.Config.HasDeletedAt,
		HasCreatedBy:          g.Config.HasCreatedBy,
		CreatedByProfileTable: g.Config.CreatedByProfileTable,
		CreatedByProfileModel: ToPascalCase(g.Config.CreatedByProfileTable),
		CreatedByProfileField: g.Config.CreatedByProfileField,
		HasCreatedByProfile:   g.Config.HasCreatedBy && g.Config.CreatedByProfileTable != "",
		DataIsolation:         g.Config.DataIsolation && g.Config.HasCreatedBy, // 数据隔离仅在有创建人字段时生效
		AdminRoleIds:          g.Config.AdminRoleIds,
		HasAudit:            g.Config.HasAudit,
		GenerateFrontendApi: g.Config.GenerateFrontendApi,
		GenerateTime:        time.Now().Format("2006-01-02 15:04:05"),
		MenuConfig:          g.Config.MenuConfig,
		ParentMenuPath:      parentMenuPath,
		ComponentPath:       componentPath,
		HasMenu:             hasMenu,
		LinkToUser:      g.Config.LinkToUser,
		ProfileName:     profileName,
		ProfileIcon:     g.Config.ProfileIcon,
		ProfileRoleCode: profileRoleCode,
		// 统计配置
		HasStats:          g.Config.StatsConfig != nil && g.Config.StatsConfig.Enabled,
		StatsGroupField:   g.getStatsField("group_field"),
		StatsGroupColumn:  g.getStatsColumn("group_field"),
		StatsGroupDisplay: g.getStatsGroupDisplay(),
		StatsSumField:     g.getStatsField("sum_field"),
		StatsSumColumn:    g.getStatsColumn("sum_field"),
		StatsTimeField:    g.getStatsField("time_field"),
		StatsTimeColumn:   g.getStatsColumn("time_field"),
		StatsChartTypes:   g.getStatsChartTypes(),
	}
}

// generateBackendCode 生成后端代码
func (g *Generator) generateBackendCode(data *TemplateData) ([]GeneratedFile, error) {
	files := make([]GeneratedFile, 0)

	// 添加 GormTag 到列配置
	for i := range data.Columns {
		data.Columns[i] = g.addGormTag(data.Columns[i])
	}
	for i := range data.FormColumns {
		data.FormColumns[i] = g.addGormTag(data.FormColumns[i])
	}

	// 核心后端文件（始终生成）
	templates := []struct {
		Name string
		Path string
	}{
		{"model.tpl", filepath.Join(g.ServerPath, "model", g.Config.ModuleName+".go")},
		{"request.tpl", filepath.Join(g.ServerPath, "model", "request", g.Config.ModuleName+"_request.go")},
		{"service.tpl", filepath.Join(g.ServerPath, "service", g.Config.ModuleName+".go")},
	}

	// API接口相关文件（仅在生成前端代码时生成）
	if g.Config.GenerateFrontend {
		templates = append(templates,
			struct {
				Name string
				Path string
			}{"api.tpl", filepath.Join(g.ServerPath, "api", "v1", g.Config.ModuleName+".go")},
			struct {
				Name string
				Path string
			}{"router.tpl", filepath.Join(g.ServerPath, "router", "modules", g.Config.ModuleName+".go")},
		)
	}

	for _, t := range templates {
		content, err := g.renderTemplate(t.Name, data)
		if err != nil {
			return nil, fmt.Errorf("渲染模板 %s 失败: %v", t.Name, err)
		}
		files = append(files, GeneratedFile{
			Path:    t.Path,
			Content: content,
			Type:    "backend",
		})
	}

	// 生成SQL文件
	if g.Config.GenerateSQL {
		// 为columns添加SqlType
		for i := range data.Columns {
			data.Columns[i] = g.addSqlType(data.Columns[i])
		}
		content, err := g.renderTemplate("sql.tpl", data)
		if err != nil {
			return nil, fmt.Errorf("渲染SQL模板失败: %v", err)
		}
		files = append(files, GeneratedFile{
			Path:    filepath.Join(g.ServerPath, "sql", g.Config.ModuleName+".sql"),
			Content: content,
			Type:    "sql",
		})

		// 如果有多对多关系，生成中间表SQL文件
		if data.HasMany2Many {
			content, err := g.renderTemplate("sql_join.tpl", data)
			if err != nil {
				return nil, fmt.Errorf("渲染中间表SQL模板失败: %v", err)
			}
			files = append(files, GeneratedFile{
				Path:    filepath.Join(g.ServerPath, "sql", g.Config.ModuleName+"_join.sql"),
				Content: content,
				Type:    "sql",
			})
		}
	}

	// 生成菜单SQL文件（仅在生成前端代码且配置了菜单时）
	if data.HasMenu && g.Config.GenerateFrontend {
		content, err := g.renderTemplate("menu_sql.tpl", data)
		if err != nil {
			return nil, fmt.Errorf("渲染菜单SQL模板失败: %v", err)
		}
		files = append(files, GeneratedFile{
			Path:    filepath.Join(g.ServerPath, "sql", g.Config.ModuleName+"_menu.sql"),
			Content: content,
			Type:    "sql",
		})
	}

	// 生成角色SQL文件（仅在启用用户关联时）
	if data.LinkToUser {
		content, err := g.renderTemplate("role_sql.tpl", data)
		if err != nil {
			return nil, fmt.Errorf("渲染角色SQL模板失败: %v", err)
		}
		files = append(files, GeneratedFile{
			Path:    filepath.Join(g.ServerPath, "sql", g.Config.ModuleName+"_role.sql"),
			Content: content,
			Type:    "sql",
		})
	}

	return files, nil
}

// generateFrontendCode 生成前端代码
func (g *Generator) generateFrontendCode(data *TemplateData) ([]GeneratedFile, error) {
	files := make([]GeneratedFile, 0)

	// 确定视图文件路径（如果有父菜单，路径需要包含父菜单路径）
	viewPath := g.Config.ModuleName
	if data.ParentMenuPath != "" {
		viewPath = filepath.Join(data.ParentMenuPath, g.Config.ModuleName)
	}

	// API文件
	content, err := g.renderTemplate("frontend_api.tpl", data)
	if err != nil {
		return nil, err
	}
	files = append(files, GeneratedFile{
		Path:    filepath.Join(g.WebPath, "src", "api", g.Config.ModuleName+".ts"),
		Content: content,
		Type:    "frontend",
	})

	// 视图文件
	content, err = g.renderTemplate("frontend_view.tpl", data)
	if err != nil {
		return nil, err
	}
	files = append(files, GeneratedFile{
		Path:    filepath.Join(g.WebPath, "src", "views", viewPath, "index.vue"),
		Content: content,
		Type:    "frontend",
	})

	// 表单抽屉组件
	content, err = g.renderTemplate("frontend_form.tpl", data)
	if err != nil {
		return nil, err
	}
	files = append(files, GeneratedFile{
		Path:    filepath.Join(g.WebPath, "src", "views", viewPath, "components", data.ModelName+"Form.vue"),
		Content: content,
		Type:    "frontend",
	})

	// 类型定义文件
	content, err = g.renderTemplate("frontend_type.tpl", data)
	if err != nil {
		return nil, err
	}
	files = append(files, GeneratedFile{
		Path:    filepath.Join(g.WebPath, "src", "types", g.Config.ModuleName+".ts"),
		Content: content,
		Type:    "frontend",
	})

	// 统计组件（仅在启用统计时生成）
	if data.HasStats {
		content, err = g.renderTemplate("frontend_stats.tpl", data)
		if err != nil {
			return nil, err
		}
		files = append(files, GeneratedFile{
			Path:    filepath.Join(g.WebPath, "src", "views", viewPath, "components", data.ModelName+"Stats.vue"),
			Content: content,
			Type:    "frontend",
		})
	}

	// LinkToUser 功能不再生成单独的 my 页面，身份信息已合并到个人中心

	return files, nil
}

// renderTemplate 渲染模板
func (g *Generator) renderTemplate(name string, data interface{}) (string, error) {
	content, err := templateFS.ReadFile("template/" + name)
	if err != nil {
		return "", err
	}

	// 注册自定义模板函数
	funcMap := template.FuncMap{
		"ToPascalCase": ToPascalCase,
		"ToSnakeCase":  ToSnakeCase,
		"TrimSuffix":   strings.TrimSuffix,
	}

	tmpl, err := template.New(name).Funcs(funcMap).Parse(string(content))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// addGormTag 添加GORM标签
func (g *Generator) addGormTag(col ColumnConfig) ColumnConfig {
	var tags []string

	// 根据字段类型添加大小限制
	switch col.FieldType {
	case "string":
		tags = append(tags, "size:255")
	}

	// 添加注释
	if col.Comment != "" {
		tags = append(tags, "comment:"+col.Comment)
	}

	// 组装GormTag
	if len(tags) > 0 {
		col.GormTag = strings.Join(tags, ";")
	}

	return col
}

// addSqlType 添加SQL类型
func (g *Generator) addSqlType(col ColumnConfig) ColumnConfig {
	// 如果已指定数据库类型，使用指定的
	if col.DbType != "" {
		switch strings.ToUpper(col.DbType) {
		case "VARCHAR", "CHAR":
			length := col.DbLength
			if length <= 0 {
				length = 255
			}
			col.SqlType = fmt.Sprintf("%s(%d)", strings.ToUpper(col.DbType), length)
		case "DECIMAL", "NUMERIC":
			// 默认精度
			col.SqlType = fmt.Sprintf("%s(10,2)", strings.ToUpper(col.DbType))
		default:
			col.SqlType = strings.ToUpper(col.DbType)
		}
		return col
	}

	// 根据Go类型推断SQL类型
	switch col.FieldType {
	case "string":
		// 富文本编辑器默认使用 LONGTEXT
		if col.FormType == "editor" {
			col.SqlType = "LONGTEXT"
			return col
		}
		length := col.DbLength
		if length <= 0 {
			length = 255
		}
		if length > 65535 {
			col.SqlType = "LONGTEXT"
		} else if length > 255 {
			col.SqlType = "TEXT"
		} else {
			col.SqlType = fmt.Sprintf("VARCHAR(%d)", length)
		}
	case "int", "int32":
		col.SqlType = "INT"
	case "int8":
		col.SqlType = "TINYINT"
	case "int16":
		col.SqlType = "SMALLINT"
	case "int64":
		col.SqlType = "BIGINT"
	case "uint", "uint32":
		col.SqlType = "INT UNSIGNED"
	case "uint8":
		col.SqlType = "TINYINT UNSIGNED"
	case "uint16":
		col.SqlType = "SMALLINT UNSIGNED"
	case "uint64":
		col.SqlType = "BIGINT UNSIGNED"
	case "float32":
		col.SqlType = "FLOAT"
	case "float64":
		col.SqlType = "DOUBLE"
	case "bool":
		col.SqlType = "TINYINT(1)"
	case "time.Time":
		col.SqlType = "DATETIME"
	default:
		// 默认类型
		col.SqlType = "VARCHAR(255)"
	}

	return col
}

// DeleteModule 删除模块生成的代码
func DeleteModule(moduleName, serverPath, webPath string) error {
	return DeleteModuleWithParentPath(moduleName, "", serverPath, webPath)
}

// DeleteModuleWithParentPath 删除模块生成的代码（支持父菜单路径）
func DeleteModuleWithParentPath(moduleName, parentMenuPath, serverPath, webPath string) error {
	// 删除后端文件
	backendFiles := []string{
		filepath.Join(serverPath, "model", moduleName+".go"),
		filepath.Join(serverPath, "model", "request", moduleName+"_request.go"),
		filepath.Join(serverPath, "api", "v1", moduleName+".go"),
		filepath.Join(serverPath, "service", moduleName+".go"),
		filepath.Join(serverPath, "router", "modules", moduleName+".go"),
		filepath.Join(serverPath, "sql", moduleName+".sql"),
		filepath.Join(serverPath, "sql", moduleName+"_menu.sql"),
	}

	for _, f := range backendFiles {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("删除文件 %s 失败: %v", f, err)
		}
	}

	// 删除前端文件
	frontendFiles := []string{
		filepath.Join(webPath, "src", "api", moduleName+".ts"),
		filepath.Join(webPath, "src", "types", moduleName+".ts"),
	}
	for _, f := range frontendFiles {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("删除文件 %s 失败: %v", f, err)
		}
	}

	// 删除前端视图目录（考虑父菜单路径）
	viewPath := moduleName
	if parentMenuPath != "" {
		viewPath = filepath.Join(parentMenuPath, moduleName)
	}
	viewDir := filepath.Join(webPath, "src", "views", viewPath)
	if err := os.RemoveAll(viewDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除目录 %s 失败: %v", viewDir, err)
	}

	return nil
}

// GetGeneratedModules 获取已生成的模块列表
// 通过检查 model 目录判断（因为 model 始终会生成）
func GetGeneratedModules(serverPath string) ([]string, error) {
	modelDir := filepath.Join(serverPath, "model")

	entries, err := os.ReadDir(modelDir)
	if err != nil {
		return nil, err
	}

	// 系统内置模型
	builtinModels := map[string]bool{
		"base.go":          true,
		"sys_user.go":      true,
		"sys_role.go":      true,
		"sys_menu.go":      true,
		"sys_api.go":       true,
		"sys_log.go":       true,
		"sys_file.go":      true,
		"sys_storage.go":   true,
		"sys_config.go":    true,
		"sys_generator.go": true,
	}

	modules := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		// 跳过系统内置模型和sys_前缀文件
		if builtinModels[name] || strings.HasPrefix(name, "sys_") {
			continue
		}
		// 移除.go后缀
		moduleName := strings.TrimSuffix(name, ".go")
		modules = append(modules, moduleName)
	}

	return modules, nil
}

// getStatsField 获取统计字段名（大驼峰形式）
func (g *Generator) getStatsField(fieldType string) string {
	if g.Config.StatsConfig == nil || !g.Config.StatsConfig.Enabled {
		return ""
	}
	switch fieldType {
	case "group_field":
		return ToPascalCase(g.Config.StatsConfig.GroupField)
	case "sum_field":
		return ToPascalCase(g.Config.StatsConfig.SumField)
	case "time_field":
		return ToPascalCase(g.Config.StatsConfig.TimeField)
	}
	return ""
}

// getStatsColumn 获取统计字段名（数据库列名）
func (g *Generator) getStatsColumn(fieldType string) string {
	if g.Config.StatsConfig == nil || !g.Config.StatsConfig.Enabled {
		return ""
	}
	switch fieldType {
	case "group_field":
		return g.Config.StatsConfig.GroupField
	case "sum_field":
		return g.Config.StatsConfig.SumField
	case "time_field":
		return g.Config.StatsConfig.TimeField
	}
	return ""
}

// getStatsGroupDisplay 获取分组显示字段
func (g *Generator) getStatsGroupDisplay() string {
	if g.Config.StatsConfig == nil || !g.Config.StatsConfig.Enabled {
		return ""
	}
	return g.Config.StatsConfig.GroupDisplay
}

// getStatsChartTypes 获取图表类型列表
func (g *Generator) getStatsChartTypes() []string {
	if g.Config.StatsConfig == nil || !g.Config.StatsConfig.Enabled {
		return nil
	}
	return g.Config.StatsConfig.ChartTypes
}
