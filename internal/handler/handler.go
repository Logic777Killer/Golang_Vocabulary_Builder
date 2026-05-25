package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"vocab-app/internal/model"
	"vocab-app/internal/service"
)

type Handler struct {
	svc    *service.WordService
	logger *slog.Logger
}

func NewHandler(svc *service.WordService, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/words", h.handleWords)
	mux.HandleFunc("/api/words/review", h.handleReview)
	mux.HandleFunc("/api/words/progress", h.handleProgress)
	mux.HandleFunc("/api/health", h.handleHealth)
}

func (h *Handler) handleWords(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.createWord(w, r)
		return
	}
	if r.Method == http.MethodGet {
		h.getAllWords(w, r)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *Handler) createWord(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Word, Translation, Example string
		Difficulty                 int
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	word := &model.Word{Word: input.Word, Translation: input.Translation, Example: input.Example, Difficulty: input.Difficulty}
	if err := h.svc.AddWord(word); err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	h.logger.Info("word created", "id", word.ID, "word", word.Word)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(word)
}

func (h *Handler) getAllWords(w http.ResponseWriter, r *http.Request) {
	words, _ := h.svc.GetAllWords()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(words)
}

func (h *Handler) handleReview(w http.ResponseWriter, r *http.Request) {
	words, _ := h.svc.GetWordsForReview()
	w.Header().Set("Content-Type", "application/json")
	if len(words) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	json.NewEncoder(w).Encode(words[0])
}

func (h *Handler) handleProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error":"missing id"}`, http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	var body struct{ Recalled bool }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if err := h.svc.MarkReviewed(id, body.Recalled); err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
