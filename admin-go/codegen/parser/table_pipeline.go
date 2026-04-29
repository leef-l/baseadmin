package parser

import (
	"fmt"
	"strings"
)

type tableIdentity struct {
	tableName  string
	appName    string
	moduleName string
}

type referenceResolution struct {
	refTable        string
	refTableDB      string
	refTableApp     string
	refTableCamel   string
	refTableLower   string
	displayField    string
	displayCamel    string
	displayLower    string
	refFieldName    string
	refFieldJSON    string
	refIsTree        bool
	refHasDeletedAt  bool
	refHasTenantID   bool
	refHasMerchantID bool
	candidateTables string
}

func splitTableIdentity(tableName string) tableIdentity {
	appName := ""
	moduleName := tableName
	if idx := strings.Index(tableName, "_"); idx > 0 {
		appName = tableName[:idx]
		moduleName = tableName[idx+1:]
	}
	return tableIdentity{
		tableName:  tableName,
		appName:    appName,
		moduleName: moduleName,
	}
}

func buildTableMetaSkeleton(identity tableIdentity, tableComment string) *TableMeta {
	return &TableMeta{
		TableName:    identity.tableName,
		AppName:      identity.appName,
		AppNameCamel: snakeToCamel(identity.appName),
		ModelName:    snakeToCamel(identity.moduleName),
		DaoName:      snakeToCamel(identity.tableName),
		ModuleName:   strings.ToLower(identity.moduleName),
		PackageName:  strings.ToLower(identity.moduleName),
		Comment:      tableComment,
	}
}

func buildExtraHiddenFieldSet(skipFields []string) map[string]struct{} {
	fields := make(map[string]struct{}, len(skipFields))
	for _, field := range skipFields {
		fields[field] = struct{}{}
	}
	return fields
}

func appendColumnFields(meta *TableMeta, columns []columnInfo, extraHidden map[string]struct{}) {
	for _, col := range columns {
		field := buildFieldMeta(col)
		if _, ok := extraHidden[field.Name]; ok {
			field.IsHidden = true
		}
		meta.Fields = append(meta.Fields, field)
	}
}

func (p *Parser) resolveReferenceFields(meta *TableMeta, identity tableIdentity) error {
	for i := range meta.Fields {
		field := &meta.Fields[i]
		if !field.IsForeignKey && !field.IsParentID {
			continue
		}

		resolution, err := p.resolveReferenceTarget(field, identity)
		if err != nil {
			return err
		}
		if referenceFieldCollides(meta.Fields, resolution.refFieldName) {
			// 用字段名前缀消歧：author_id → AuthorUsername
			prefix := snakeToCamel(strings.TrimSuffix(field.Name, "_id"))
			resolution.refFieldName = prefix + resolution.displayCamel
			resolution.refFieldJSON = snakeToCamelLower(strings.TrimSuffix(field.Name, "_id")) + resolution.displayCamel
			if referenceFieldCollides(meta.Fields, resolution.refFieldName) {
				continue
			}
		}
		applyReferenceResolution(field, resolution)
	}
	return nil
}

func (p *Parser) resolveReferenceTarget(field *FieldMeta, identity tableIdentity) (referenceResolution, error) {
	if field == nil {
		return referenceResolution{}, fmt.Errorf("关联字段不能为空")
	}

	refTable := strings.TrimSuffix(field.Name, "_id")
	if field.IsParentID {
		refTable = identity.moduleName
	}

	resolution := referenceResolution{
		refTable:        refTable,
		refTableDB:      refTable,
		candidateTables: refTable,
	}

	if field.RefTableHint != "" {
		resolution.refTableDB = field.RefTableHint
		resolution.candidateTables = field.RefTableHint
		if idx := strings.Index(resolution.refTableDB, "_"); idx > 0 {
			resolution.refTable = resolution.refTableDB[idx+1:]
		} else {
			resolution.refTable = resolution.refTableDB
		}
		resolution.displayField = field.RefDisplayHint
		if resolution.displayField == "" {
			resolution.displayField = p.findDisplayField(resolution.refTableDB)
		}
	} else {
		if identity.appName != "" {
			prefixed := identity.appName + "_" + resolution.refTable
			resolution.candidateTables = prefixed + " 或 " + resolution.refTable
			if displayField := p.findDisplayField(prefixed); displayField != "" {
				resolution.refTableDB = prefixed
				resolution.displayField = displayField
			}
		}
		if resolution.displayField == "" {
			resolution.displayField = p.findDisplayField(resolution.refTable)
		}
	}

	if resolution.displayField == "" {
		return referenceResolution{}, fmt.Errorf(
			"字段 %s 是外键，但找不到关联表（尝试了 %s）。\n  请先创建关联表，或将字段名改为非 _id 后缀",
			field.Name, resolution.candidateTables,
		)
	}

	if idx := strings.Index(resolution.refTableDB, "_"); idx > 0 {
		resolution.refTableApp = resolution.refTableDB[:idx]
	} else {
		resolution.refTableApp = identity.appName
	}
	resolution.refTableCamel = snakeToCamel(resolution.refTable)
	resolution.refTableLower = snakeToCamelLower(resolution.refTable)
	resolution.displayCamel = snakeToCamel(resolution.displayField)
	resolution.displayLower = snakeToCamelLower(resolution.displayField)
	resolution.refFieldName = resolution.refTableCamel + resolution.displayCamel
	resolution.refFieldJSON = resolution.refTableLower + resolution.displayCamel
	resolution.refIsTree = p.tableHasColumn(resolution.refTableDB, "parent_id")
	resolution.refHasDeletedAt = p.tableHasColumn(resolution.refTableDB, "deleted_at")
	resolution.refHasTenantID = p.tableHasColumn(resolution.refTableDB, "tenant_id")
	resolution.refHasMerchantID = p.tableHasColumn(resolution.refTableDB, "merchant_id")
	return resolution, nil
}

func referenceFieldCollides(fields []FieldMeta, refFieldName string) bool {
	for _, other := range fields {
		if other.NameCamel == refFieldName || other.RefFieldName == refFieldName {
			return true
		}
	}
	return false
}

func applyReferenceResolution(field *FieldMeta, resolution referenceResolution) {
	field.RefTable = resolution.refTable
	field.RefTableDB = resolution.refTableDB
	field.RefTableApp = resolution.refTableApp
	field.RefTableCamel = resolution.refTableCamel
	field.RefTableLower = resolution.refTableLower
	field.RefDisplayField = resolution.displayField
	field.RefDisplayCamel = resolution.displayCamel
	field.RefDisplayLower = resolution.displayLower
	field.RefFieldName = resolution.refFieldName
	field.RefFieldJSON = resolution.refFieldJSON
	field.RefIsTree = resolution.refIsTree
	field.RefHasDeletedAt = resolution.refHasDeletedAt
	field.RefHasTenantID = resolution.refHasTenantID
	field.RefHasMerchantID = resolution.refHasMerchantID
}
