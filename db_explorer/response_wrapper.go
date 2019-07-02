package main

import (
	"encoding/json"
	"net/http"
)

const UnknownTblErr = "unknown table"
const InternalErr = "internal error"

func SuccessResponseWrapper(w http.ResponseWriter, req *http.Request, data map[string]interface{}) {
	response := make(map[string]interface{})
	response["response"] = data

	result, _ := json.Marshal(response)
	w.Write([]byte(string(result)))
}

func ErrorResponseWrapper(w http.ResponseWriter, req *http.Request, data string, status int) {
	response := make(map[string]interface{})
	response["error"] = data

	result, _ := json.Marshal(response)
	w.WriteHeader(status)
	w.Write([]byte(string(result)))
}
