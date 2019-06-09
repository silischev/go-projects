package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type dbHandler struct {
	db *sql.DB
}

type responseData struct {
	Response map[string]interface{} `json:"response"`
}

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
func NewDbExplorer(db *sql.DB) (*http.ServeMux, error) {
	handler := &dbHandler{
		db: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.getTables)
	//mux.HandleFunc("/{table}", handler.getTables)

	return mux, nil
}

func (h *dbHandler) getTables(w http.ResponseWriter, req *http.Request) {
	r := responseData{}
	r.Response = make(map[string]interface{})

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

	r.Response["tables"] = tblNames
	result, _ := json.Marshal(r)

	w.Write([]byte(string(result)))
}

func (h *dbHandler) getTableRows(w http.ResponseWriter, req *http.Request) {
	//keys, ok := r.URL.Query()["key"]
	limit := req.Form.Get("limit")
	offset := req.Form.Get("offset")

	log.Println(limit)
	log.Println(offset)
}
