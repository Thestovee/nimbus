package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetTimetable(t *testing.T) {
	// Этот handler изображает сервер Librus.
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method GET, got %s", r.Method)
			}

			if r.URL.Path != "/Timetables" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}

			query := r.URL.Query()
			if got := query.Get("weekStart"); got != "2026-09-14" {
				t.Errorf("expected weekStart 2026-09-14, got %q", got)
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"Timetable":{"2026-09-14":[[{"LessonNo":"0"}]]}}`))
		},
	))
	defer server.Close()

	// Наш клиент отправит запрос fake-серверу, а не Librus.
	c := NewClient()
	c.APIBaseURL = server.URL

	weekStart := time.Date(2026, time.September, 14, 0, 0, 0, 0, time.UTC)
	data, err := c.GetTimetable(context.Background(), weekStart)
	if err != nil {
		t.Fatalf("GetTimetable returned an error: %v", err)
	}

	if data == nil {
		t.Fatal("GetTimetable returned no response")
	}
	slots := data.Timetable["2026-09-14"]
	if len(slots) != 1 || len(slots[0]) != 1 || slots[0][0].LessonNo != "0" {
		t.Errorf("unexpected timetable: %+v", data.Timetable)
	}
}

func TestGetTimetableRejectsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "private response", http.StatusBadGateway)
	}))
	defer server.Close()

	c := NewClient()
	c.APIBaseURL = server.URL
	weekStart := time.Date(2026, time.September, 14, 0, 0, 0, 0, time.UTC)
	data, err := c.GetTimetable(context.Background(), weekStart)
	if data != nil {
		t.Errorf("GetTimetable returned data: %v", data)
	}
	if err == nil {
		t.Fatal("GetTimetable did not return an error")
	}
	if strings.Contains(err.Error(), "private response") {
		t.Error("upstream response body leaked into error")
	}
}

func TestGetTimetableRejectsInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{broken`))
	}))
	defer server.Close()

	c := NewClient()
	c.APIBaseURL = server.URL
	weekStart := time.Date(2026, time.September, 14, 0, 0, 0, 0, time.UTC)
	data, err := c.GetTimetable(context.Background(), weekStart)
	if data != nil {
		t.Errorf("GetTimetable returned data: %v", data)
	}
	if err == nil {
		t.Errorf("GetTimetable did not return an error")
	}
}
