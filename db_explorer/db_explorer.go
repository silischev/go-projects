package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

type dbHandler struct {
	db *sql.DB
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
	switch req.URL.Path {
	case "/":
		h.getTables()
	}
}

func (h *dbHandler) getTables() {
	rows, err := h.db.Query("SHOW TABLES;")
	if err != nil {
		panic(err)
	}

	var tblName string

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&tblName)
		fmt.Println(tblName)
	}
}
