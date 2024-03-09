package sprint

import (
	"net/http"
	"tasklify/internal/pages"
)

type GetCreateSprintHandler struct{}

func NewGetCreateSprintHandler() *GetCreateSprintHandler {
	return &GetCreateSprintHandler{}
}

func (h *GetCreateSprintHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := CreateSprintDialog()
	err := pages.Layout(c, "My website").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
