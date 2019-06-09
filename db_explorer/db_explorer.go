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

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
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

	h.getTableRowsFromDb(tblName)
	//log.Println(rows)

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

func (h *dbHandler) getTableRowsFromDb(table string) {
	rows, err := h.db.Query(fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		log.Fatal(err)
	}

	//log.Println(rows)

	var tblRows interface{}

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&tblRows)
		log.Println(tblRows)
	}

	//return tblRows
}
