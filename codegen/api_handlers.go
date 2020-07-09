// Code generated automatically. DO NOT EDIT.
package main

import (
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
		h.handlerCreate(w, r)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/user/create":
		h.handlerCreate(w, r)

	default:
		w.WriteHeader(http.StatusNotFound)
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
	}
	params.Login = login

	res, err := s.Profile(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

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
	}

	if utf8.RuneCountInString(login) < 10 {
	}
	params.Login = login

	rawVal = ""

	if len(r.Form["name"]) != 0 {
		rawVal = r.Form["name"][0]
	}

	name := rawVal
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

	}

	if age < 0 {
	}

	if age > 128 {
	}
	params.Age = age

	res, err := s.Create(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

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

	if utf8.RuneCountInString(username) < 3 {
	}

	if rawVal == "" {
	}
	params.Username = username

	rawVal = ""

	if len(r.Form["name"]) != 0 {
		rawVal = r.Form["name"][0]
	}

	name := rawVal
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

	}

	if level < 1 {
	}

	if level > 50 {
	}
	params.Level = level

	res, err := s.Create(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}
