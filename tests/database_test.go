package tests

import (
	"tasklify/internal/config"
	"tasklify/internal/database"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	db    = database.GetDatabase(config.GetConfig())
	users []database.User
)

func TestDatabase(t *testing.T) {
	// Users
	userCases := []struct {
		desc string
		user database.User
	}{
		{"Generated user", database.User{Username: gofakeit.Username(), Password: gofakeit.Password(gofakeit.Bool(), gofakeit.Bool(), gofakeit.Bool(), gofakeit.Bool(), gofakeit.Bool(), int(gofakeit.UintRange(20, 999))), FirstName: gofakeit.Name(), LastName: gofakeit.LastName(), Email: gofakeit.Email()}},
		// {"Manual user", database.User{Username: "custom user", Password: "strongpassword342", FirstName: "FirstName", LastName: "LastName", Email: "someuser@example.com", SystemRole: database.SystemRoleAdmin}},
	}

	for _, c := range userCases {
		t.Run(c.desc, func(t *testing.T) {
			if c.user.SystemRole.Val == "" {
				users := database.SystemRoles.Members()
				gofakeit.ShuffleAnySlice(users)
				c.user.SystemRole = users[0]
			}

			t.Cleanup(func() {
				db.DeleteUserByID(c.user.ID)
			})

			err := db.CreateUser(&c.user)
			assert.NoError(t, err)
		})
	}

	t.Run("GetUsers", func(t *testing.T) {
		var err error
		users, err = db.GetUsers()
		require.NoError(t, err)
		assert.NotEmpty(t, users)
	})

	// Projects
	projectCases := []struct {
		desc    string
		project database.Project
	}{
		{"Generated project", database.Project{Title: gofakeit.BookTitle(), Description: gofakeit.Paragraph(1, 5, 12, " "), Docs: gofakeit.Paragraph(3, 4, 7, " "), Developers: nil}},
	}

	for _, c := range projectCases {
		t.Run(c.desc, func(t *testing.T) {
			if c.project.Developers == nil {
				usersToKeep := gofakeit.Number(3, len(users))

				// Select the first 'itemsToKeep' elements from the shuffled slice
				filteredUsers := users[:usersToKeep]

				c.project.ProductOwnerID = filteredUsers[0].ID
				c.project.ScrumMasterID = filteredUsers[1].ID
				c.project.Developers = filteredUsers[2:]
			}

			t.Cleanup(func() {
				db.Delete(c.user.ID)
			})

			id, err := db.CreateProject(&c.project)
			ok := assert.NoError(t, err)
			if ok {
				assert.Greater(t, id, uint(0))
			}
		})
	}
}
