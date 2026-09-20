package client

type LibrusTimetableResponse struct {
	Timetable map[string][][]LibrusLesson `json:"Timetable"`
}

type LibrusLesson struct {
	LessonNo   string        `json:"LessonNo"`
	HourFrom   string        `json:"HourFrom"`
	HourTo     string        `json:"HourTo"`
	IsCanceled bool          `json:"IsCanceled"`
	Subject    LibrusSubject `json:"Subject"`
	Classroom  Classroom     `json:"Classroom"`
}

type LibrusSubject struct {
	Name  string
	Short string
}

type Classroom struct {
	Id  string
	Url string
}

type LibrusClassroomsResponse struct {
	Classrooms []LibrusClassroom `json:"Classrooms"`
}

type LibrusClassroom struct {
	ID     int64  `json:"Id"`
	Name   string `json:"Name"`
	Symbol string `json:"Symbol"`
}
