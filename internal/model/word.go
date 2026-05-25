package model

import "time"

type Word struct {
	ID          int64     `json:"id"`
	Word        string    `json:"word"`
	Translation string    `json:"translation"`
	Example     string    `json:"example,omitempty"`
	Difficulty  int       `json:"difficulty"`
	NextReview  time.Time `json:"next_review"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

const (
	StatusNew      = "new"
	StatusLearning = "learning"
	StatusMastered = "mastered"
)
