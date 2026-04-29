package parser

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// columnInfo 从 information_schema.COLUMNS 查询到的字段信息
type columnInfo struct {
	ColumnName    string
	DataType      string
	ColumnType    string
	IsNullable    string
	ColumnKey     string
	ColumnDefault sql.NullString
	ColumnComment string
	CharMaxLength sql.NullInt64
	Extra         string
}

// Parser 数据库表结构解析器
type Parser struct {
	DSN                string   // "user:pass@tcp(host:port)/dbname"
	SkipFields         []string // 额外隐藏的字段列表（从 codegen.yaml skip_fields 加载）
	db                 *sql.DB  // 复用数据库连接
	dbName             string   // 缓存的数据库名
	tableColumnsCache  map[string]map[string]struct{}
	tableColumnsLoader func(tableName string) (map[string]struct{}, error)
}

// New 创建解析器实例
func New(dsn string, skipFields ...[]string) (*Parser, error) {
	p := &Parser{
		DSN:               dsn,
		tableColumnsCache: make(map[string]map[string]struct{}),
	}
	if len(skipFields) > 0 {
		p.SkipFields = skipFields[0]
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}
	if _, err := db.Exec("SET NAMES utf8mb4"); err != nil {
		db.Close()
		return nil, fmt.Errorf("设置字符集失败: %w", err)
	}
	dbName, err := extractDBName(dsn)
	if err != nil {
		db.Close()
		return nil, err
	}
	p.db = db
	p.dbName = dbName
	p.tableColumnsLoader = func(tableName string) (map[string]struct{}, error) {
		return queryTableColumnSet(db, dbName, tableName)
	}
	return p, nil
}

// Close 释放数据库连接
func (p *Parser) Close() {
	if p.db != nil {
		p.db.Close()
	}
}

// ParseTable 解析单张表
func (p *Parser) ParseTable(tableName string) (*TableMeta, error) {
	db := p.db
	dbName := p.dbName

	// 查询表注释
	tableComment, err := queryTableComment(db, dbName, tableName)
	if err != nil {
		return nil, err
	}

	// 查询字段信息
	columns, err := queryColumns(db, dbName, tableName)
	if err != nil {
		return nil, err
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("表 %s 不存在或没有字段", tableName)
	}

	identity := splitTableIdentity(tableName)
	meta := buildTableMetaSkeleton(identity, tableComment)
	appendColumnFields(meta, columns, buildExtraHiddenFieldSet(p.SkipFields))
	if err := p.resolveReferenceFields(meta, identity); err != nil {
		return nil, err
	}

	FinalizeTemplateMeta(meta)
	if err := validateRequiredScopeFields(meta); err != nil {
		return nil, err
	}

	return meta, nil
}

func validateRequiredScopeFields(meta *TableMeta) error {
	if meta == nil {
		return fmt.Errorf("表元数据不能为空")
	}

	var missing []string
	if !meta.HasTenantID {
		missing = append(missing, "tenant_id")
	}
	if !meta.HasMerchantID {
		missing = append(missing, "merchant_id")
	}
	if !meta.HasCreatedBy {
		missing = append(missing, "created_by")
	}
	if !meta.HasDeptID {
		missing = append(missing, "dept_id")
	}
	if len(missing) > 0 {
		return fmt.Errorf(
			"表 %s 缺少必需的权限归属字段: %s。请在迁移中补齐 tenant_id BIGINT UNSIGNED、merchant_id BIGINT UNSIGNED、created_by BIGINT UNSIGNED 和 dept_id BIGINT UNSIGNED",
			meta.TableName,
			strings.Join(missing, ", "),
		)
	}
	return nil
}

// FinalizeTemplateMeta 统一收口模板依赖的派生元数据。
// 真实解析、离线验证和测试应共享这套逻辑，避免多处手工拼装后出现行为漂移。
func FinalizeTemplateMeta(meta *TableMeta) {
	if meta == nil {
		return
	}

	meta.HasParentID = false
	meta.HasStatus = false
	meta.HasSort = false
	meta.HasPassword = false
	meta.HasTooltip = false
	meta.HasRichText = false
	meta.HasRichTextComponent = false
	meta.HasMoney = false
	meta.HasSearchable = false
	meta.HasTreeSelect = false
	meta.HasCreatedBy = false
	meta.HasDeptID = false
	meta.HasTenantID = false
	meta.HasMerchantID = false
	meta.HasTenantScope = false
	meta.HasDict = false
	meta.HasBatchEdit = false
	meta.HasImport = false
	meta.HasEnum = false
	meta.HasImage = false
	meta.HasForeignKey = false
	meta.HasKeywordSearch = false
	meta.ParentDisplayField = ""
	meta.SearchFields = nil
	meta.KeywordSearchFields = nil
	hasBatchEditStatus := false

	for i := range meta.Fields {
		field := &meta.Fields[i]

		if field.SearchEnabled && (field.IsForeignKey || field.IsParentID) && field.RefIsTree && field.SearchModeHint != "select" {
			field.SearchComponent = "TreeSelect"
		}

		if field.Name == "parent_id" {
			meta.HasParentID = true
			if field.RefDisplayLower != "" {
				meta.ParentDisplayField = field.RefDisplayLower
			}
		}
		if field.Name == "status" && !field.IsHidden {
			meta.HasStatus = true
		}
		if field.Name == "sort" && !field.IsHidden {
			meta.HasSort = true
		}
		if field.IsPassword {
			meta.HasPassword = true
		}
		if field.TooltipText != "" && !field.IsHidden && !field.IsID {
			meta.HasTooltip = true
		}
		if field.Component == "RichText" {
			meta.HasRichText = true
			meta.HasRichTextComponent = true
		}
		if field.Component == "JsonEditor" {
			meta.HasRichText = true
		}
		if field.IsMoney {
			meta.HasMoney = true
		}
		if field.Name == "created_by" {
			meta.HasCreatedBy = true
		}
		if field.Name == "dept_id" {
			meta.HasDeptID = true
		}
		if field.Name == "tenant_id" {
			meta.HasTenantID = true
		}
		if field.Name == "merchant_id" {
			meta.HasMerchantID = true
		}
		if field.DictType != "" {
			meta.HasDict = true
		}
		if field.IsSearchable {
			meta.HasSearchable = true
		}
		if field.Component == ComponentTreeSelectSingle || field.Component == ComponentTreeSelectMulti {
			meta.HasTreeSelect = true
		}
		if field.IsEnum && !field.IsHidden {
			meta.HasEnum = true
			if field.Name == "status" {
				hasBatchEditStatus = true
			}
		}
		if field.Component == ComponentImageUpload && !field.IsHidden {
			meta.HasImage = true
		}
		if field.IsForeignKey && !field.IsHidden {
			meta.HasForeignKey = true
		}
	}

	meta.HasBatchEdit = hasBatchEditStatus
	meta.HasImport = !meta.HasParentID
	meta.HasTenantScope = meta.HasTenantID

	finalizeSearchMeta(meta)
}

// findDisplayField 在关联表中按优先级查找显示字段。
// 关联表的列集合会按表缓存，避免对同一张表反复查询 information_schema。
func (p *Parser) findDisplayField(tableName string) string {
	columns, err := p.getTableColumns(tableName)
	if err != nil {
		return ""
	}
	for _, col := range displayFieldPriorityOrder() {
		if _, ok := columns[col]; ok {
			return col
		}
	}
	return ""
}

// tableHasColumn 检查表是否存在指定列。
func (p *Parser) tableHasColumn(tableName, columnName string) bool {
	columns, err := p.getTableColumns(tableName)
	if err != nil {
		return false
	}
	_, ok := columns[columnName]
	return ok
}

func (p *Parser) getTableColumns(tableName string) (map[string]struct{}, error) {
	if p.tableColumnsCache == nil {
		p.tableColumnsCache = make(map[string]map[string]struct{})
	}
	if columns, ok := p.tableColumnsCache[tableName]; ok {
		return columns, nil
	}

	loader := p.tableColumnsLoader
	if loader == nil {
		loader = func(tableName string) (map[string]struct{}, error) {
			return queryTableColumnSet(p.db, p.dbName, tableName)
		}
	}

	columns, err := loader(tableName)
	if err != nil {
		return nil, err
	}
	if columns == nil {
		columns = map[string]struct{}{}
	}
	p.tableColumnsCache[tableName] = columns
	return columns, nil
}

// ParseTables 解析多张表
func (p *Parser) ParseTables(tableNames []string) ([]*TableMeta, error) {
	var result []*TableMeta
	for _, name := range tableNames {
		meta, err := p.ParseTable(name)
		if err != nil {
			return nil, fmt.Errorf("解析表 %s 失败: %w", name, err)
		}
		result = append(result, meta)
	}
	return result, nil
}

// extractDBName 从 DSN 中提取数据库名
func extractDBName(dsn string) (string, error) {
	// DSN 格式: user:pass@tcp(host:port)/dbname?params
	slashIdx := strings.LastIndex(dsn, "/")
	if slashIdx < 0 {
		return "", fmt.Errorf("DSN 格式错误，无法提取数据库名: %s", dsn)
	}
	rest := dsn[slashIdx+1:]
	qIdx := strings.Index(rest, "?")
	if qIdx >= 0 {
		rest = rest[:qIdx]
	}
	if rest == "" {
		return "", fmt.Errorf("DSN 中未指定数据库名: %s", dsn)
	}
	return rest, nil
}

// queryTableComment 查询表注释
func queryTableComment(db *sql.DB, dbName, tableName string) (string, error) {
	var comment sql.NullString
	err := db.QueryRow(
		"SELECT TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?",
		dbName, tableName,
	).Scan(&comment)
	if err != nil {
		return "", fmt.Errorf("查询表 %s 注释失败: %w", tableName, err)
	}
	return comment.String, nil
}

// queryColumns 查询表的所有字段信息
func queryColumns(db *sql.DB, dbName, tableName string) ([]columnInfo, error) {
	rows, err := db.Query(
		`SELECT COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY,
		        COLUMN_DEFAULT, COLUMN_COMMENT, CHARACTER_MAXIMUM_LENGTH, EXTRA
		 FROM information_schema.COLUMNS
		 WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		 ORDER BY ORDINAL_POSITION`,
		dbName, tableName,
	)
	if err != nil {
		return nil, fmt.Errorf("查询表 %s 字段失败: %w", tableName, err)
	}
	defer rows.Close()

	var columns []columnInfo
	for rows.Next() {
		var col columnInfo
		if err := rows.Scan(
			&col.ColumnName, &col.DataType, &col.ColumnType, &col.IsNullable,
			&col.ColumnKey, &col.ColumnDefault, &col.ColumnComment,
			&col.CharMaxLength, &col.Extra,
		); err != nil {
			return nil, fmt.Errorf("扫描字段信息失败: %w", err)
		}
		columns = append(columns, col)
	}
	return columns, rows.Err()
}

func queryTableColumnSet(db *sql.DB, dbName, tableName string) (map[string]struct{}, error) {
	rows, err := db.Query(
		`SELECT COLUMN_NAME
		 FROM information_schema.COLUMNS
		 WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?`,
		dbName, tableName,
	)
	if err != nil {
		return nil, fmt.Errorf("查询表 %s 列集合失败: %w", tableName, err)
	}
	defer rows.Close()

	columns := make(map[string]struct{})
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, fmt.Errorf("扫描表 %s 列名失败: %w", tableName, err)
		}
		columns[columnName] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("读取表 %s 列集合失败: %w", tableName, err)
	}
	return columns, nil
}

// buildFieldMeta 根据列信息构建 FieldMeta
func buildFieldMeta(col columnInfo) FieldMeta {
	name := col.ColumnName
	isID := name == "id"
	// 外键判断：_id 后缀 + 排除特殊字段 + 必须是整数类型（varchar/char 类型的 _id 字段视为业务关联ID，非真正外键）
	isIntType := col.DataType == "bigint" || col.DataType == "int" || col.DataType == "smallint" || col.DataType == "tinyint" || col.DataType == "mediumint"
	isPassword := name == "password" || strings.HasSuffix(name, "_password") || strings.HasSuffix(name, "_pwd") ||
		strings.HasSuffix(name, "_secret") || strings.HasSuffix(name, "_secret_key") ||
		strings.HasSuffix(name, "_secret_id") || strings.HasSuffix(name, "_access_key") ||
		strings.HasSuffix(name, "_token")
	isForeignKey := strings.HasSuffix(name, "_id") && name != "id" && name != "dept_id" && name != "parent_id" && isIntType && !isPassword
	isMultiFK := strings.HasSuffix(name, "_ids")
	isParentID := name == "parent_id"

	// 解析备注
	commentMeta := ParseCommentMeta(col.ColumnComment)
	// 如果 comment 为空，回退为字段名的 CamelCase 形式作为 Label
	if commentMeta.Label == "" {
		commentMeta.Label = snakeToCamelDao(name)
		commentMeta.ShortLabel = commentMeta.Label
	}

	// 构建基础数据库类型（简化，去掉长度信息用于映射）
	dbType := col.ColumnType

	field := FieldMeta{
		Name:               name,
		NameCamel:          snakeToCamel(name),
		NameDao:            snakeToCamelDao(name),
		NameLower:          snakeToCamelLower(name),
		DBType:             dbType,
		GoType:             MapGoType(col.DataType, isID || isForeignKey || isParentID || name == "dept_id" || name == "created_by"),
		TSType:             MapTSType(col.DataType, isID || isForeignKey || isParentID || name == "dept_id" || name == "created_by"),
		Comment:            col.ColumnComment,
		Label:              commentMeta.Label,
		ShortLabel:         commentMeta.ShortLabel,
		TooltipText:        commentMeta.TooltipText,
		EnumValues:         commentMeta.EnumValues,
		IsRequired:         col.IsNullable == "NO" && !col.ColumnDefault.Valid && name != "id" && col.Extra != "auto_increment",
		IsID:               isID,
		IsParentID:         isParentID,
		IsForeignKey:       isForeignKey,
		IsMultiFK:          isMultiFK,
		IsTimeField:        strings.HasSuffix(name, "_at"),
		IsHidden:           IsHiddenField(name),
		IsEnum:             len(commentMeta.EnumValues) > 0,
		IsPassword:         isPassword,
		DictType:           commentMeta.DictType,
		DefaultValue:       col.ColumnDefault.String,
		SearchFormField:    snakeToCamelLower(name),
		SearchModeHint:     commentMeta.SearchMode,
		KeywordModeHint:    commentMeta.KeywordMode,
		SearchPriorityHint: commentMeta.SearchPriority,
		RefTableHint:       commentMeta.RefTableHint,
		RefDisplayHint:     commentMeta.RefDisplayHint,
	}
	applySystemScopeRefHint(&field)
	applySystemScopeFieldPresentation(&field)

	// 判断是否可搜索的文本字段（用于关键词模糊查询）
	goType := field.GoType
	if isSearchableTextFieldName(name) {
		if goType == "string" && !isID && !isForeignKey && !isPassword {
			field.IsSearchable = true
		}
	}

	// 判断是否精确搜索字段（编号类，用 = 而非 LIKE）
	if isExactSearchFieldName(name) {
		field.IsExactSearch = true
	}

	// 判断是否金额字段（单位：分，列表需要分→元格式化）
	if isMoneyFieldName(name) {
		field.IsMoney = true
	}

	if col.CharMaxLength.Valid {
		field.MaxLength = int(col.CharMaxLength.Int64)
	}

	if field.DictType != "" {
		field.EnumValues = nil
		field.IsEnum = false
	}

	// 自动推导验证规则
	field.ValidationRules, field.FrontendRules = buildValidationRules(field)
	// 更新时的验证规则（去掉 required）
	for _, r := range field.ValidationRules {
		if r != "required" {
			field.UpdateValidationRules = append(field.UpdateValidationRules, r)
		}
	}

	// 映射前端组件
	field.Component = MapComponent(field)

	// 隐藏字段不需要前端必填校验（后端自动注入）
	if field.IsHidden {
		field.IsRequired = false
	}

	applySearchMeta(&field)

	return field
}

func applySystemScopeRefHint(field *FieldMeta) {
	if field == nil || field.RefTableHint != "" {
		return
	}
	switch field.Name {
	case "tenant_id":
		field.RefTableHint = "system_tenant"
		field.RefDisplayHint = "name"
	case "merchant_id":
		field.RefTableHint = "system_merchant"
		field.RefDisplayHint = "name"
	}
}

func applySystemScopeFieldPresentation(field *FieldMeta) {
	if field == nil {
		return
	}
	switch field.Name {
	case "tenant_id":
		field.Label = "租户"
		field.ShortLabel = "租户"
		field.TooltipText = ""
	case "merchant_id":
		field.Label = "商户"
		field.ShortLabel = "商户"
		field.TooltipText = ""
	}
}

// buildValidationRules 根据字段名和类型自动推导验证规则
func buildValidationRules(f FieldMeta) (goRules []string, frontendRule string) {
	if f.IsID || f.IsHidden {
		return nil, ""
	}
	// 必填
	if f.IsRequired {
		goRules = append(goRules, "required")
	}
	// 邮箱
	if f.Name == "email" || strings.HasSuffix(f.Name, "_email") {
		goRules = append(goRules, "email")
		if f.IsRequired {
			frontendRule = "requiredEmail"
		} else {
			frontendRule = "email"
		}
	}
	// 手机号
	if f.Name == "phone" || f.Name == "mobile" || strings.HasSuffix(f.Name, "_phone") || strings.HasSuffix(f.Name, "_mobile") {
		goRules = append(goRules, "phone-loose")
		if f.IsRequired {
			frontendRule = "requiredPhone"
		} else {
			frontendRule = "phone"
		}
	}
	// URL
	if f.Name == "url" || strings.HasSuffix(f.Name, "_url") || strings.HasSuffix(f.Name, "_link") {
		goRules = append(goRules, "url")
		if f.IsRequired {
			frontendRule = "requiredUrl"
		} else {
			frontendRule = "url"
		}
	}
	// 长度限制（仅 string 类型且有 MaxLength）
	if f.GoType == "string" && f.MaxLength > 0 && !f.IsPassword {
		goRules = append(goRules, fmt.Sprintf("max-length:%d", f.MaxLength))
	}
	// 密码
	if f.IsPassword {
		goRules = append(goRules, "length:6,32")
	}
	return
}

func applySearchMeta(field *FieldMeta) {
	resetSearchMeta(field)
	field.SearchFormField = field.NameLower
	if field.IsHidden || field.IsID || field.IsPassword {
		return
	}

	switch field.SearchModeHint {
	case "off":
		field.SearchEnabled = false
	case "eq":
		field.SearchExplicit = true
		applyExplicitEqSearchMeta(field)
	case "like":
		field.SearchExplicit = true
		setInputSearch(field, "like")
	case "range":
		field.SearchExplicit = true
		setRangeSearch(field)
	case "select":
		field.SearchExplicit = true
		setSelectSearch(field)
	case "tree":
		field.SearchExplicit = true
		setTreeSearch(field)
	case "on":
		field.SearchExplicit = true
		applyHeuristicSearchMeta(field)
	default:
		applyHeuristicSearchMeta(field)
	}

	if field.SearchEnabled {
		field.SearchPriority = inferSearchPriority(*field)
	}
	if field.SearchPriorityHint > 0 {
		field.SearchPriority = field.SearchPriorityHint
	}
	applyKeywordMeta(field)
	if field.KeywordOnly {
		field.SearchEnabled = false
	}
}

func resetSearchMeta(field *FieldMeta) {
	field.SearchEnabled = false
	field.SearchComponent = ""
	field.SearchOperator = ""
	field.SearchGoType = ""
	field.SearchTSType = ""
	field.SearchPointer = false
	field.SearchRange = false
	field.SearchFormField = ""
	field.SearchPriority = 0
	field.SearchExplicit = false
	field.KeywordEnabled = false
	field.KeywordOnly = false
}

// ApplySearchMeta 对单个字段应用搜索启发式和注释覆盖规则。
// 离线验证等手工拼装 FieldMeta 的场景应复用这里，避免和真实 parser 逻辑漂移。
func ApplySearchMeta(field *FieldMeta) {
	if field == nil {
		return
	}
	applySearchMeta(field)
}

func shouldAutoSearchTextField(field FieldMeta) bool {
	if field.GoType != "string" {
		return false
	}
	switch field.Component {
	case ComponentImageUpload, ComponentFileUpload, ComponentRichText, ComponentJsonEditor, ComponentPassword:
		return false
	}
	if field.IsSearchable {
		return true
	}
	return isSearchableTextFieldName(field.Name) || isExactSearchFieldName(field.Name)
}

func applyHeuristicSearchMeta(field *FieldMeta) {
	switch {
	case field.IsTimeField:
		setRangeSearch(field)
	case field.IsForeignKey || field.IsParentID:
		if field.Component == ComponentTreeSelectSingle || field.Component == ComponentTreeSelectMulti {
			setTreeSearch(field)
		} else {
			setSelectSearch(field)
		}
	case field.IsEnum || field.DictType != "":
		setSelectSearch(field)
	case shouldAutoSearchTextField(*field):
		if field.IsExactSearch {
			setInputSearch(field, "eq")
		} else {
			setInputSearch(field, "like")
		}
	}
}

func applyExplicitEqSearchMeta(field *FieldMeta) {
	switch {
	case field.IsForeignKey || field.IsParentID:
		if field.Component == ComponentTreeSelectSingle || field.Component == ComponentTreeSelectMulti {
			setTreeSearch(field)
		} else {
			setSelectSearch(field)
		}
	case field.IsEnum || field.DictType != "":
		setSelectSearch(field)
	default:
		setInputSearch(field, "eq")
	}
}

func setRangeSearch(field *FieldMeta) {
	field.SearchEnabled = true
	field.SearchRange = true
	field.SearchComponent = "RangePicker"
	field.SearchOperator = "range"
	field.SearchGoType = "string"
	field.SearchTSType = "string"
	field.SearchPointer = false
	field.SearchFormField = field.NameLower + "Range"
}

func setSelectSearch(field *FieldMeta) {
	field.SearchEnabled = true
	field.SearchRange = false
	field.SearchComponent = "Select"
	field.SearchOperator = "eq"
	field.SearchFormField = field.NameLower
	if field.IsForeignKey || field.IsParentID {
		field.SearchGoType = "snowflake.JsonInt64"
		field.SearchTSType = "string"
		field.SearchPointer = true
		return
	}
	field.SearchGoType = field.GoType
	field.SearchTSType = field.TSType
	field.SearchPointer = true
}

func setTreeSearch(field *FieldMeta) {
	setSelectSearch(field)
	field.SearchComponent = "TreeSelect"
}

func setInputSearch(field *FieldMeta, operator string) {
	field.SearchEnabled = true
	field.SearchRange = false
	field.SearchComponent = "Input"
	field.SearchOperator = operator
	field.SearchGoType = "string"
	field.SearchTSType = "string"
	field.SearchPointer = false
	field.SearchFormField = field.NameLower
}

func applyKeywordMeta(field *FieldMeta) {
	switch field.KeywordModeHint {
	case "off":
		field.KeywordEnabled = false
		field.KeywordOnly = false
		return
	case "only":
		field.KeywordEnabled = true
		field.KeywordOnly = true
		return
	case "on":
		field.KeywordEnabled = true
		field.KeywordOnly = false
		return
	}
	field.KeywordEnabled = field.SearchEnabled && field.SearchOperator == "like"
}

func inferSearchPriority(field FieldMeta) int {
	switch {
	case field.IsExactSearch:
		return 100
	case field.Name == "title" || field.Name == "name" || field.Name == "username" || field.Name == "nickname" || field.Name == "real_name":
		return 95
	case field.Name == "phone" || field.Name == "mobile" || field.Name == "email" || strings.HasSuffix(field.Name, "_phone") || strings.HasSuffix(field.Name, "_mobile") || strings.HasSuffix(field.Name, "_email"):
		return 90
	case field.IsForeignKey || field.IsParentID:
		return 82
	case field.IsEnum || field.DictType != "":
		return 78
	case field.IsTimeField:
		return 60
	case field.SearchOperator == "like":
		return 70
	default:
		return 50
	}
}

func finalizeSearchMeta(meta *TableMeta) {
	type searchItem struct {
		index int
		field FieldMeta
	}

	sortItems := func(items []searchItem) {
		sort.SliceStable(items, func(i, j int) bool {
			if items[i].field.SearchPriority == items[j].field.SearchPriority {
				return items[i].index < items[j].index
			}
			return items[i].field.SearchPriority > items[j].field.SearchPriority
		})
	}

	direct := make([]searchItem, 0)
	keyword := make([]searchItem, 0)
	for i, field := range meta.Fields {
		if field.SearchEnabled {
			direct = append(direct, searchItem{index: i, field: field})
		}
		if field.KeywordEnabled {
			keyword = append(keyword, searchItem{index: i, field: field})
		}
	}

	sortItems(direct)
	sortItems(keyword)

	visible := make([]FieldMeta, 0, len(direct))
	fuzzyDirectCount := 0
	for _, item := range direct {
		field := item.field
		if field.KeywordOnly {
			continue
		}
		if field.SearchOperator == "like" && field.KeywordEnabled && !field.SearchExplicit {
			if fuzzyDirectCount >= 2 {
				continue
			}
			fuzzyDirectCount++
		}
		visible = append(visible, field)
	}

	meta.SearchFields = visible

	if len(keyword) == 0 {
		return
	}

	keywordFields := make([]FieldMeta, 0, len(keyword))
	keywordSeen := make(map[string]struct{}, len(keyword))
	for _, item := range keyword {
		if _, ok := keywordSeen[item.field.Name]; ok {
			continue
		}
		keywordSeen[item.field.Name] = struct{}{}
		keywordFields = append(keywordFields, item.field)
	}
	meta.KeywordSearchFields = keywordFields
	meta.HasKeywordSearch = len(keywordFields) > 1
	for _, field := range keywordFields {
		if field.KeywordOnly {
			meta.HasKeywordSearch = true
			break
		}
	}
}
