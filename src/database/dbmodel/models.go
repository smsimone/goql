package dbmodel

import (
	"database/sql"
)

type Table struct {
	SchemaName  *string `json:"schemaname" db:"schemaname"`
	TableName   *string `json:"tablename" db:"tablename"`
	TableOwner  *string `json:"tableowner" db:"tableowner"`
	TableSpace  *string `json:"tablespace" db:"tablespace"`
	HasIndexes  *bool   `json:"hasindexes" db:"hasindexes"`
	HasRules    *bool   `json:"hasrules" db:"hasrules"`
	HasTriggers *bool   `json:"hastriggers" db:"hastriggers"`
	RowSecurity *bool   `json:"rowsecurity" db:"rowsecurity"`
}

type TableInformation struct {
	Columns []ColumnInformation `json:"columns"`
}

type ColumnInformation struct {
	ColumnName   *string `json:"columnname" db:"columnname"`
	DataType     *string `json:"datatype" db:"data_type"`
	IsNullable   *string `json:"isnullable" db:"is_nullable"`
	DefaultValue *string `json:"column_default" db:"column_default"`
}

func (t *Table) FromRow(row *sql.Rows) error {
	var schemaName string
	var tableName string

	if err := row.Scan(&schemaName, &tableName); err != nil {
		return err
	}

	t.SchemaName = &schemaName
	t.TableName = &tableName

	return nil
}
