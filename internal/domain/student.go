package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

type Student struct {
	StudentID uuid.UUID `json:"student_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Valid     bool      `json:"valid"`
	PhotoURL  *string   `json:"photo_url"`
}

func DatabaseStudentToStudent(student database.Student) Student {
	return Student{
		StudentID: student.StudentID,
		Username:  student.Username,
		Email:     student.Email,
		CreatedAt: student.CreatedAt,
		Name:      student.Name,
		Valid:     student.Valid,
		PhotoURL:  util.NullStringToString(student.PhotoUrl),
	}
}
