package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const UnknownTblErr = "unknown table"
const RecordNotFound = "record not found"
const InternalErr = "internal error"

func SuccessResponseWrapper(w http.ResponseWriter, req *http.Request, data map[string]interface{}) {
	response := make(map[string]interface{})
	response["response"] = data

	result, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(string(result)))
}

func ErrorResponseWrapper(w http.ResponseWriter, req *http.Request, data string, status int) {
	response := make(map[string]interface{})
	response["error"] = data

	result, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(status)
	w.Write([]byte(string(result)))
}
