package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeAuthenticator struct {
	err      error
	username string
	password string
	calls    int
}

func (f *fakeAuthenticator) Authenticate(username, password string) error {
	f.username = username
	f.password = password
	f.calls++
	return f.err
}

func newLoginTestRouter(authenticator Authenticator) *gin.Engine {
	router := gin.New()
	router.POST("/login", LoginHandler(func() Authenticator {
		return authenticator
	}))

	return router
}

func TestLoginHandlerSuccess(t *testing.T) {
	fake := &fakeAuthenticator{}
	router := newLoginTestRouter(fake)
	body := `{"username":"student","password":"secret"}`
	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(body),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	if fake.username != "student" {
		t.Errorf(
			"expected username %q, got %q",
			"student",
			fake.username,
		)
	}

	if fake.password != "secret" {
		t.Errorf(
			"expected password %q, got %q",
			"secret",
			fake.password,
		)
	}
	if fake.calls != 1 {
		t.Errorf("expected Authenticate to be called once, got %d", fake.calls)
	}
}

func TestLoginHandlerRejectsInvalidJSON(t *testing.T) {
	fake := &fakeAuthenticator{}
	router := newLoginTestRouter(fake)
	body := `{`
	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(body),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}
	if fake.calls != 0 {
		t.Fatalf(
			"expected Authenticate to not be called, but got %d",
			fake.calls)
	}

}

func TestLoginHandlerRejectsAuthenticationError(t *testing.T) {
	fake := &fakeAuthenticator{
		err: errors.New("authentication failed"),
	}
	router := newLoginTestRouter(fake)
	body := `{"username":"student","password":"secret"}`
	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(body),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnauthorized,
			response.Code)
	}
	if fake.calls != 1 {
		t.Fatalf(
			"expected Authenticate to be called once, got %d",
			fake.calls)
	}
}
