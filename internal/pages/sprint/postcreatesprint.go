package sprint

import (
	"math/rand"
	"net/http"
	"strconv"
	"tasklify/internal/database"
	"time"
)

type PostCreateSprintHandler struct{}

func NewPostCreateSprintHandler() *PostCreateSprintHandler {
	return &PostCreateSprintHandler{}
}

func (h *PostCreateSprintHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Parse form data
	layout := "Mon, 01/02/06, 03:04PM"

	startDate, _ := time.Parse(layout, r.FormValue("start_date"))
	endDate, _ := time.Parse(layout, r.FormValue("end_date"))
	convertedVelocity, _ := strconv.ParseFloat(r.FormValue("velocity"), 32)
	velocity := new(float32)
	*velocity = float32(convertedVelocity)

	var sprint = &database.Sprint{
		Title:     strconv.Itoa(rand.Int()), // TODO ask if title is needed for sprint
		StartDate: startDate,
		EndDate:   endDate,
		Velocity:  velocity,
		ProjectID: 1, // Todo, when projects are implemented, change this
	}

	err := database.GetDatabase().CreateSprint(sprint)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/about", http.StatusSeeOther)
}
