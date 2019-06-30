package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type dbHandler struct {
	db *sql.DB
}

type dbResponseResultSet struct {
	Records map[string]interface{} `json:"records"`
}

type errorResponseData struct {
	Error  string `json:"error"`
	status int
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

func SuccessResponseWrapper(w http.ResponseWriter, req *http.Request, data map[string]interface{}) {
	response := make(map[string]interface{})
	response["response"] = data

	result, _ := json.Marshal(response)
	w.Write([]byte(string(result)))
}

func ErrorResponseWrapper(w http.ResponseWriter, req *http.Request, responseData errorResponseData) {
	result, _ := json.Marshal(responseData)
	w.WriteHeader(responseData.status)
	w.Write([]byte(string(result)))
}

func (h *dbHandler) getTables(w http.ResponseWriter, req *http.Request) {
	data := make(map[string]interface{})
	data["tables"] = h.getTablesFromDb()

	SuccessResponseWrapper(w, req, data)
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
		r := errorResponseData{
			Error:  "unknown table",
			status: http.StatusNotFound,
		}

		ErrorResponseWrapper(w, req, r)

		return
	}

	cols, err := getColumns(h.db, tblName)
	if err != nil {
		log.Fatal(err)
	}

	res, err := getRows(h.db, tblName, cols)
	if err != nil {
		log.Fatal(err)
	}

	var dbResponseRs []map[string]interface{}
	for _, row := range res {
		rowData := make(map[string]interface{})

		for _, val := range row.Row {
			rowData[val.Attr] = val.Val
		}

		dbResponseRs = append(dbResponseRs, rowData)
	}

	data := make(map[string]interface{})
	data["records"] = dbResponseRs

	SuccessResponseWrapper(w, req, data)
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
