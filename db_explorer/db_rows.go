package main

import (
	"database/sql"
	"fmt"
	"strconv"
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
	cols, _ := rows.Columns()
	var dbTblRs []dbTuple

	for rows.Next() {
		dbTblRs = append(dbTblRs, getTuple(cols, rows, columns))
	}

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
