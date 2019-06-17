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

/* type successResponseData struct {
	Response map[string]interface{} `json:"response"`
} */

type dbResponseResultSet struct {
	Records map[string]interface{} `json:"records"`
}

/* type dbResponseResultSet struct {
	Records []map[string]interface{}
} */

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

func ResponseWriter(w http.ResponseWriter, req *http.Request, data map[string]interface{}) {
	response := make(map[string]interface{})
	response["response"] = data

	result, _ := json.Marshal(response)
	w.Write([]byte(string(result)))
}

func (h *dbHandler) getTables(w http.ResponseWriter, req *http.Request) {
	data := make(map[string]interface{})
	data["tables"] = h.getTablesFromDb()

	ResponseWriter(w, req, data)
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

	var dbResponseRs []map[string]interface{}
	for _, row := range res.Records {
		rowData := make(map[string]interface{})

		for _, val := range row.Row {
			rowData[val.Attr] = val.Val
		}

		dbResponseRs = append(dbResponseRs, rowData)
	}

	data := make(map[string]interface{})
	data["records"] = dbResponseRs

	ResponseWriter(w, req, data)
}

func (h *dbHandler) getTablesFromDb() []string {
	rows, err := h.db.Query("SHOW TABLES;")
	if err != nil {
		log.Fatal(err)
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

			//log.Println(colName)
			//log.Fatal(values[i])

			if ok {
				//log.Println("*1*")
				//log.Println(string(byteVal))
				v = string(byteVal)
			} else {
				//log.Println("*2*")
				//log.Println(values[i])
				v = values[i]
			}

			dbTblRow.Row = append(dbTblRow.Row, dbTblRowAttrValue{Attr: colName, Val: v})
		}

		//log.Fatal("***")

		dbTblRs.Records = append(dbTblRs.Records, dbTblRow)
	}

	return dbTblRs
}
