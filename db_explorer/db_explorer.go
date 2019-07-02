package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type dbHandler struct {
	db *sql.DB
}

type dbResponseResultSet struct {
	Records map[string]interface{} `json:"records"`
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
	data := make(map[string]interface{})
	data["tables"] = h.getTablesFromDb()

	SuccessResponseWrapper(w, req, data)
}

func (h *dbHandler) getTableRows(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tblName := vars["table"]

	var limit int
	var offset int
	var err error

	if req.FormValue("limit") != "" {
		limit, err = strconv.Atoi(req.FormValue("limit"))
		if err != nil {
			ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
			return
		}
	} else {
		limit = 0
	}

	if req.FormValue("offset") != "" {
		offset, err = strconv.Atoi(req.FormValue("offset"))
		if err != nil {
			log.Fatal(err)
			ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
			return
		}
	} else {
		offset = 0
	}

	tblNames := h.getTablesFromDb()
	found := false

	for _, value := range tblNames {
		if value == tblName {
			found = true
			break
		}
	}

	if !found {
		ErrorResponseWrapper(w, req, UnknownTblErr, http.StatusNotFound)
		return
	}

	cols, err := getColumns(h.db, tblName)
	if err != nil {
		log.Fatal(err)
	}

	res, err := getRows(h.db, tblName, cols, limit, offset)
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
