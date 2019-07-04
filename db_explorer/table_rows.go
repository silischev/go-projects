package main

import (
	"database/sql"
	"fmt"
	"strconv"
)

type dbTblRowAttrValue struct {
	Attr string
	Val  interface{}
}

type dbTblRow struct {
	Row []dbTblRowAttrValue
}

func getRows(db *sql.DB, table string, columns []dbColumn, limit int, offset int) ([]dbTblRow, error) {
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
	var dbTblRs []dbTblRow

	for rows.Next() {
		values := make([]interface{}, len(cols))
		pointers := make([]interface{}, len(cols))

		for i := range values {
			pointers[i] = &values[i]
		}

		rows.Scan(pointers...)

		dbTblRow := dbTblRow{}
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

					dbTblRow.Row = append(dbTblRow.Row, dbTblRowAttrValue{Attr: colName, Val: commonVal})
				}

			}
		}

		dbTblRs = append(dbTblRs, dbTblRow)
	}

	return dbTblRs, nil
}
