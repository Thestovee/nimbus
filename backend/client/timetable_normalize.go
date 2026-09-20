package client

import (
	"context"
	"strconv"
	"time"
)

type NimbusLesson struct {
	LessonNo   string `json:"lessonNo"`
	Start      string `json:"start"`
	End        string `json:"end"`
	Subject    string `json:"subject"`
	Room       string `json:"room"`
	IsCanceled bool   `json:"isCanceled"`
}

func NormalizeTimetable(
	timetable LibrusTimetableResponse,
	classrooms LibrusClassroomsResponse,
) map[string][]NimbusLesson {
	roomNames := make(map[string]string, len(classrooms.Classrooms))
	for _, room := range classrooms.Classrooms {
		id := strconv.FormatInt(room.ID, 10)
		roomNames[id] = room.Name
	}

	result := make(map[string][]NimbusLesson, len(timetable.Timetable))

	for date, slots := range timetable.Timetable {
		result[date] = []NimbusLesson{}

		for _, lessons := range slots {
			for _, lesson := range lessons {
				roomName := roomNames[lesson.Classroom.Id]

				normal := NimbusLesson{
					LessonNo:   lesson.LessonNo,
					Start:      lesson.HourFrom,
					End:        lesson.HourTo,
					Subject:    lesson.Subject.Name,
					Room:       roomName,
					IsCanceled: lesson.IsCanceled,
				}

				result[date] = append(result[date], normal)
			}
		}
	}
	return result
}

func (c *LibrusClient) GetNimbusTimetable(
	ctx context.Context,
	weekStart time.Time,
) (map[string][]NimbusLesson, error) {
	timetable, err := c.GetTimetable(ctx, weekStart)
	if err != nil {
		return nil, err
	}
	classrooms, err := c.GetClassrooms(ctx)
	if err != nil {
		return nil, err
	}

	return NormalizeTimetable(*timetable, *classrooms), nil
}
