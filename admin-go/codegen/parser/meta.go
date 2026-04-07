package parser

// 前端组件类型常量
const (
	ComponentInput            = "Input"
	ComponentInputNumber      = "InputNumber"
	ComponentTextarea         = "Textarea"
	ComponentSwitch           = "Switch"
	ComponentRadio            = "Radio"
	ComponentSelect           = "Select"
	ComponentTreeSelectSingle = "TreeSelectSingle"
	ComponentTreeSelectMulti  = "TreeSelectMulti"
	ComponentSelectMulti      = "SelectMulti"
	ComponentImageUpload      = "ImageUpload"
	ComponentFileUpload       = "FileUpload"
	ComponentRichText         = "RichText"
	ComponentJsonEditor       = "JsonEditor"
	ComponentPassword         = "Password"
	ComponentInputUrl         = "InputUrl"
	ComponentDateTimePicker   = "DateTimePicker"
	ComponentIconPicker       = "IconPicker"
)

// EnumValue 枚举值
type EnumValue struct {
	Value     string
	Label     string
	NameIdent string // 语义化标识符，如 "Enabled"/"Disabled"（为空则用 V+Value 兜底）
}

// FieldMeta 字段元数据
type FieldMeta struct {
	Name                  string      // snake_case
	NameCamel             string      // CamelCase（Go 风格，ID/URL 全大写）
	NameDao               string      // CamelCase（GoFrame DAO 风格，Id/Url 首字母大写）
	NameLower             string      // camelCase（首字母小写）
	DBType                string      // varchar/int/bigint/tinyint/text 等
	GoType                string      // string/int/int64/JsonInt64 等
	TSType                string      // string/number/boolean
	Comment               string      // 原始备注
	Label                 string      // 前端 Label（完整，含括号）
	ShortLabel            string      // 精简标签（去掉括号部分，用于列头等紧凑场景）
	TooltipText           string      // 括号内的提示文字（为空则无需 Tooltip）
	EnumValues            []EnumValue // 枚举值列表
	Component             string      // 前端组件类型
	IsRequired            bool
	IsID                  bool   // 是否是主键 id
	IsParentID            bool   // 是否是 parent_id
	IsForeignKey          bool   // 是否是 *_id（单选外键）
	IsMultiFK             bool   // 是否是 *_ids（多选外键）
	IsTimeField           bool   // 是否是 *_at 时间字段
	IsHidden              bool   // 表单中隐藏（id/created_at/updated_at/deleted_at/created_by/dept_id）
	IsEnum                bool   // 是否有枚举值
	IsPassword            bool   // 是否是密码字段
	IsSearchable          bool   // 是否可用于关键词搜索（title/name/nickname/phone/email/order_no 等文本字段）
	IsExactSearch         bool   // 是否精确搜索（编号类字段 _no/_code/_sn，用 = 而非 LIKE）
	SearchEnabled         bool   // 是否进入列表搜索区
	SearchComponent       string // 搜索控件类型：Input/Select/TreeSelect/RangePicker
	SearchOperator        string // 搜索操作符：like/eq/range
	SearchGoType          string // 搜索参数 Go 类型（基础类型）
	SearchTSType          string // 搜索参数 TS 类型
	SearchPointer         bool   // 搜索参数是否使用指针（区分未传与零值）
	SearchRange           bool   // 是否为区间搜索（生成 start/end 参数）
	SearchFormField       string // 前端搜索表单字段名
	SearchPriority        int    // 搜索优先级，值越大越靠前
	SearchExplicit        bool   // 是否通过注释显式声明了搜索模式
	SearchModeHint        string // 注释中的搜索模式提示：off/on/eq/like/range/select/tree
	KeywordModeHint       string // 注释中的关键词模式：on/off/only
	KeywordEnabled        bool   // 是否参与全局关键词搜索
	KeywordOnly           bool   // 是否仅参与全局关键词搜索，不单独生成控件
	SearchPriorityHint    int    // 注释中的优先级提示
	IsMoney               bool   // 是否是金额字段（*_price/*_amount/*_balance/*_income，单位：分）
	MaxLength             int
	ValidationRules       []string // 后端验证规则列表，如 ["required", "email", "length:1,50"]
	UpdateValidationRules []string // 更新时的验证规则（去掉 required）
	FrontendRules         string   // 前端验证规则标识：email/phone/url/''
	DictType              string   // 字典表类型标识，如 "gender"（非空表示使用字典表动态加载）
	DefaultValue          string
	RefTableHint          string // 注释显式指定的关联表，如 system_users
	RefDisplayHint        string // 注释显式指定的显示字段，如 username
	// 关联字段信息（仅 IsForeignKey 或 IsParentID 时有值）
	RefTable        string // 关联模块名，如 dept（用于 dao 引用和代码生成）
	RefTableDB      string // 关联表实际数据库表名，如 system_dept（用于 g.DB().Model()）
	RefTableApp     string // 关联表的应用名，如 system（用于前端 API 导入路径）
	RefTableCamel   string // 关联表 CamelCase，如 Article
	RefTableLower   string // 关联表 camelCase，如 article
	RefDisplayField string // 关联表显示字段 snake_case，如 title
	RefDisplayCamel string // 关联显示字段 CamelCase，如 Title
	RefDisplayLower string // 关联显示字段 camelCase，如 title
	RefFieldName    string // 结构体字段名 = RefTableCamel + RefDisplayCamel，如 ArticleTitle
	RefFieldJSON    string // json 名 = RefTableLower + RefDisplayCamel，如 articleTitle
	RefIsTree       bool   // 关联表是否有 parent_id（树形结构）
}

// TableMeta 表元数据
type TableMeta struct {
	TableName           string
	AppName             string // 应用名，如 system、demo
	AppNameCamel        string // 应用名 CamelCase，如 System、Demo
	ModelName           string // CamelCase（模块名），如 Dept（用于 service/model/controller 命名）
	DaoName             string // CamelCase（完整表名），如 DemoDemo（用于 dao 引用，gf gen dao 生成的名称）
	ModuleName          string // 小写，如 dept
	PackageName         string // 包名，如 dept
	Comment             string
	Fields              []FieldMeta
	HasParentID         bool   // 有 parent_id 字段
	HasStatus           bool   // 有 status 字段
	HasSort             bool   // 有 sort 字段
	HasPassword         bool   // 有 password 字段
	HasTooltip          bool   // 有字段需要 Tooltip 提示
	HasRichText         bool   // 有 RichText 或 JsonEditor 字段（用于弹窗加宽）
	HasMoney            bool   // 有金额字段（用于列表格式化）
	HasSearchable       bool   // 有可搜索的文本字段
	HasTreeSelect       bool   // 有 TreeSelectSingle/TreeSelectMulti 字段（不含外键 TreeSelect）
	HasCreatedBy        bool   // 有 created_by 字段（用于数据权限注入）
	HasDeptID           bool   // 有 dept_id 字段（用于数据权限注入）
	HasDict             bool   // 有字典字段（需要导入字典 API）
	HasBatchEdit        bool   // 有可批量编辑的枚举字段（status 等）
	HasImport           bool   // 是否生成导入功能（默认 true，除树形表外）
	HasEnum             bool   // 有非隐藏的枚举字段（用于前端 Tag 导入判断）
	HasImage            bool   // 有图片上传字段（用于前端列表图片 slot）
	HasForeignKey       bool   // 有非隐藏的外键字段（用于前端 form 外键选项加载）
	HasKeywordSearch    bool   // 是否启用全局关键词搜索
	ParentDisplayField  string // parent_id 关联的显示字段 camelCase（如 title/name），用于前端 TreeSelect fieldNames
	EnableOpLog         bool   // 是否生成操作日志（由配置控制）
	SearchFields        []FieldMeta
	KeywordSearchFields []FieldMeta
}
