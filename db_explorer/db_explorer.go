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
	Error string `json:"error"`
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

	cols, err := getColumns(h.db, tblName)
	if err != nil {
		log.Fatal(err)
	}

	res, err := getRows(h.db, tblName, cols)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(res)

	data := make(map[string]interface{})

	/*
		var dbResponseRs []map[string]interface{}
		for _, row := range res.Records {
			rowData := make(map[string]interface{})

			for _, val := range row.Row {
				rowData[val.Attr] = val.Val
			}

			dbResponseRs = append(dbResponseRs, rowData)
		} */

	//data := make(map[string]interface{})
	//data["records"] = dbResponseRs

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
