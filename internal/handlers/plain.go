package handlers

import (
	"net/http"
)

type PlainHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (phf PlainHandlerFunc) Serve(w http.ResponseWriter, r *http.Request) error {
	return phf(w, r)
}
