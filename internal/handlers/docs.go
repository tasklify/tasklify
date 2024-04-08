package handlers

import (
	"log"
	"net/http"
)

func DocsHandler(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	}
}
