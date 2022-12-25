package main

import (
	"github.com/manny-e1/snippetbox/internal/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, body := ts.get(t, "/ping")

	//rr := httptest.NewRecorder()
	//r, err := http.NewRequest(http.MethodGet, "/", nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//Ping(rr, r)
	//rs := rr.Result()
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: 200,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
			wantBody: "Not Found",
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
			wantBody: "Not Found",
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			statusCode, _, body := ts.get(t, test.urlPath)
			assert.Equal(t, statusCode, test.wantCode)
			assert.StringContains(t, body, test.wantBody)
		})
	}

}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	_, _, body := ts.get(t, "/user/signup")
	validCsrfToken := extractCSRFToken(t, body)

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action=\"/user/signup\" method=\"POST\" novalidate>"
	)
	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCsrfToken,
			wantCode:     http.StatusSeeOther,
		}, {
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "Wrong token",
			wantCode:     http.StatusBadRequest,
		}, {
			name:         "Empty Name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCsrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		}, {
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    validCsrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		}, {
			name:         "Empty password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCsrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		}, {
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "dagim@gmail.",
			userPassword: validPassword,
			csrfToken:    validCsrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		}, {
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "1234",
			csrfToken:    validCsrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		}, {
			name:         "Duplicate email",
			userName:     validName,
			userEmail:    "dagim@gmail.com",
			userPassword: validPassword,
			csrfToken:    validCsrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", test.userName)
			form.Add("email", test.userEmail)
			form.Add("password", test.userPassword)
			form.Add("csrf_token", test.csrfToken)
			code, _, _ := ts.postForm(t, "/user/signup", form)
			assert.Equal(t, code, test.wantCode)
			if test.wantFormTag != "" {
				assert.StringContains(t, body, test.wantFormTag)
			}
		},
		)
	}

}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	_, _, body := ts.get(t, "/user/login")
	validCsrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name          string
		wantCode      int
		wantLoginCode int
		wantLocation  string
		wantBody      string
		email         string
		password      string
		csrfToken     string
	}{
		{
			name:         "Unauthenticated",
			wantCode:     http.StatusSeeOther,
			wantLocation: "/user/login",
		},
		{
			name:          "Authenticated",
			wantCode:      http.StatusOK,
			wantBody:      "<form action='/snippet/create' method='POST'>",
			email:         "manny@gmail.com",
			password:      "12345678",
			wantLoginCode: http.StatusSeeOther,
			csrfToken:     validCsrfToken,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.wantBody == "" {
				code, header, _ := ts.get(t, "/snippet/create")
				assert.Equal(t, code, test.wantCode)
				assert.Equal(t, header.Get("Location"), test.wantLocation)
			} else {
				form := url.Values{}
				form.Add("email", test.email)
				form.Add("password", test.password)
				form.Add("csrf_token", test.csrfToken)
				loginStatusCode, _, _ := ts.postForm(t, "/user/login", form)
				assert.Equal(t, loginStatusCode, test.wantLoginCode)

				authSnippetCreateStatusCode, _, authSnippetCreateBody := ts.get(t, "/snippet/create")
				assert.Equal(t, authSnippetCreateStatusCode, test.wantCode)
				assert.StringContains(t, authSnippetCreateBody, test.wantBody)

			}
		})
	}
}
