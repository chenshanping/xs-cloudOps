package generator

// StatsConfig 统计配置
type StatsConfig struct {
	Enabled      bool     `json:"enabled"`       // 是否启用统计
	GroupField   string   `json:"group_field"`   // 分组字段（如 category_id, status）
	GroupDisplay string   `json:"group_display"` // 分组显示字段（用于显示名称，如 Category.Name）
	SumField     string   `json:"sum_field"`     // 求和字段（可选，如 price, amount）
	TimeField    string   `json:"time_field"`    // 时间字段（用于趋势统计，如 created_at）
	ChartTypes   []string `json:"chart_types"`   // 图表类型: pie/bar/line
}

// GeneratorConfig 代码生成器配置
type GeneratorConfig struct {
	ID               uint             `json:"id,omitempty"`      // 配置ID（编辑时传递）
	TableName        string           `json:"table_name"`        // 数据库表名
	ModuleName       string           `json:"module_name"`       // 模块名称（英文，用于路由/文件名）
	Description      string           `json:"description"`       // 模块描述（中文）
	Author           string           `json:"author"`            // 作者
	GenerateBackend  bool             `json:"generate_backend"`  // 是否生成后端代码
	GenerateFrontend bool             `json:"generate_frontend"` // 是否生成前端代码
	GenerateSQL      bool             `json:"generate_sql"`      // 是否生成建表SQL
	FrontendPath     string           `json:"frontend_path"`     // 前端代码生成目录（留空则使用默认）
	Columns          []ColumnConfig   `json:"columns"`           // 字段配置
	Relations        []RelationConfig `json:"relations"`         // 关联关系配置
	MenuConfig       *MenuConfig      `json:"menu_config"`       // 菜单配置
	StatsConfig      *StatsConfig     `json:"stats_config"`      // 统计配置
	// 时间字段开关（用于决定是否在模型中包含这些字段）
	HasCreatedAt bool `json:"has_created_at"`
	HasUpdatedAt bool `json:"has_updated_at"`
	HasDeletedAt bool `json:"has_deleted_at"`
	HasCreatedBy bool `json:"has_created_by"` // 是否包含创建人字段
	// 创建人身份关联配置（显示创建者的身份信息，如医生/商家等）
	CreatedByProfileTable string `json:"created_by_profile_table"` // 创建者身份关联表名（如 doctor）
	CreatedByProfileField string `json:"created_by_profile_field"` // 身份表中要显示的字段（如 name）
	// 数据隔离配置（仅当 HasCreatedBy 为 true 时生效）
	DataIsolation bool   `json:"data_isolation"` // 是否启用数据隔离（非管理员只能看自己创建的数据）
	AdminRoleIds  string `json:"admin_role_ids"` // 管理员角色ID列表（逗号分隔，如 "1,2"）
	// 审批功能配置
	HasAudit bool `json:"has_audit"` // 是否启用审批功能（包含审批人、审批时间、审批备注字段）
	// 前台接口配置（仅当 HasCreatedBy 或 HasAudit 为 true 时有意义）
	GenerateFrontendApi bool `json:"generate_frontend_api"` // 是否生成前台用户使用的接口（私有路由，不做 created_by 过滤，仅返回已启用/审批通过的数据）
	// 用户关联配置
	LinkToUser      bool   `json:"link_to_user"`       // 是否关联用户表（一对一关系，生成 /my 接口供当前用户使用）
	ProfileName     string `json:"profile_name"`       // 身份显示名称（如"医生"、"商家"），默认使用 Description
	ProfileIcon     string `json:"profile_icon"`       // 身份图标（ant-design图标名，如"UserOutlined"）
	ProfileRoleCode string `json:"profile_role_code"` // 身份限定角色编码（为空表示不限制，任何有数据的用户都能看到）
}

// ColumnConfig 字段配置
type ColumnConfig struct {
	ColumnName    string `json:"column_name"`     // 数据库字段名
	FieldName     string `json:"field_name"`      // Go字段名（大驼峰）
	FieldType     string `json:"field_type"`      // Go类型
	JsonName      string `json:"json_name"`       // JSON名称（小写下划线）
	TsType        string `json:"ts_type"`         // TypeScript类型
	Comment       string `json:"comment"`         // 字段注释
	GormTag       string `json:"gorm_tag"`        // GORM标签
	DbType        string `json:"db_type"`         // 数据库字段类型（用于生成SQL）
	DbLength      int    `json:"db_length"`       // 数据库字段长度
	SqlType       string `json:"-"`               // SQL完整类型（自动生成，不序列化）
	DefaultValue  string `json:"default_value"`   // 默认值
	IsPrimaryKey  bool   `json:"is_primary_key"`  // 是否主键
	IsRequired    bool   `json:"is_required"`     // 是否必填
	IsSearchable  bool   `json:"is_searchable"`   // 是否可搜索
	SearchType    string `json:"search_type"`     // 搜索类型: eq(等于)/like(模糊)/gt(大于)/lt(小于)/between(范围)
	IsListVisible bool   `json:"is_list_visible"` // 是否列表显示
	IsFormVisible bool   `json:"is_form_visible"` // 是否表单显示
	IsSortable    bool   `json:"is_sortable"`     // 是否可排序
	SortOrder     string `json:"sort_order"`      // 默认排序方式: asc/desc
	IsUnique      bool   `json:"is_unique"`       // 是否唯一
	FormType      string `json:"form_type"`       // 表单组件类型
	DictType      string `json:"dict_type"`       // 字典类型
	// 下拉选项（form_type=select时使用）
	SelectOptions []SelectOption `json:"select_options"`
	// 开关值（form_type=switch时使用）
	SwitchValues *SwitchValue `json:"switch_values"`
}

// SelectOption 下拉选项
type SelectOption struct {
	Label string      `json:"label"` // 显示文本
	Value interface{} `json:"value"` // 值
}

// SwitchValue 开关值
type SwitchValue struct {
	ActiveValue   interface{} `json:"active_value"`   // 开启时的值
	InactiveValue interface{} `json:"inactive_value"` // 关闭时的值
	ActiveText    string      `json:"active_text"`    // 开启时的文本
	InactiveText  string      `json:"inactive_text"`  // 关闭时的文本
}

// RelationConfig 关联关系配置
type RelationConfig struct {
	RelationType  string `json:"relation_type"`   // 关联类型: hasOne, hasMany, belongsTo, many2many
	RelatedTable  string `json:"related_table"`   // 关联表名
	RelatedModule string `json:"related_module"`  // 关联模块名（用于 API import，留空则使用表名）
	RelatedModel  string `json:"related_model"`   // 关联模型名称
	ForeignKey    string `json:"foreign_key"`     // 外键字段
	ReferenceKey  string `json:"reference_key"`   // 引用字段
	JoinTable     string `json:"join_table"`      // 中间表(多对多时使用)
	DisplayField  string `json:"display_field"`   // 显示字段(如name)，用于下拉框和表格显示
	Comment       string `json:"comment"`         // 关联字段注释（如"产品类型"）
	IsRequired    bool   `json:"is_required"`     // 是否必填
	UseOptionsApi bool   `json:"use_options_api"` // 使用轻量options接口（返回id,name,count）而非分页列表
	UseTreeLayout bool   `json:"use_tree_layout"` // 使用左树右表布局（仅belongsTo生效）
}

// OptionItem 选项项（用于下拉选择）
type OptionItem struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// MenuConfig 菜单配置
type MenuConfig struct {
	ParentID   uint   `json:"parent_id"`  // 父菜单ID
	MenuName   string `json:"menu_name"`  // 菜单名称
	MenuIcon   string `json:"menu_icon"`  // 菜单图标
	MenuSort   int    `json:"menu_sort"`  // 菜单排序
	Permission string `json:"permission"` // 权限标识前缀
}

// GeneratedFile 生成的文件
type GeneratedFile struct {
	Path    string `json:"path"`    // 文件路径
	Content string `json:"content"` // 文件内容
	Type    string `json:"type"`    // 文件类型: backend/frontend
}

// PreviewResult 预览结果
type PreviewResult struct {
	Files []GeneratedFile `json:"files"`
}

// TableInfo 数据库表信息
type TableInfo struct {
	TableName    string `json:"table_name"`
	TableComment string `json:"table_comment"`
}

// ColumnInfo 数据库字段信息
type ColumnInfo struct {
	ColumnName    string  `json:"column_name"`
	DataType      string  `json:"data_type"`
	ColumnType    string  `json:"column_type"`
	ColumnComment string  `json:"column_comment"`
	ColumnKey     string  `json:"column_key"`
	IsNullable    string  `json:"is_nullable"`
	ColumnDefault *string `json:"column_default"`
}

// FormTypes 表单组件类型
var FormTypes = map[string]string{
	"input":    "输入框",
	"textarea": "文本域",
	"number":   "数字输入",
	"select":   "下拉选择",
	"radio":    "单选框",
	"checkbox": "复选框",
	"switch":   "开关",
	"date":     "日期选择",
	"datetime": "日期时间",
	// 上传与图片
	"image":  "图片上传(单个)",
	"images": "图片上传(多个)",
	"file":   "文件上传",
	"upload": "文件上传(兼容)",
	"editor": "富文本编辑器",
}
