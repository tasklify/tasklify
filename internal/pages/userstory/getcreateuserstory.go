package userstory

import (
	"net/http"
)

type GetCreateUserStoryHandler struct{}

func NewGetCreateUserStoryHandler() *GetCreateUserStoryHandler {
	return &GetCreateUserStoryHandler{}
}

func (h *GetCreateUserStoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := CreateUserStoryDialog().Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
