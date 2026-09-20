package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nimbus/client"
	"nimbus/session"

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

func newLoginTestRouter(authenticator *fakeAuthenticator) *gin.Engine {
	router := gin.New()
	router.POST("/login", LoginHandler(func(username, password string) (*client.LibrusClient, error) {
		if err := authenticator.Authenticate(username, password); err != nil {
			return nil, err
		}
		return client.NewClient(), nil
	}, session.NewStore()))

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

	cookies := response.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one session cookie, got %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != "nimbus_session" || cookie.Value == "" {
		t.Errorf("expected a nonempty nimbus_session cookie, got %q=%q", cookie.Name, cookie.Value)
	}
	if cookie.Path != "/" || cookie.MaxAge != 24*60*60 || !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("unexpected session cookie attributes: %+v", cookie)
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
	if cookies := response.Result().Cookies(); len(cookies) != 0 {
		t.Fatalf("expected no cookie after failed authentication, got %d", len(cookies))
	}
}

func TestLoginHandlerRejectsSessionCreationError(t *testing.T) {
	router := gin.New()
	router.POST("/login", LoginHandler(func(username, password string) (*client.LibrusClient, error) {
		return nil, nil
	}, session.NewStore()))

	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"student","password":"secret"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}
	if strings.Contains(response.Body.String(), "Logged in successfully") {
		t.Fatalf("success response followed a session error: %q", response.Body.String())
	}
}
