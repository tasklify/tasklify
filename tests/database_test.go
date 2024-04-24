package tests

import (
	"fmt"
	"tasklify/internal/config"
	"tasklify/internal/database"
	"testing"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mroth/weightedrand/v2"
	"github.com/stretchr/testify/assert"
)

var (
	db           = database.GetDatabase(config.GetConfig())
	users        []database.User
	projects     []database.Project
	sprints      []database.Sprint
	userStories  []database.UserStory
	tasks        []database.Task
	workSessions []database.WorkSession
)

type userCase struct {
	desc string
	user *database.User
}

type projectCase struct {
	desc    string
	project *database.Project
}

type sprintCase struct {
	desc   string
	sprint *database.Sprint
}

type userStoryCase struct {
	desc      string
	userStory *database.UserStory
}

type taskCase struct {
	desc string
	task *database.Task
}

type workSessionCase struct {
	desc        string
	workSession *database.WorkSession
}

func TestDatabase(t *testing.T) {
	// Users
	var userCases []userCase

	for i := 0; i < gofakeit.IntRange(20, 30); i++ {
		generatedUserCase := userCase{
			desc: "Generated user",
			user: &database.User{
				Username:  gofakeit.Username() + " " + fmt.Sprint(gofakeit.Uint()),
				Password:  gofakeit.Password(gofakeit.Bool(), gofakeit.Bool(), gofakeit.Bool(), gofakeit.Bool(), gofakeit.Bool(), int(gofakeit.UintRange(20, 999))),
				FirstName: gofakeit.Name(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
			},
		}

		generatedUserCase.user.SystemRole = database.SystemRoleUser

		userCases = append(userCases, generatedUserCase)
	}

	/*
		manualUserCase := userCase{
			"Manual user",
			&database.User{
				Username:   "custom user 34",
				Password:   "strongpassword342",
				FirstName:  "FirstName rert",
				LastName:   "LastName t45t5",
				Email:      "someuser_5465@example.com",
				SystemRole: database.SystemRoleAdmin,
			},
		}

		userCases = append(userCases, manualUserCase)
	*/

	for _, c := range userCases {
		t.Run(c.desc, func(t *testing.T) {
			err := db.CreateUser(c.user)
			ok := assert.NoError(t, err)
			if ok && assert.NotZero(t, c.user.ID) {
				users = append(users, *c.user)
			}
		})
	}

	// Projects
	var projectCases []projectCase
	for i := 0; i < gofakeit.IntRange(5, 10); i++ {
		generatedProjectCase := projectCase{
			desc: "Generated project",
			project: &database.Project{
				Title:       gofakeit.BookTitle() + " " + fmt.Sprint(gofakeit.UintRange(1, 9999)),
				Description: gofakeit.Paragraph(1, 5, 12, " "),
				Docs:        gofakeit.Paragraph(3, 4, 7, " "),
			},
		}

		gofakeit.ShuffleAnySlice(users)
		usersToKeep := gofakeit.Number(3, len(users))
		filteredUsers := users[:usersToKeep]

		generatedProjectCase.project.ProductOwnerID = filteredUsers[0].ID
		generatedProjectCase.project.ScrumMasterID = filteredUsers[1].ID
		generatedProjectCase.project.Developers = filteredUsers[2:]

		projectCases = append(projectCases, generatedProjectCase)
	}

	for _, c := range projectCases {
		t.Run(c.desc, func(t *testing.T) {
			ID, err := db.CreateProject(c.project)
			ok := assert.NoError(t, err)
			if ok && assert.NotZero(t, c.project.ID) && assert.Equal(t, ID, c.project.ID) {
				projects = append(projects, *c.project)
			}
		})
	}

	for _, project := range projects {
		t.Run("Test if projects are in database", func(t *testing.T) {

			assert.NotZero(t, project.ID)
			projectOutput, err := db.GetProjectByID(project.ID)
			ok := assert.NoError(t, err)
			if ok {
				assert.NotZero(t, projectOutput)
				assert.Equal(t, projectOutput.ID, project.ID)
			}
		})
	}

	// Sprints
	var sprintCases []sprintCase
	for _, project := range projects {
		t.Run("Generate sprint cases", func(t *testing.T) {
			startDate := gofakeit.DateRange(
				time.Now().Add(-1*time.Duration(gofakeit.IntRange(1, 99)*24*30*int(time.Hour))),
				time.Now().Add(time.Duration(gofakeit.IntRange(1, 10)*24*30*int(time.Hour))),
			)

			// Consecutive sprints
			for ii := 0; ii < gofakeit.IntRange(0, 20); ii++ {
				generatedSprintCase := sprintCase{
					desc: "Generated sprint",
					sprint: &database.Sprint{
						Title:     gofakeit.BookTitle() + " " + fmt.Sprint(ii),
						StartDate: startDate,
						EndDate:   startDate.Add(time.Duration(gofakeit.IntRange(2, 99) * 24 * int(time.Hour))),
						Velocity:  ptr.Float32(gofakeit.Float32Range(10, 200)),
					},
				}

				generatedSprintCase.sprint.ProjectID = project.ID

				sprintCases = append(sprintCases, generatedSprintCase)

				startDate = generatedSprintCase.sprint.EndDate.Add(time.Duration(gofakeit.IntRange(1, 10) * 24 * int(time.Hour)))
			}
		})
	}

	for _, c := range sprintCases {
		t.Run(c.desc, func(t *testing.T) {
			err := db.CreateSprint(c.sprint)
			ok := assert.NoError(t, err)
			if ok {
				sprints = append(sprints, *c.sprint)
			}
		})
	}

	// UserStories
	realizedChooser, err := weightedrand.NewChooser(
		weightedrand.NewChoice(true, 9),
		weightedrand.NewChoice(false, 1),
	)
	assert.NoError(t, err)

	var userStoryCases []userStoryCase
	for i := 0; i < gofakeit.IntRange(150, 400); i++ {
		t.Run("Generate user story cases", func(t *testing.T) {
			generatedUserStoryCase := userStoryCase{
				desc: "Generated userStory",
				userStory: &database.UserStory{
					Title:         gofakeit.BookTitle() + " " + fmt.Sprint(gofakeit.UintRange(1, 9999)),
					Description:   ptr.String(gofakeit.LoremIpsumSentence(22)),
					BusinessValue: gofakeit.UintRange(0, 10),
					StoryPoints:   gofakeit.Float64Range(10, 30),
					Realized:      ptr.Bool(realizedChooser.Pick()),
				},
			}

			priorities := database.Priorities.Members()
			gofakeit.ShuffleAnySlice(priorities)
			generatedUserStoryCase.userStory.Priority = priorities[0]

			gofakeit.ShuffleAnySlice(projects)

			projectID := projects[0].ID

			generatedUserStoryCase.userStory.ProjectID = projectID

			currentProjectSprints, err := db.GetSprintByProject(projectID)
			assert.NoError(t, err)

			if len(currentProjectSprints) != 0 {
				gofakeit.ShuffleAnySlice(currentProjectSprints)
				generatedUserStoryCase.userStory.SprintID = ptr.Uint(currentProjectSprints[0].ID)
			}

			currentProjectUsers, err := db.GetUsersWithRoleOnProject(projectID, database.ProjectRoleDeveloper)
			assert.NoError(t, err)
			assert.NotEmpty(t, currentProjectUsers)

			gofakeit.ShuffleAnySlice(currentProjectUsers)
			generatedUserStoryCase.userStory.UserID = ptr.Uint(currentProjectUsers[0].ID)

			userStoryCases = append(userStoryCases, generatedUserStoryCase)
		})
	}

	for _, c := range userStoryCases {
		t.Run(c.desc, func(t *testing.T) {
			err := db.CreateUserStory(c.userStory)
			ok := assert.NoError(t, err)
			if ok {
				userStories = append(userStories, *c.userStory)
			}
		})
	}

	// Tasks
	var taskCases []taskCase
	for _, userStory := range userStories {
		t.Run("Generate task cases", func(t *testing.T) {
			taskCase := taskCase{
				desc: "Generated task case",
				task: &database.Task{
					Title:       ptr.String(gofakeit.BookTitle() + " " + fmt.Sprint(gofakeit.UintRange(1, 9999))),
					Description: ptr.String(gofakeit.LoremIpsumSentence(22)),
					Status:      &database.StatusTodo,
					ProjectID:   userStory.ProjectID,
					UserID:      userStory.UserID,
					UserStoryID: userStory.ID,
				},
			}

			assert.NotNil(t, userStory.Realized)

			if *userStory.Realized {
				taskCase.task.Status = &database.StatusDone
				taskCases = append(taskCases, taskCase)
			} else {
				if gofakeit.Bool() {
					taskCases = append(taskCases, taskCase)
				}
			}

		})
	}

	for _, c := range taskCases {
		t.Run(c.desc, func(t *testing.T) {
			err := db.CreateTask(c.task)
			ok := assert.NoError(t, err)
			if ok {
				tasks = append(tasks, *c.task)
			}
		})
	}

	// Work sessions
	var workSessionCases []workSessionCase
	for _, task := range tasks {
		t.Run("Generate work cases", func(t *testing.T) {
			userStory, err := db.GetUserStoryByID(task.UserStoryID)
			assert.NoError(t, err)

			sprint, err := db.GetSprintByID(*userStory.SprintID)
			assert.NoError(t, err)

			workSession := workSessionCase{
				desc: "Generated work session",
				workSession: &database.WorkSession{
					StartTime: sprint.StartDate,
					EndTime:   ptr.Time(gofakeit.DateRange(sprint.StartDate, sprint.EndDate)),
					Duration:  time.Duration(userStory.StoryPoints),
					Remaining: time.Duration(gofakeit.Float32Range(0.1, 10) * float32(time.Hour)),
					TaskID:    task.ID,
					UserID:    *task.UserID,
				},
			}

			if userStory.Realized != nil && *userStory.Realized {
				workSession.workSession.Remaining = 0
			}

			workSessionCases = append(workSessionCases, workSession)

		})
	}

	for _, c := range workSessionCases {
		t.Run(c.desc, func(t *testing.T) {
			err = db.CreateWorkSession(c.workSession)
			ok := assert.NoError(t, err)
			if ok {
				workSessions = append(workSessions, *c.workSession)
			}
		})
	}

	// Cleanup
	/*
		t.Cleanup(func() {
			for _, workSession := range workSessions {
				err := db.RawDB().Unscoped().Delete(&workSession).Error
				assert.NoError(t, err)
			}

			for _, task := range tasks {
				err := db.RawDB().Unscoped().Delete(&task).Error
				assert.NoError(t, err)
			}

			for _, userStory := range userStories {
				err := db.RawDB().Unscoped().Delete(&userStory).Error
				assert.NoError(t, err)
			}

			for _, sprint := range sprints {
				err := db.RawDB().Unscoped().Delete(&sprint).Error
				assert.NoError(t, err)
			}

			for _, project := range projects {
				err := db.RawDB().Unscoped().Delete(&project).Error
				assert.NoError(t, err)
			}

			for _, user := range users {
				err := db.RawDB().Unscoped().Delete(&user).Error
				assert.NoError(t, err)
			}
		})
	*/
}
