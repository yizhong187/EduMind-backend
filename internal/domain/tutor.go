package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
)

type Tutor struct {
	TutorID     uuid.UUID `json:"tutor_id"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Valid       bool      `json:"valid"`
	YOE         int32     `json:"yoe"`
	Subject     string    `json:"subject"`
	Verified    bool      `json:"verified"`
	Rating      float64   `json:"rating"`
	RatingCount int32     `json:"rating_count"`
}

func DatabaseTutorToTutor(tutor database.Tutor) Tutor {
	return Tutor{
		TutorID:     tutor.TutorID,
		Username:    tutor.Username,
		CreatedAt:   tutor.CreatedAt,
		Name:        tutor.Name,
		Valid:       tutor.Valid,
		YOE:         tutor.Yoe,
		Subject:     tutor.Subject,
		Verified:    tutor.Verified,
		Rating:      tutor.Rating.Float64,
		RatingCount: tutor.RatingCount,
	}
}
