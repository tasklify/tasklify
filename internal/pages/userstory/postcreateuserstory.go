package userstory

import (
	"net/http"
	"strconv"
	"tasklify/internal/database"
)

type PostCreateUserStoryHandler struct{}

func NewPostCreateUserStoryHandler() *PostCreateUserStoryHandler {
	return &PostCreateUserStoryHandler{}
}

func (h *PostCreateUserStoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	title := r.FormValue("title")
	description := r.FormValue("description")
	priority := r.FormValue("priority")
	businessValue, _ := strconv.Atoi(r.FormValue("business_value"))
	projectID, _ := strconv.Atoi(r.FormValue("project_id"))

	realized := false // Default value for new stories

	parsedPriority := database.Priorities.Parse(priority)
	if parsedPriority == nil {
		http.Error(w, "Invalid priority value", http.StatusBadRequest)
		return
	}

	var userStory = &database.UserStory{
		Title:         title,
		Description:   &description,
		Priority:      *parsedPriority,
		BusinessValue: businessValue,
		Realized:      &realized,
		ProjectID:     uint(projectID),
	}

	// Check for duplication and validate priority and business value within the creation process
	err := database.GetDatabase().CreateUserStory(userStory)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/about", http.StatusSeeOther)
}
