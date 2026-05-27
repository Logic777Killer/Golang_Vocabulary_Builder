package repository

import (
	"context"
	"database/sql"
	"time"
	"vocab-app/internal/model"
)

type WordRepository struct {
	db *sql.DB
}

func NewWordRepository(db *sql.DB) *WordRepository {
	return &WordRepository{db: db}
}

func (r *WordRepository) Create(w *model.Word) error {
	query := `INSERT INTO words (word, translation, example, difficulty, next_review, status, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	if w.NextReview.IsZero() {
		w.NextReview = time.Now()
	}
	if w.Status == "" {
		w.Status = model.StatusNew
	}
	w.CreatedAt = time.Now()

	return r.db.QueryRowContext(context.Background(), query,
		w.Word, w.Translation, w.Example, w.Difficulty, w.NextReview, w.Status, w.CreatedAt).Scan(&w.ID)
}

func (r *WordRepository) GetAll() ([]*model.Word, error) {
	rows, err := r.db.QueryContext(context.Background(),
		"SELECT id, word, translation, example, difficulty, next_review, status, created_at FROM words")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*model.Word
	for rows.Next() {
		w := &model.Word{}
		if err := rows.Scan(&w.ID, &w.Word, &w.Translation, &w.Example, &w.Difficulty, &w.NextReview, &w.Status, &w.CreatedAt); err != nil {
			return nil, err
		}
		words = append(words, w)
	}
	return words, rows.Err()
}

func (r *WordRepository) GetForReview() (*model.Word, error) {
	w := &model.Word{}
	err := r.db.QueryRowContext(context.Background(),
		"SELECT id, word, translation, example, difficulty, next_review, status, created_at FROM words WHERE status != 'mastered' AND next_review <= NOW() LIMIT 1").
		Scan(&w.ID, &w.Word, &w.Translation, &w.Example, &w.Difficulty, &w.NextReview, &w.Status, &w.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *WordRepository) UpdateProgress(id int64, status string, next time.Time) error {
	_, err := r.db.ExecContext(context.Background(),
		"UPDATE words SET status = $1, next_review = $2 WHERE id = $3", status, next, id)
	return err
}
