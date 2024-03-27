package handlers

import (
	"log"
	"net/http"
	"slices"
)

type RequestParams struct {
	UserID uint
}

type Handler interface {
	Serve(w http.ResponseWriter, r *http.Request) error
}

func UnifiedHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If appropriate request, parse form
		if slices.Contains([]string{http.MethodPost, http.MethodPut, http.MethodPatch}, r.Method) {
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
