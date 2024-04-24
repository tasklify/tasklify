package tests

import (
	"net/http"
	"net/url"
	"tasklify/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticatedAPIRequests(t *testing.T) {
	cfg := config.GetConfig()
	baseURL := "http://localhost:" + cfg.Port

	// Step 1: Login to get session token
	loginURL := baseURL + "/login"
	formData := url.Values{
		"username": {cfg.Admin.Username},
		"password": {cfg.Admin.Password},
	}
	resp, err := http.PostForm(loginURL, formData)
	require.NoError(t, err, "Failed to post login form")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK from login")
	defer resp.Body.Close()

	// Step 2: Extract session token from response cookies
	cookies := resp.Cookies()
	require.NotEmpty(t, cookies, "Expected cookies to contain session token")
	sessionToken := cookies[0].Value
	require.Greater(t, len(sessionToken), 10, "Expected session token to be longer than 10 characters")

	// Step 3: Make an authenticated request using the session token
	apiURL := baseURL + "/users"
	req, err := http.NewRequest("GET", apiURL, nil)
	require.NoError(t, err, "Failed to create GET request")
	req.Header.Add("Cookie", sessionToken)

	client := &http.Client{}
	response, err := client.Do(req)
	require.NoError(t, err, "Failed to execute GET request")
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected HTTP 200 OK from API endpoint")
}
