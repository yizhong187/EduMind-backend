package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
)

type Student struct {
	StudentID uuid.UUID `json:"student_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Valid     bool      `json:"valid"`
}

func DatabaseStudentToStudent(student database.Student) Student {
	return Student{
		StudentID: student.StudentID,
		Username:  student.Username,
		CreatedAt: student.CreatedAt,
		Name:      student.Name,
		Valid:     student.Valid,
	}
}
