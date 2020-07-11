// Code generated automatically. DO NOT EDIT.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"unicode/utf8"
)

func (h *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/user/profile":

		h.handlerProfile(w, r)

	case "/user/create":

		if r.Method != http.MethodPost {
			http.Error(w, "{\"error\": \"bad method\"}", http.StatusNotAcceptable)

			return
		}

		h.handlerCreate(w, r)

	default:
		http.Error(w, "{\"error\": \"unknown method\"}", http.StatusNotFound)
	}
}

func (h *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/user/create":

		if r.Method != http.MethodPost {
			http.Error(w, "{\"error\": \"bad method\"}", http.StatusNotAcceptable)

			return
		}

		h.handlerCreate(w, r)

	default:
		http.Error(w, "{\"error\": \"unknown method\"}", http.StatusNotFound)
	}
}

func (s *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatalln("Error parse query: ", err)
	}
	params := ProfileParams{}
	var rawVal string
	rawVal = ""

	if len(r.Form["login"]) != 0 {
		rawVal = r.Form["login"][0]
	}
	login := rawVal

	if rawVal == "" {
		http.Error(w, "{\"error\": \"login must me not empty\"}", http.StatusBadRequest)

		return
	}
	params.Login = login

	res, err := s.Profile(r.Context(), params)
	if err != nil {
		switch err.(type) {
		case ApiError:
			apiErr := err.(ApiError)
			http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", apiErr.Error()), apiErr.HTTPStatus)

			return
		default:
			http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), http.StatusInternalServerError)

			return
		}
	}

	resp, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"error": "", "response": %s}`, resp)))
}
func (s *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatalln("Error parse query: ", err)
	}
	params := CreateParams{}
	var rawVal string
	rawVal = ""

	if len(r.Form["login"]) != 0 {
		rawVal = r.Form["login"][0]
	}
	login := rawVal

	if rawVal == "" {
		http.Error(w, "{\"error\": \"login must me not empty\"}", http.StatusBadRequest)

		return
	}

	if utf8.RuneCountInString(login) < 10 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	params.Login = login

	rawVal = ""

	if len(r.Form["name"]) != 0 {
		rawVal = r.Form["name"][0]
	}
	name := rawVal
	name = r.Form["full_name"][0]
	params.Name = name

	rawVal = ""

	if len(r.Form["status"]) != 0 {
		rawVal = r.Form["status"][0]
	}
	status := rawVal
	params.Status = status

	rawVal = ""

	if len(r.Form["age"]) != 0 {
		rawVal = r.Form["age"][0]
	}

	age, err := strconv.Atoi(rawVal)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if age < 0 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if age > 128 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	params.Age = age

	res, err := s.Create(r.Context(), params)
	if err != nil {
		switch err.(type) {
		case ApiError:
			apiErr := err.(ApiError)
			http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", apiErr.Error()), apiErr.HTTPStatus)

			return
		default:
			http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), http.StatusInternalServerError)

			return
		}
	}

	resp, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"error": "", "response": %s}`, resp)))
}

func (s *OtherApi) handlerCreate(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatalln("Error parse query: ", err)
	}
	params := OtherCreateParams{}
	var rawVal string
	rawVal = ""

	if len(r.Form["username"]) != 0 {
		rawVal = r.Form["username"][0]
	}
	username := rawVal

	if rawVal == "" {
		http.Error(w, "{\"error\": \"username must me not empty\"}", http.StatusBadRequest)

		return
	}

	if utf8.RuneCountInString(username) < 3 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	params.Username = username

	rawVal = ""

	if len(r.Form["name"]) != 0 {
		rawVal = r.Form["name"][0]
	}
	name := rawVal
	name = r.Form["account_name"][0]
	params.Name = name

	rawVal = ""

	if len(r.Form["class"]) != 0 {
		rawVal = r.Form["class"][0]
	}
	class := rawVal
	params.Class = class

	rawVal = ""

	if len(r.Form["level"]) != 0 {
		rawVal = r.Form["level"][0]
	}

	level, err := strconv.Atoi(rawVal)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if level < 1 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if level > 50 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	params.Level = level

	res, err := s.Create(r.Context(), params)
	if err != nil {
		switch err.(type) {
		case ApiError:
			apiErr := err.(ApiError)
			http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", apiErr.Error()), apiErr.HTTPStatus)

			return
		default:
			http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), http.StatusInternalServerError)

			return
		}
	}

	resp, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"error": "", "response": %s}`, resp)))
}
