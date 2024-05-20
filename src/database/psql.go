package database

import (
	"context"
	"database/sql"
	"fmt"
	"goql/src/database/dbmodel"
	"reflect"
)

type PostgresConnection struct {
	Id   int
	Conn *sql.Conn
}

func (p *PostgresConnection) GetId() int {
	return p.Id
}

func (p *PostgresConnection) Disconnect() {
	if p.Conn != nil {
		p.Conn.Close()
	}
}

func (p *PostgresConnection) GetTables() (*[]dbmodel.Table, error) {
	conn := p.Conn
	if conn == nil {
		fmt.Println("no active connection")
		return nil, fmt.Errorf("no active connection")
	}

	rows, err := conn.QueryContext(context.Background(), "SELECT schemaname, tablename FROM pg_catalog.pg_tables;")
	if err != nil {
		fmt.Println("failed to get tables")
		return nil, fmt.Errorf("failed to get tables")
	}

	fmt.Println("Fetched databse tables")

	tables := []dbmodel.Table{}

	for rows.Next() {
		table := dbmodel.Table{}
		err := table.FromRow(rows)
		if err != nil {
			fmt.Printf("failed to scan table: %s\n", err.Error())
			continue
		}
		tables = append(tables, table)
	}

	fmt.Printf("Converted %d tables\n", len(tables))

	return &tables, nil
}

func (p *PostgresConnection) GetTableInformation(schema string, table string) {

}

func (p *PostgresConnection) getKeys(table string) (*[]Constraint, error) {
	query := fmt.Sprintf("select constraint_name,column_name from information_schema.key_column_usage where table_name = '%s'", table)
	fmt.Printf("Get key query: %s\n", query)

	rows, err := p.Conn.QueryContext(context.Background(), query)
	if err != nil {
		fmt.Printf("Failed to get keys for table %s: %s\n", table, err.Error())
		return nil, err
	}

	constraints := []Constraint{}
	for rows.Next() {
		var (
			name   string
			column string
		)

		if err := rows.Scan(&name, &column); err != nil {
			fmt.Printf("Failed to parse constraint rows: %s\n", err.Error())
			return nil, err
		}

		constr, err := ConstraintFromName(name)
		if err != nil {
			fmt.Printf("Failed to get constraint from %s: %s\n", name, err.Error())
			continue
		}

		constraints = append(constraints, Constraint{
			Column: column,
			Type:   constr,
		})
	}
	defer rows.Close()

	return &constraints, nil
}

func (p *PostgresConnection) GetTableData(schema string, table string) (*TableData, error) {
	constr, err := p.getKeys(table)
	if err != nil {
		return nil, err
	}

	rows, err := p.Conn.QueryContext(context.Background(), fmt.Sprintf("SELECT * from %s.\"%s\";", schema, table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var t TableData
	t.Columns = cols
	t.PrimaryKey = []string{}
	for _, val := range *constr {
		if val.Type == PKEY {
			t.PrimaryKey = append(t.PrimaryKey, val.Column)
		}
	}

	for rows.Next() {
		var columnTypes []reflect.Type
		for _, c := range types {
			rowT := c.ScanType()
			columnTypes = append(columnTypes, rowT)
		}

		values := []any{}
		for range columnTypes {
			values = append(values, &[]byte{})
		}

		err := rows.Scan(values...)
		if err != nil {
			fmt.Printf("Failed to scan row: %s\n", err.Error())
			continue
		}

		currentRow := RowData{
			Columns: []Column{},
		}

		for idx, item := range columnTypes {
			var col Column
			col.DataType = item.Name()
			rawData := values[idx]
			if rawData != nil {
				bytes := rawData.(*[]byte)
				col.Value = string(*bytes)
			}
			fmt.Printf("Row data: %+v\n", col)
			currentRow.Columns = append(currentRow.Columns, col)
		}
		t.Rows = append(t.Rows, currentRow)
	}

	return &t, nil
}

func (p *PostgresConnection) UpdateValue(request UpdateValueRequest) error {
	query := fmt.Sprintf("UPDATE %s.\"%s\" SET \"%s\"=%s WHERE %s;", request.Schema, request.Table, request.Field, request.Value, request.Where)

	_, err := p.Conn.ExecContext(context.Background(), query)

	if err != nil {
		fmt.Printf("Failed to update field: %s\n", err.Error())
		return err
	}

	return nil
}
