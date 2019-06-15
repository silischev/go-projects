package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type dbHandler struct {
	db *sql.DB
}

type successResponseData struct {
	Response map[string]interface{} `json:"response"`
}

type errorResponseData struct {
	Error string `json:"error"`
}

type dbTblRowAttrValue struct {
	Attr string
	Val  interface{}
}

type dbTblRow struct {
	Row []dbTblRowAttrValue
}

type dbTblResultSet struct {
	Records []dbTblRow
}

func NewDbExplorer(db *sql.DB) (*mux.Router, error) {
	handler := &dbHandler{
		db: db,
	}

	mux := mux.NewRouter()
	mux.HandleFunc("/", handler.getTables)
	mux.HandleFunc("/{table}", handler.getTableRows)

	return mux, nil
}

func (h *dbHandler) getTables(w http.ResponseWriter, req *http.Request) {
	r := successResponseData{}
	r.Response = make(map[string]interface{})

	r.Response["tables"] = h.getTablesFromDb()
	result, _ := json.Marshal(r)

	w.Write([]byte(string(result)))
}

func (h *dbHandler) getTableRows(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tblName := vars["table"]
	tblNames := h.getTablesFromDb()
	found := false

	for _, value := range tblNames {
		if value == tblName {
			found = true
			break
		}
	}

	if !found {
		r := errorResponseData{Error: "unknown table"}
		result, _ := json.Marshal(r)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(string(result)))

		return
	}

	res := h.getTableRowsFromDb(tblName)

}

func (h *dbHandler) getTablesFromDb() []string {
	rows, err := h.db.Query("SHOW TABLES;")
	if err != nil {
		panic(err)
	}

	var tblName string
	var tblNames []string

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&tblName)
		tblNames = append(tblNames, tblName)
	}

	return tblNames
}

func (h *dbHandler) getTableRowsFromDb(table string) dbTblResultSet {
	rows, err := h.db.Query(fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		log.Fatal(err)
	}

	cols, _ := rows.Columns()
	dbTblRs := dbTblResultSet{}

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
			var v interface{}
			byteVal, ok := values[i].([]byte)

			if ok {
				v = string(byteVal)
			} else {
				v = values[i]
			}

			dbTblRow.Row = append(dbTblRow.Row, dbTblRowAttrValue{Attr: colName, Val: v})
		}

		dbTblRs.Records = append(dbTblRs.Records, dbTblRow)
	}

	return dbTblRs
}
