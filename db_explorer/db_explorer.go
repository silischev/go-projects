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

	res := h.getTableRowsFromDb(tblName)
	log.Println(res)

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

func (h *dbHandler) getTableRowsFromDb(table string) []map[string]interface{} {
	rows, err := h.db.Query(fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		log.Fatal(err)
	}

	cols, _ := rows.Columns()
	res := make(map[string]interface{})
	res2 := make([]map[string]interface{}, 2)

	defer rows.Close()
	for rows.Next() {
		values := make([]interface{}, len(cols))
		pointers := make([]interface{}, len(cols))
		for i := range values {
			pointers[i] = &values[i]
		}

		rows.Scan(pointers...)

		for i, colName := range cols {
			var v interface{}
			byteVal, ok := values[i].([]byte)

			if ok {
				v = string(byteVal)
			} else {
				v = values[i]
			}

			res[colName] = v
		}

		//log.Fatal(res)

		res2 = append(res2, res)
		//log.Fatal(res2)
	}

	log.Fatal(res2)

	return res2
}
