package main

import (
	"database/sql"
	"encoding/json"
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
	mux.HandleFunc("/{table}/", handler.createItem).Methods("PUT")
	mux.HandleFunc("/{table}/{id:[0-9]+}", handler.updateItem).Methods("POST")
	mux.HandleFunc("/{table}/{id:[0-9]+}", handler.deleteItem).Methods("DELETE")

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
		ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
		return
	}

	res, err := getItem(h.db, tblName, id, cols)
	if err != nil {
		ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
		return
	}

	if len(res.Value) == 0 {
		ErrorResponseWrapper(w, req, RecordNotFound, http.StatusNotFound)
		return
	}

	rowData := make(map[string]interface{})

	for _, val := range res.Value {
		rowData[val.AttrName] = val.Val
	}

	data := make(map[string]interface{})
	data["record"] = rowData

	SuccessResponseWrapper(w, req, data)
}

func (h *dbHandler) createItem(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tblName := vars["table"]

	if !isTableExist(tblName, h.getTablesFromDb()) {
		ErrorResponseWrapper(w, req, UnknownTblErr, http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var reqBodyParams map[string]interface{}
	err := decoder.Decode(&reqBodyParams)
	if err != nil {
		ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
		return
	}

	cols, err := getColumns(h.db, tblName)
	if err != nil {
		ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	result, err := createItem(h.db, tblName, cols, reqBodyParams)
	if err != nil {
		if httpErr, ok := err.(httpError); ok {
			ErrorResponseWrapper(w, req, httpErr.OriginalError.Error(), httpErr.Status)
			return
		}

		ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
		return
	}

	for _, row := range result {
		data[row.Value[0].AttrName] = row.Value[0].Val
	}

	SuccessResponseWrapper(w, req, data)
}

func (h *dbHandler) updateItem(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tblName := vars["table"]

	if !isTableExist(tblName, h.getTablesFromDb()) {
		ErrorResponseWrapper(w, req, UnknownTblErr, http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var reqBodyParams map[string]interface{}
	err = decoder.Decode(&reqBodyParams)
	if err != nil {
		ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
		return
	}

	cols, err := getColumns(h.db, tblName)
	if err != nil {
		ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	result, err := updateItem(h.db, tblName, id, cols, reqBodyParams)
	if err != nil {
		if httpErr, ok := err.(httpError); ok {
			ErrorResponseWrapper(w, req, httpErr.OriginalError.Error(), httpErr.Status)
			return
		}

		ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
		return
	}

	for _, row := range result {
		data[row.Value[0].AttrName] = row.Value[0].Val
	}

	SuccessResponseWrapper(w, req, data)
}

func (h *dbHandler) deleteItem(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tblName := vars["table"]

	if !isTableExist(tblName, h.getTablesFromDb()) {
		ErrorResponseWrapper(w, req, UnknownTblErr, http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	result, err := deleteItem(h.db, tblName, id)
	if err != nil {
		if httpErr, ok := err.(httpError); ok {
			ErrorResponseWrapper(w, req, httpErr.OriginalError.Error(), httpErr.Status)
			return
		}

		ErrorResponseWrapper(w, req, InternalErr, http.StatusInternalServerError)
		return
	}

	data["deleted"] = result

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
