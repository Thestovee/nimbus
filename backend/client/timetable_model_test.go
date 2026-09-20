package client

import (
	"encoding/json"
	"testing"
)

func TestLibrusTimetableResponseJSON(t *testing.T) {
	raw := []byte(`{
  "Timetable": {
    "2026-09-14": [
      [],
      [{
        "LessonNo": "1",
        "HourFrom": "08:00",
        "HourTo": "08:45",
        "IsCanceled": false
      }]
    ]
  }
}`)

	var got LibrusTimetableResponse
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("decode timetable: %v", err)
	}

	slots, ok := got.Timetable["2026-09-14"]
	if !ok {
		t.Fatal("missing timetable day")
	}
	if len(slots) != 2 {
		t.Fatalf("expected 2 slots, got %d", len(slots))
	}
	if len(slots[0]) != 0 {
		t.Errorf("expected empty first slot, got %d lessons", len(slots[0]))
	}
	if len(slots[1]) != 1 {
		t.Fatalf("expected one lesson in second slot, got %d", len(slots[1]))
	}
	lesson := slots[1][0]
	if lesson.LessonNo != "1" || lesson.HourFrom != "08:00" || lesson.HourTo != "08:45" || lesson.IsCanceled {
		t.Errorf("unexpected decoded lesson: %+v", lesson)
	}
}
