package sprint

import (
	"net/http"
	"tasklify/internal/pages"
)

type SprintHandler struct{}

func NewSprintHandler() *SprintHandler {
	return &SprintHandler{}
}

func (h *SprintHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := Sprint()
	err := pages.Layout(c, "My website").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
