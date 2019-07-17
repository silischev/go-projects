package main

import (
	"database/sql"
)

const PrimaryKey = "PRI"

type dbColumn struct {
	Name       string
	ColumnKey  sql.NullString
	Type       string
	IsNullable bool
}

func getColumns(db *sql.DB, table string) ([]dbColumn, error) {
	var cols []dbColumn
	rows, err := db.Query("SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE from information_schema.columns where table_name = ?", table)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var isNullable string
		dbColumn := dbColumn{}
		err := rows.Scan(&dbColumn.Name, &dbColumn.ColumnKey, &dbColumn.Type, &isNullable)
		if err != nil {
			return nil, err
		}

		switch dbColumn.Type {
		case "varchar", "text":
			dbColumn.Type = "string"
		}

		switch isNullable {
		case "NO":
			dbColumn.IsNullable = false
		case "YES":
			dbColumn.IsNullable = true
		}

		cols = append(cols, dbColumn)
	}

	return cols, nil
}

func getPrimaryKeyAttr(dbColumns []dbColumn) dbColumn {
	var dbCol dbColumn

	for _, dbColumn := range dbColumns {
		if dbColumn.ColumnKey.String == PrimaryKey {
			dbCol = dbColumn
		}
	}

	return dbCol
}
