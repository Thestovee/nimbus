package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"nimbus/client"
	"nimbus/session"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTimetableHandlerRejectsUnknownSession(t *testing.T) {
	router := gin.New()
	router.GET("/timetable", TimetableHandler(session.NewStore()))

	request := httptest.NewRequest(http.MethodGet, "/timetable?weekStart=2026-09-14", nil)
	request.AddCookie(&http.Cookie{Name: "nimbus_session", Value: "unknown-token"})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
}

func TestTimetableHandlerRejectsInvalidWeekStart(t *testing.T) {
	sessions := session.NewStore()
	token, err := sessions.Create(client.NewClient())
	if err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	router.GET("/timetable", TimetableHandler(sessions))
	request := httptest.NewRequest(http.MethodGet, "/timetable?weekStart=not-a-date", nil)
	request.AddCookie(&http.Cookie{Name: "nimbus_session", Value: token})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestTimetableHandlerReturnsNormalizedLessons(t *testing.T) {
	fakeLibrus := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/Timetables":
			if got := r.URL.Query().Get("weekStart"); got != "2026-09-14" {
				t.Errorf("weekStart = %q, want 2026-09-14", got)
			}
			_, _ = w.Write([]byte(`{"Timetable":{"2026-09-14":[[{"LessonNo":"1","HourFrom":"08:00","HourTo":"08:45","Subject":{"Name":"Math"},"Classroom":{"Id":"42"}}]]}}`))
		case "/Classrooms":
			_, _ = w.Write([]byte(`{"Classrooms":[{"Id":42,"Name":"Room 42"}]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer fakeLibrus.Close()

	librusClient := client.NewClient()
	librusClient.APIBaseURL = fakeLibrus.URL
	librusClient.HTTPClient = fakeLibrus.Client()
	sessions := session.NewStore()
	token, err := sessions.Create(librusClient)
	if err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	router.GET("/timetable", TimetableHandler(sessions))
	request := httptest.NewRequest(http.MethodGet, "/timetable?weekStart=2026-09-14", nil)
	request.AddCookie(&http.Cookie{Name: "nimbus_session", Value: token})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
	var body struct {
		Days map[string][]client.NimbusLesson `json:"days"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	lessons := body.Days["2026-09-14"]
	if len(lessons) != 1 {
		t.Fatalf("expected one lesson, got %+v", lessons)
	}
	if lessons[0].Subject != "Math" || lessons[0].Room != "Room 42" {
		t.Errorf("unexpected normalized lesson: %+v", lessons[0])
	}
}
