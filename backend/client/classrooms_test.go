package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetClassrooms(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/Classrooms" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"Classrooms":[{"Id":42,"Name":"Room 42","Symbol":"42"}]}`))
	}))
	defer server.Close()

	c := NewClient()
	c.APIBaseURL = server.URL

	data, err := c.GetClassrooms(context.Background())
	if err != nil {
		t.Fatalf("GetClassrooms returned an error: %v", err)
	}
	if data == nil || len(data.Classrooms) != 1 {
		t.Fatalf("expected one classroom, got %+v", data)
	}
	room := data.Classrooms[0]
	if room.ID != 42 || room.Name != "Room 42" || room.Symbol != "42" {
		t.Errorf("unexpected classroom: %+v", room)
	}
}
