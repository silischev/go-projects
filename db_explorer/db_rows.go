package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type dbTblAttrValue struct {
	AttrName string
	Val      interface{}
}

type dbTuple struct {
	Value []dbTblAttrValue
}

func getItem(db *sql.DB, table string, id int, columns []dbColumn) (dbTuple, error) {
	tuple := dbTuple{}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE id = ?", table), id)
	if err != nil {
		return tuple, err
	}

	defer rows.Close()

	cols, _ := rows.Columns()
	var dbTblRs []dbTuple

	for rows.Next() {
		dbTblRs = append(dbTblRs, getTuple(cols, rows, columns))
	}

	if len(dbTblRs) > 0 {
		tuple = dbTblRs[0]
	}

	return tuple, nil
}

func getRows(db *sql.DB, table string, columns []dbColumn, limit int, offset int) ([]dbTuple, error) {
	var rows *sql.Rows
	var err error

	query := fmt.Sprintf("SELECT * FROM %s", table)

	if limit > 0 {
		rows, err = db.Query(query+" LIMIT ? OFFSET ?", limit, offset)
	} else {
		rows, err = db.Query(query)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cols, _ := rows.Columns()
	var dbTblRs []dbTuple

	for rows.Next() {
		dbTblRs = append(dbTblRs, getTuple(cols, rows, columns))
	}

	return dbTblRs, nil
}

func createItem(db *sql.DB, table string, columns []dbColumn, data map[string]interface{}) ([]dbTuple, error) {
	var dbTblRs []dbTuple
	var columnsNames []string
	var placeholders []string
	var values []interface{}

	for _, column := range columns {
		if column.ColumnKey.String == PrimaryKey {
			continue
		}

		if val, ok := data[column.Name]; ok {
			columnsNames = append(columnsNames, column.Name)
			placeholders = append(placeholders, "?")
			values = append(values, val)
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columnsNames, ","), strings.Join(placeholders, ","))

	result, err := db.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	row := dbTblAttrValue{"id", id}
	dbTuple := dbTuple{[]dbTblAttrValue{row}}
	dbTblRs = append(dbTblRs, dbTuple)

	return dbTblRs, nil
}

func updateItem(db *sql.DB, table string, id int, columns []dbColumn, data map[string]interface{}) ([]dbTuple, error) {
	var dbTblRs []dbTuple
	var columnsNames []string
	var values []interface{}

	for _, column := range columns {
		if column.ColumnKey.String == PrimaryKey {
			continue
		}

		if val, ok := data[column.Name]; ok {
			columnsNames = append(columnsNames, column.Name+" = ?")
			values = append(values, val)
		}
	}

	if len(columnsNames) == 0 && data["id"] != nil {
		return nil, NewHttpError("field id have invalid type", http.StatusBadRequest)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %v", table, strings.Join(columnsNames, ","), id)

	result, err := db.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	row := dbTblAttrValue{"updated", rowsAffected}
	dbTuple := dbTuple{[]dbTblAttrValue{row}}
	dbTblRs = append(dbTblRs, dbTuple)

	return dbTblRs, nil
}

func getTuple(cols []string, rows *sql.Rows, columns []dbColumn) dbTuple {
	dbTuple := dbTuple{}
	values := make([]interface{}, len(cols))
	pointers := make([]interface{}, len(cols))

	for i := range values {
		pointers[i] = &values[i]
	}

	rows.Scan(pointers...)

	for i, colName := range cols {
		for _, column := range columns {
			if colName == column.Name {
				var commonVal interface{}
				var strVal string

				byteVal, ok := values[i].([]byte)

				if ok {
					switch column.Type {
					case "int":
						strVal = string(byteVal)
						commonVal, _ = strconv.ParseInt(strVal, 10, 64)
					case "varchar", "text":
						commonVal = string(byteVal)
					}
				} else {
					commonVal = values[i]
				}

				dbTuple.Value = append(dbTuple.Value, dbTblAttrValue{AttrName: colName, Val: commonVal})
			}

		}
	}

	return dbTuple
}
