package sprint

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"tasklify/internal/database"
	"testing"
	"time"
)

// todo move to common place and maintain
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetUsers() ([]database.User, error) {
	args := m.Called()
	return args.Get(0).([]database.User), args.Error(1)
}

func (m *MockDatabase) GetUserByUsername(username string) (*database.User, error) {
	args := m.Called(username)
	return args.Get(0).(*database.User), args.Error(1)
}

func (m *MockDatabase) GetUserByID(id uint) (*database.User, error) {
	args := m.Called(id)
	return args.Get(0).(*database.User), args.Error(1)
}

func (m *MockDatabase) UpdateUser(user *database.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDatabase) CreateUser(user *database.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDatabase) GetSprintByID(id uint) *database.Sprint {
	args := m.Called(id)
	return args.Get(0).(*database.Sprint)
}

func (m *MockDatabase) CreateUserStory(userStory *database.UserStory) error {
	args := m.Called(userStory)
	return args.Error(0)
}

func (m *MockDatabase) UpdateUserStory(userStory *database.UserStory) error {
	args := m.Called(userStory)
	return args.Error(0)
}

func (m *MockDatabase) AddUserStoryToSprint(sprintID uint, userStories []uint) (*database.Sprint, error) {
	args := m.Called(sprintID, userStories)
	return args.Get(0).(*database.Sprint), args.Error(1)
}

func (m *MockDatabase) GetUserStoriesByProject(projectID uint) ([]database.UserStory, error) {
	args := m.Called(projectID)
	return args.Get(0).([]database.UserStory), args.Error(1)
}

func (m *MockDatabase) GetUserStoriesBySprint(sprintID uint) ([]database.UserStory, error) {
	args := m.Called(sprintID)
	return args.Get(0).([]database.UserStory), args.Error(1)
}

func (m *MockDatabase) GetUserStoryByID(id uint) (*database.UserStory, error) {
	args := m.Called(id)
	return args.Get(0).(*database.UserStory), args.Error(1)
}

func (m *MockDatabase) UserStoryWithTitleExists(title string) bool {
	args := m.Called(title)
	return args.Bool(0)
}

func (m *MockDatabase) GetTasksByUserStory(userStoryID uint) ([]database.Task, error) {
	args := m.Called(userStoryID)
	return args.Get(0).([]database.Task), args.Error(1)
}

func (m *MockDatabase) CreateTask(task *database.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockDatabase) GetProjectByID(id uint) (*database.Project, error) {
	args := m.Called(id)
	return args.Get(0).(*database.Project), args.Error(1)
}

func (m *MockDatabase) GetProjectRole(userID uint, projectID uint) database.ProjectRole {
	args := m.Called(userID, projectID)
	return args.Get(0).(database.ProjectRole)
}

func (m *MockDatabase) GetProjectHasUserByProjectAndUser(userID uint, projectID uint) (*database.ProjectHasUser, error) {
	args := m.Called(userID, projectID)
	return args.Get(0).(*database.ProjectHasUser), args.Error(1)
}

func (m *MockDatabase) CreateProject(project *database.Project) (uint, error) {
	args := m.Called(project)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockDatabase) AddUserToProject(projectID uint, userID uint, projectRole string) error {
	args := m.Called(projectID, userID, projectRole)
	return args.Error(0)
}

func (m *MockDatabase) GetUsersOnProject(projectID uint) ([]database.User, error) {
	args := m.Called(projectID)
	return args.Get(0).([]database.User), args.Error(1)
}

func (m *MockDatabase) GetUsersNotOnProject(projectID uint) ([]database.User, error) {
	args := m.Called(projectID)
	return args.Get(0).([]database.User), args.Error(1)
}

func (m *MockDatabase) ProjectWithTitleExists(title string) bool {
	args := m.Called(title)
	return args.Bool(0)
}

func (m *MockDatabase) RemoveUserFromProject(projectID uint, userID uint) error {
	args := m.Called(projectID, userID)
	return args.Error(0)
}

func (m *MockDatabase) GetUserProjects(userID uint) ([]database.Project, error) {
	args := m.Called(userID)
	return args.Get(0).([]database.Project), args.Error(1)
}

func (m *MockDatabase) CreateAcceptanceTest(acceptanceTest *database.AcceptanceTest) error {
	args := m.Called(acceptanceTest)
	return args.Error(0)
}

func (m *MockDatabase) RawDB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDatabase) GetSprintByProject(projectID uint) ([]database.Sprint, error) {
	args := m.Called(projectID)
	return args.Get(0).([]database.Sprint), args.Error(1)
}

func (m *MockDatabase) CreateSprint(sprint *database.Sprint) error {
	args := m.Called(sprint)
	return args.Error(0)
}
func TestPostSprint(t *testing.T) {
	// arrange
	mockDB := new(MockDatabase)
	//mockDB.On("GetSprintByProject", mock.Anything).Return([]database.Sprint{}, nil)

	now := time.Date(2024, time.March, 18, 0, 0, 0, 0, time.UTC)
	formData := url.Values{}
	formData.Set("start_date", now.Add(24*time.Hour).Format(time.DateOnly))
	formData.Set("end_date", now.Add(48*time.Hour).Format(time.DateOnly))
	formData.Set("velocity", "3")
	formData.Set("project_id", "3")
	req := httptest.NewRequest("POST", "/createsprint", nil)
	req.PostForm = formData
	w := httptest.NewRecorder()

	// act
	err := postSprint(w, req, mockDB)

	// assert
	if err != nil {
		t.Errorf("PostSprint returned an error: %v", err)
	}
	resp := w.Result()

	// Check response status code
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Expected status code %d, got %d", http.StatusSeeOther, resp.StatusCode)
	}
}

func Test_fieldValidation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		sprintToAdd sprintFormData
		sprints     []database.Sprint
		want        bool
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Start before End",
			sprintToAdd: sprintFormData{
				StartDate: now,
				EndDate:   now.Add(24 * time.Hour),
			},
			sprints:     []database.Sprint{},
			want:        true,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "Start and End same date",
			sprintToAdd: sprintFormData{
				StartDate: now,
				EndDate:   now,
			},
			sprints:     []database.Sprint{},
			want:        false,
			wantErr:     true,
			expectedErr: "start date should be before end date",
		},
		{
			name: "Start date after end date",
			sprintToAdd: sprintFormData{
				StartDate: now.Add(24 * time.Hour),
				EndDate:   now,
			},
			sprints:     []database.Sprint{},
			want:        false,
			wantErr:     true,
			expectedErr: "start date should be before end date",
		},

		{
			name: "Start date in past",
			sprintToAdd: sprintFormData{
				StartDate: now.Add(-24 * time.Hour),
				EndDate:   now,
			},
			sprints:     []database.Sprint{},
			want:        false,
			wantErr:     true,
			expectedErr: "start date should not be in the past",
		},
		{
			name: "Overlapping sprint before",
			sprintToAdd: sprintFormData{
				StartDate: now,
				EndDate:   now.Add(24 * time.Hour),
			},
			sprints: []database.Sprint{
				{StartDate: now.Add(48 * time.Hour), EndDate: now.Add(72 * time.Hour)},
			},
			want:        true,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "Overlapping sprint touching start",
			sprintToAdd: sprintFormData{
				StartDate: now,
				EndDate:   now.Add(24 * time.Hour),
			},
			sprints: []database.Sprint{
				{StartDate: now.Add(24 * time.Hour), EndDate: now.Add(48 * time.Hour)},
			},
			want:        false,
			wantErr:     true,
			expectedErr: "sprint should not overlap with an existing one",
		},
		// todo conditions
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fieldValidation(tt.sprintToAdd, tt.sprints)
			if (err != nil) != tt.wantErr {
				t.Errorf("fieldValidation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.expectedErr {
				t.Errorf("fieldValidation() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			if got != tt.want {
				t.Errorf("fieldValidation() got = %v, want %v", got, tt.want)
			}
		})
	}
}
