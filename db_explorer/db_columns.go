package main

import (
	"database/sql"
)

const PrimaryKey = "PRI"

type dbColumn struct {
	Name      string
	ColumnKey sql.NullString
	Type      string
}

func getColumns(db *sql.DB, table string) ([]dbColumn, error) {
	var cols []dbColumn
	rows, err := db.Query("SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE from information_schema.columns where table_name = ?", table)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		dbColumn := dbColumn{}
		err := rows.Scan(&dbColumn.Name, &dbColumn.ColumnKey, &dbColumn.Type)
		if err != nil {
			return nil, err
		}

		cols = append(cols, dbColumn)
	}

	return cols, nil
}
