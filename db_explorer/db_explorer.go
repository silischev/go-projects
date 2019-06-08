package main

import (
	"database/sql"
	"encoding/json"
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
func NewDbExplorer(db *sql.DB) (*dbHandler, error) {
	handler := &dbHandler{
		db: db,
	}

	return handler, nil
}

func (h *dbHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	res := ""

	switch req.URL.Path {
	case "/":
		res = h.getTables()
	}

	w.Write([]byte(res))
}

func (h *dbHandler) getTables() string {
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
	/* fmt.Println(string(result))
	log.Println(r)
	log.Fatal("***") */

	return string(result)
}
