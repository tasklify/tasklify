package userstory

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"tasklify/internal/database"
	"tasklify/internal/handlers"
	"tasklify/internal/web/components/common"
	"tasklify/internal/web/pages"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.RegisterConverter(database.Priority{""}, priorityConverter)
}

func priorityConverter(value string) reflect.Value {
	priorityPtr := database.Priorities.Parse(value)
	if priorityPtr == nil {
		return reflect.Zero(reflect.TypeOf((*database.Priority)(nil)).Elem())
	}
	return reflect.ValueOf(*priorityPtr)
}

func PostUserStory(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type UserStoryFormData struct {
		Title           string            `schema:"title,required"`
		Description     string            `schema:"description,required"`
		Priority        database.Priority `schema:"priority,required"`
		BusinessValue   int               `schema:"business_value,required"`
		AcceptanceTests []string          `schema:"acceptanceTests"`
	}
	var userStoryData UserStoryFormData
	if err := decoder.Decode(&userStoryData, r.PostForm); err != nil {
		return err
	}

	ProjectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}

	userStoryExists := database.GetDatabase().UserStoryInThisProjectAlreadyExists(userStoryData.Title, uint(ProjectID))
	if userStoryExists {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("User story with same title already exists.")
		return c.Render(r.Context(), w)
	}
	if (userStoryData.BusinessValue < 0) || (userStoryData.BusinessValue > 10) {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Business value must be between 0 and 10.")
		return c.Render(r.Context(), w)
	}

	userStory := &database.UserStory{
		Title:           userStoryData.Title,
		Description:     &userStoryData.Description,
		BusinessValue:   userStoryData.BusinessValue,
		Priority:        userStoryData.Priority,
		ProjectID:       uint(ProjectID),
		Realized:        new(bool), // Defaults to false
		AcceptanceTests: []database.AcceptanceTest{},
	}

	if err := database.GetDatabase().CreateUserStory(userStory); err != nil {
		return err
	}

	for _, testDescription := range userStoryData.AcceptanceTests {
		acceptanceTest := &database.AcceptanceTest{
			Description: &testDescription,
			Realized:    new(bool), // Defaults to false
			UserStoryID: userStory.ID,
		}
		if err := database.GetDatabase().CreateAcceptanceTest(acceptanceTest); err != nil {
			return err
		}
	}

	redirectURL := fmt.Sprintf("/productbacklog?projectID=%d", uint(ProjectID))
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

func GetUserStory(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	ProjectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}

	projectRoles, err := database.GetDatabase().GetProjectRoles(params.UserID, uint(ProjectID))
	if err != nil {
		return err
	}

	if (len(projectRoles) == 0 || slices.Contains(projectRoles, database.ProjectRoleDeveloper)) && !slices.Contains(projectRoles, database.ProjectRoleMaster) {
		return pages.NotFound(w, r)
	}

	c := CreateUserStoryDialog(uint(ProjectID))
	return c.Render(r.Context(), w)
}

func DeleteUserStory(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	UserStoryID, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
	if err != nil {
		return err
	}

	userStory, err := database.GetDatabase().GetUserStoryByID(uint(UserStoryID))
	if err != nil {
		return err
	}

	if err := database.GetDatabase().DeleteUserStory(uint(UserStoryID)); err != nil {
		return err
	}

	redirectURL := fmt.Sprintf("/productbacklog?projectID=%d", userStory.ProjectID)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

func generateUUID() string {
	x := [16]byte{}
	_, _ = rand.Read(x[:])
	x[6] = (x[6] & 0x0f) | 0x40
	x[8] = (x[8] & 0x3f) | 0x80

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", x[0:4], x[4:6], x[6:8], x[8:10], x[10:])
	return uuid
}


func AddAcceptanceTest(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	randUUID := generateUUID()
	c := AcceptanceTestDialog(randUUID)
	return c.Render(r.Context(), w)
}

func DeleteAcceptanceTest(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	return nil
}


func GetUserStoryDetails(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type RequestData struct {
		TabID *string `schema:"tabID"`
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	userStoryID, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))

	if err != nil {
		return err
	}

	var requestData RequestData
	err = decoder.Decode(&requestData, r.URL.Query())
	if err != nil {
		return err
	}

	var activeTab string
	if requestData.TabID == nil {
		activeTab = "details"
	} else {
		activeTab = *requestData.TabID
	}

	userStory, err := database.GetDatabase().GetUserStoryByID(uint(userStoryID))
	if err != nil {
		return err
	}

	userStoryComments, err := database.GetDatabase().GetUserStoryComments(uint(userStoryID))
	if err != nil {
		return err
	}
	userStory.UserStoryComments = append(userStory.UserStoryComments, userStoryComments...)

	currentUser, err := database.GetDatabase().GetUserByID(params.UserID)
	if err != nil {
		return err
	}

	projectRoles, err := database.GetDatabase().GetProjectRoles(params.UserID, userStory.ProjectID)
	if err != nil {
		return err
	}

	c := UserStoryDetailsDialog(*userStory, activeTab, *currentUser, projectRoles)
	return c.Render(r.Context(), w)

}


func GetEditUserStory(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	UserStoryID, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
	if err != nil {
		return err
	}

	ProjectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}

	projectRoles, err := database.GetDatabase().GetProjectRoles(params.UserID, uint(ProjectID))
	if err != nil {
		return err
	}

	if (len(projectRoles) == 0 || slices.Contains(projectRoles, database.ProjectRoleDeveloper)) && !slices.Contains(projectRoles, database.ProjectRoleMaster) {
		return pages.NotFound(w, r)
	}

	userStory, err := database.GetDatabase().GetUserStoryByID(uint(UserStoryID))
	if err != nil {
		return err
	}

	c := EditUserStoryDialog(userStory)
	return c.Render(r.Context(), w)
}


func PutUserStory(w http.ResponseWriter, r *http.Request, params handlers.RequestParams) error {
	type UserStoryFormData struct {
		Title           string            `schema:"title,required"`
		Description     string            `schema:"description,required"`
		Priority        database.Priority `schema:"priority,required"`
		BusinessValue   int               `schema:"business_value,required"`
		AcceptanceTests []string          `schema:"acceptanceTests"`
		StoryPoints	 uint              `schema:"story_points"`
	}
	var userStoryData UserStoryFormData
	if err := decoder.Decode(&userStoryData, r.PostForm); err != nil {
		return err
	}

	UserStoryID, err := strconv.Atoi(chi.URLParam(r, "userStoryID"))
	if err != nil {
		return err
	}

	projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		return err
	}

	userStoryExists := database.GetDatabase().UserStoryInThisProjectAlreadyExistsEdit(userStoryData.Title, uint(projectID), uint(UserStoryID))
	if userStoryExists {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("User story with same title already exists.")
		c := common.ValidationError("User story with same title already exists.")
		return c.Render(r.Context(), w)
	}

	if (userStoryData.BusinessValue < 0) || (userStoryData.BusinessValue > 10) {
		w.WriteHeader(http.StatusBadRequest)
		c := common.ValidationError("Business value must be between 0 and 10.")
		return c.Render(r.Context(), w)
	}

	userStory, err := database.GetDatabase().GetUserStoryByID(uint(UserStoryID))
	if err != nil {
		return err
	}

	userStory.Title = userStoryData.Title
	userStory.Description = &userStoryData.Description
	userStory.Priority = userStoryData.Priority
	userStory.BusinessValue = userStoryData.BusinessValue
	userStory.StoryPoints = userStoryData.StoryPoints

	if err := database.GetDatabase().UpdateUserStory(userStory); err != nil {
		return err
	}

	acceptanceTests, err := database.GetDatabase().GetAcceptanceTestsByUserStory(userStory.ID)
	if err != nil {
		return err
	}

	for _, test := range acceptanceTests {
		err := database.GetDatabase().DeleteAcceptanceTest(&test)
		if err != nil {
			return err
		}
	}

	for _, testDescription := range userStoryData.AcceptanceTests {
		acceptanceTest := &database.AcceptanceTest{
			Description: &testDescription,
			Realized:   new(bool), // Defaults to false
			UserStoryID: userStory.ID,
		}
		if err := database.GetDatabase().CreateAcceptanceTest(acceptanceTest); err != nil {
			return err
		}
	}

	redirectURL := fmt.Sprintf("/productbacklog?projectID=%d", userStory.ProjectID)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusSeeOther)

	return nil
}