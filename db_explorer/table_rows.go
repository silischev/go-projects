package main

import (
	"database/sql"
	"fmt"
	"log"
)

type dbTblRowAttrValue struct {
	Attr string
	Val  interface{}
}

type dbTblRow struct {
	Row []dbTblRowAttrValue
}

func getRows(db *sql.DB, table string, columns []dbColumn) ([]dbTblRow, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		return nil, err
	}

	cols, _ := rows.Columns()
	var dbTblRs []dbTblRow

	defer rows.Close()
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
					switch column.Type {
					//case int
					}

					var v interface{}
					byteVal, ok := values[i].([]byte)

					if ok {
						v = string(byteVal)
					} else {
						v = values[i]
					}

					log.Println(v)
				}
			}

			//dbTblRow.Row = append(dbTblRow.Row, dbTblRowAttrValue{Attr: colName, Val: v})
		}

		dbTblRs = append(dbTblRs, dbTblRow)
	}

	return dbTblRs, nil
}
