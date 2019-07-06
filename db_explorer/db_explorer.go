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

func NewDbExplorer(db *sql.DB) (*mux.Router, error) {
	handler := &dbHandler{
		db: db,
	}

	mux := mux.NewRouter()
	mux.HandleFunc("/", handler.getTables).Methods("GET")
	mux.HandleFunc("/{table}", handler.getRows).Methods("GET")
	mux.HandleFunc("/{table}/{id:[0-9]+}", handler.getItem).Methods("GET")

	return mux, nil
}

func (h *dbHandler) getTables(w http.ResponseWriter, req *http.Request) {
	data := make(map[string]interface{})
	data["tables"] = h.getTablesFromDb()

	SuccessResponseWrapper(w, req, data)
}

func (h *dbHandler) getRows(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tblName := vars["table"]

	limit := 0
	offset := 0
	var err error

	if req.FormValue("limit") != "" {
		limit, err = strconv.Atoi(req.FormValue("limit"))
		if err != nil {
			ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
			return
		}
	}

	if req.FormValue("offset") != "" {
		offset, err = strconv.Atoi(req.FormValue("offset"))
		if err != nil {
			ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
			return
		}
	}

	if !isTableExist(tblName, h.getTablesFromDb()) {
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

		for _, val := range row.Value {
			rowData[val.AttrName] = val.Val
		}

		dbResponseRs = append(dbResponseRs, rowData)
	}

	data := make(map[string]interface{})
	data["records"] = dbResponseRs

	SuccessResponseWrapper(w, req, data)
}

func (h *dbHandler) getItem(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tblName := vars["table"]
	var id int

	if !isTableExist(tblName, h.getTablesFromDb()) {
		ErrorResponseWrapper(w, req, UnknownTblErr, http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
		return
	}

	if !isTableExist(tblName, h.getTablesFromDb()) {
		ErrorResponseWrapper(w, req, UnknownTblErr, http.StatusNotFound)
		return
	}

	cols, err := getColumns(h.db, tblName)
	if err != nil {
		log.Fatal(err)
	}

	res, err := getItem(h.db, tblName, id, cols)
	if err != nil {
		log.Fatal(err)
	}

	rowData := make(map[string]interface{})

	for _, val := range res.Value {
		rowData[val.AttrName] = val.Val
	}

	data := make(map[string]interface{})
	data["record"] = rowData

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

func isTableExist(tblName string, tblNames []string) bool {
	for _, value := range tblNames {
		if value == tblName {
			return true
		}
	}

	return false
}
