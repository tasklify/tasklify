package tests

import (
	"tasklify/internal/config"
	"tasklify/internal/database"
	"testing"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

var (
	db          = database.GetDatabase(config.GetConfig())
	users       []database.User
	projects    []database.Project
	sprints     []database.Sprint
	userStories []database.UserStory
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

func TestDatabase(t *testing.T) {
	// Users
	var userCases []userCase

	for i := 0; i < gofakeit.IntRange(3, 9); i++ {
		generatedUserCase := userCase{
			desc: "Generated user",
			user: &database.User{
				Username:  gofakeit.Username(),
				Password:  gofakeit.Password(gofakeit.Bool(), gofakeit.Bool(), gofakeit.Bool(), gofakeit.Bool(), gofakeit.Bool(), int(gofakeit.UintRange(20, 999))),
				FirstName: gofakeit.Name(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
			},
		}

		systemRoles := database.SystemRoles.Members()
		gofakeit.ShuffleAnySlice(systemRoles)
		generatedUserCase.user.SystemRole = systemRoles[0]

		userCases = append(userCases, generatedUserCase)
	}

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

	for _, c := range userCases {
		t.Run(c.desc, func(t *testing.T) {
			err := db.CreateUser(c.user)
			ok := assert.NoError(t, err)
			if ok && assert.NotZero(t, c.user.ID) {
				users = append(users, *c.user)
			}

			t.Cleanup(func() {
				db.DeleteUserByID(c.user.ID)
			})
		})
	}

	// Projects
	var projectCases []projectCase
	for i := 0; i < gofakeit.IntRange(1, 9); i++ {
		generatedProjectCase := projectCase{
			desc: "Generated project",
			project: &database.Project{
				Title:       gofakeit.BookTitle(),
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

			t.Cleanup(func() {
				db.DeleteProject(c.project.ID)
			})
		})
	}

	// Sprints
	var sprintCases []sprintCase
	for i := 0; i < gofakeit.IntRange(1, 9); i++ {

		startDate := gofakeit.DateRange(
			time.Now().Add(-1*time.Duration(gofakeit.IntRange(1, 99)*24*30*int(time.Hour))),
			time.Now().Add(time.Duration(gofakeit.IntRange(1, 10)*24*30*int(time.Hour))),
		)

		gofakeit.ShuffleAnySlice(projects)
		projectID := projects[0].ID
		assert.NotZero(t, projectID)

		// Consecutive sprints
		for i := 0; i < gofakeit.IntRange(1, 1); i++ {
			generatedSprintCase := sprintCase{
				desc: "Generated sprint",
				sprint: &database.Sprint{
					Title:     gofakeit.BookTitle(),
					StartDate: startDate,
					EndDate:   startDate.Add(time.Duration(gofakeit.IntRange(2, 99) * 24 * int(time.Hour))),
					Velocity:  ptr.Float32(gofakeit.Float32Range(10, 9999)),
				},
			}

			generatedSprintCase.sprint.ProjectID = projectID

			sprintCases = append(sprintCases, generatedSprintCase)

			startDate = generatedSprintCase.sprint.EndDate.Add(time.Duration(gofakeit.IntRange(1, 10) * 24 * int(time.Hour)))
		}
	}

	for _, c := range sprintCases {
		t.Run(c.desc, func(t *testing.T) {
			err := db.CreateSprint(c.sprint)
			ok := assert.NoError(t, err)
			if ok {
				sprints = append(sprints, *c.sprint)
			}

			t.Cleanup(func() {
				db.DeleteSprint(c.sprint.ProjectID, c.sprint.ID)
			})
		})
	}

	/*
		// UserStories
		var userStoryCases []userStoryCase
		for i := 1; i < gofakeit.IntRange(999, 99999); i++ {
			generatedUserStoryCase := userStoryCase{
				desc: "Generated userStory",
				userStory: &database.UserStory{
					Title:         gofakeit.BookTitle(),
					Description:   ptr.String(gofakeit.LoremIpsumSentence(22)),
					BusinessValue: gofakeit.Uint(),
					StoryPoints:   gofakeit.Float64Range(0.1, 99),
				},
			}

			priorities := database.Priorities.Members()
			gofakeit.ShuffleAnySlice(priorities)
			generatedUserStoryCase.userStory.Priority = priorities[0]

			options := []any{true, false}
			weights := []float32{0.9, 0.1}
			realized, err := gofakeit.Weighted(options, weights)
			assert.NoError(t, err)

			generatedUserStoryCase.userStory.Realized = ptr.Bool(realized.(bool))

			gofakeit.ShuffleAnySlice(projects)
			generatedUserStoryCase.userStory.ProjectID = projects[0].ID

			currentProjectSprints, err := db.GetSprintByProject(projects[0].ID)
			assert.NoError(t, err)
			assert.NotEmpty(t, currentProjectSprints)

			gofakeit.ShuffleAnySlice(currentProjectSprints)
			generatedUserStoryCase.userStory.SprintID = &currentProjectSprints[0].ID

			currentProjectUsers, err := db.GetUsersWithRoleOnProject(projects[0].ID, database.ProjectRoleDeveloper)
			assert.NoError(t, err)
			assert.NotEmpty(t, currentProjectUsers)

			gofakeit.ShuffleAnySlice(currentProjectUsers)
			generatedUserStoryCase.userStory.UserID = &currentProjectUsers[0].ID

			userStoryCases = append(userStoryCases, generatedUserStoryCase)
		}

		for _, c := range userStoryCases {
			t.Run(c.desc, func(t *testing.T) {
				err := db.CreateUserStory(c.userStory)
				ok := assert.NoError(t, err)
				if ok {
					userStories = append(userStories, *c.userStory)
				}

				t.Cleanup(func() {
					db.DeleteUserStory(c.userStory.ID)
				})
			})
		}
	*/
}
