// Code generated. DO NOT EDIT.
package main

import "net/http"

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
	res, err := s.Profile(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	res, err := s.Create(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *OtherApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	res, err := s.Create(ctx, params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
