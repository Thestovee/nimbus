package client

import (
	"reflect"
	"testing"
)

func TestNormalizeTimetable(t *testing.T) {
	timetable := LibrusTimetableResponse{
		Timetable: map[string][][]LibrusLesson{
			"2026-09-14": {
				{},
				{{LessonNo: "1", HourFrom: "08:00", HourTo: "08:45", Subject: LibrusSubject{Name: "Math"}, Classroom: Classroom{Id: "42"}}},
				{{LessonNo: "2", HourFrom: "08:55", HourTo: "09:40", Subject: LibrusSubject{Name: "Art"}, Classroom: Classroom{Id: "999"}, IsCanceled: true}},
			},
			"2026-09-15": {},
		},
	}
	classrooms := LibrusClassroomsResponse{
		Classrooms: []LibrusClassroom{{ID: 42, Name: "Room 42"}},
	}

	got := NormalizeTimetable(timetable, classrooms)
	want := map[string][]NimbusLesson{
		"2026-09-14": {
			{LessonNo: "1", Start: "08:00", End: "08:45", Subject: "Math", Room: "Room 42"},
			{LessonNo: "2", Start: "08:55", End: "09:40", Subject: "Art", Room: "", IsCanceled: true},
		},
		"2026-09-15": {},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NormalizeTimetable() = %#v, want %#v", got, want)
	}
}
