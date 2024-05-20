package database

import (
	"fmt"
	"goql/src/database/dbmodel"
	"strings"
)

type ConstraintType string

const (
	PKEY ConstraintType = "primary"
	FKEY ConstraintType = "foreign"
)

type ActiveConnection interface {
	GetId() int
	Disconnect()
	GetTables() (*[]dbmodel.Table, error)
	GetTableInformation(schema string, table string)
	GetTableData(schema string, table string) (*TableData, error)
	UpdateValue(request UpdateValueRequest) error
}

type UpdateValueRequest struct {
	Schema string `json:"schema"`
	Table  string `json:"table"`
	Value  any    `json:"value"`
	Field  string `json:"field"`
	/// must be a query
	Where string `json:"where"`
}

type Constraint struct {
	Column string         `json:"column"`
	Type   ConstraintType `json:"type"`
}

type TableData struct {
	Columns    []string  `json:"columns"`
	Rows       []RowData `json:"rows"`
	PrimaryKey []string  `json:"primay_key"`
}

type RowData struct {
	Columns []Column `json:"columns"`
}

type Column struct {
	Value    any    `json:"value"`
	DataType string `json:"data_type"`
}

func ConstraintFromName(name string) (ConstraintType, error) {
	if strings.HasSuffix(name, "pkey") {
		return PKEY, nil
	} else if strings.HasSuffix(name, "fkey") {
		return FKEY, nil
	}
	return PKEY, fmt.Errorf("type %s does not exists", name)
}
