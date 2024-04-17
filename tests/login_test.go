package tests

import (
	"net/http"
	"net/url"
	"strings"
	"tasklify/internal/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("API POST Form Requests", func() {
	var sessionToken string

	config := config.GetConfig()

	baseURL := "http://localhost:" + config.Port

	BeforeEach(func() {
		// Login to get session token
		loginUrl := baseURL + "/login"
		formData := url.Values{
			"username": {config.Admin.Username},
			"password": {config.Admin.Password},
		}
		resp, err := http.PostForm(loginUrl, formData)
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		// Extract the session token from the response
		// This is an example, you'll need to extract the token based on your actual response structure
		sessionToken = strings.Split(resp.Header.Get("Set-Cookie"), ";")[0]
	})

	Describe("Making authenticated requests", func() {
		It("Should make a successful POST request with session token", func() {
			apiUrl := baseURL + "/api/resource"
			form := url.Values{
				"key1": {"value1"},
				"key2": {"value2"},
			}
			req, _ := http.NewRequest("POST", apiUrl, strings.NewReader(form.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("Cookie", sessionToken)

			client := &http.Client{}
			response, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			// Further assertions can be added here based on response body
		})
	})
})
