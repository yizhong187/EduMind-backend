package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/yizhong187/EduMind-backend/internal/database"
	"github.com/yizhong187/EduMind-backend/internal/util"
)

type Subject struct {
	SubjectID int32 `json:"subject_id"`
	Yoe       int32 `json:"yoe"`
}

type Tutor struct {
	TutorID     uuid.UUID `json:"tutor_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Valid       bool      `json:"valid"`
	Subjects    []Subject `json:"subjects"`
	Verified    bool      `json:"verified"`
	Rating      float64   `json:"rating"`
	RatingCount int32     `json:"rating_count"`
	PhotoURL    *string   `json:"photo_url"`
}

func DatabaseTutorToTutor(tutor database.Tutor, subjects []database.GetTutorSubjectsRow) Tutor {
	subjectsList := make([]Subject, len(subjects))
	for i, subject := range subjects {
		subjectsList[i] = Subject{
			SubjectID: subject.SubjectID,
			Yoe:       subject.Yoe,
		}
	}

	return Tutor{
		TutorID:     tutor.TutorID,
		Username:    tutor.Username,
		Email:       tutor.Email,
		CreatedAt:   tutor.CreatedAt,
		Name:        tutor.Name,
		Valid:       tutor.Valid,
		Verified:    tutor.Verified,
		Rating:      tutor.Rating.Float64,
		RatingCount: tutor.RatingCount,
		Subjects:    subjectsList,
		PhotoURL:    util.NullStringToString(tutor.PhotoUrl),
	}
}
