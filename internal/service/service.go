package service

import (
	"time"
	"vocab-app/internal/model"
	"vocab-app/internal/repository"
)

type WordService struct {
	repo     *repository.WordRepository
	interval time.Duration
}

func NewWordService(repo *repository.WordRepository, hours int) *WordService {
	return &WordService{repo, time.Duration(hours) * time.Hour}
}

func (s *WordService) AddWord(w *model.Word) error { return s.repo.Create(w) }

func (s *WordService) GetWordsForReview() ([]*model.Word, error) {
	w, err := s.repo.GetForReview()
	if err != nil || w == nil {
		return []*model.Word{}, nil
	}
	return []*model.Word{w}, nil
}

func (s *WordService) MarkReviewed(id int64, recalled bool) error {
	var next time.Time
	if recalled {
		next = time.Now().Add(s.interval)
	} else {
		next = time.Now().Add(1 * time.Hour)
	}
	return s.repo.UpdateProgress(id, model.StatusLearning, next)
}

func (s *WordService) GetAllWords() ([]*model.Word, error) { return s.repo.GetAll() }
