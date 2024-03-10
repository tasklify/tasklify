package handlers

import (
	"log"
	"net/http"
)

type RequestParams struct {
	UserID uint
}

type Handler interface {
	Serve(w http.ResponseWriter, r *http.Request) error
}

func UnifiedHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If POST request, parse form
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				http.Error(w, "Error parsing form", http.StatusInternalServerError)
				return
			}
		}

		err := h.Serve(w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	}
}
