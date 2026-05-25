package repository

import (
	"sync"
	"time"
	"vocab-app/internal/model"
)

type WordRepository struct {
	words  map[int64]*model.Word
	nextID int64
	mu     sync.RWMutex
}

func NewWordRepository() *WordRepository {
	return &WordRepository{
		words:  make(map[int64]*model.Word),
		nextID: 1,
	}
}

func (r *WordRepository) Create(w *model.Word) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	w.ID = r.nextID
	r.nextID++
	w.CreatedAt = time.Now()
	if w.NextReview.IsZero() {
		w.NextReview = w.CreatedAt
	}
	if w.Status == "" {
		w.Status = model.StatusNew
	}
	r.words[w.ID] = w
	return nil
}

func (r *WordRepository) GetAll() ([]*model.Word, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*model.Word, 0, len(r.words))
	for _, w := range r.words {
		res = append(res, w)
	}
	return res, nil
}

func (r *WordRepository) GetForReview() (*model.Word, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	now := time.Now()
	for _, w := range r.words {
		if w.Status != model.StatusMastered && !w.NextReview.After(now) {
			return w, nil
		}
	}
	return nil, nil
}

func (r *WordRepository) UpdateProgress(id int64, status string, next time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	w, ok := r.words[id]
	if !ok {
		return nil
	}
	w.Status = status
	w.NextReview = next
	return nil
}
