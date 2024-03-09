package sprint

import (
	"net/http"
)

type GetCreateSprintHandler struct{}

func NewGetCreateSprintHandler() *GetCreateSprintHandler {
	return &GetCreateSprintHandler{}
}

func (h *GetCreateSprintHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := CreateSprintDialog().Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
