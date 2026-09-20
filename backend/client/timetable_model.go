package client

type LibrusTimetableResponse struct {
	Timetable map[string][][]LibrusLesson `json:"Timetable"`
}

type LibrusLesson struct {
	LessonNo   string `json:"LessonNo"`
	HourFrom   string `json:"HourFrom"`
	HourTo     string `json:"HourTo"`
	IsCanceled bool   `json:"IsCanceled"`
}
