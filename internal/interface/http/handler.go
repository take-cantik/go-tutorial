package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	domaintask "github.com/take_cantik/go-tutorial/internal/domain/task"
	taskuc "github.com/take_cantik/go-tutorial/internal/usecase/task"
)

type Handler struct {
	uc *taskuc.UseCase
}

func NewHandler(uc *taskuc.UseCase) *Handler {
	return &Handler{uc: uc}
}

type taskResponse struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type createRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.uc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := make([]taskResponse, 0, len(tasks))
	for _, t := range tasks {
		res = append(res, toResponse(t))
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	task, err := h.uc.Get(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(task))
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	task, err := h.uc.Create(r.Context(), taskuc.CreateInput{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(task))
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	task, err := h.uc.Update(r.Context(), taskuc.UpdateInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(task))
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.uc.Delete(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", h.GetTask)
	mux.HandleFunc("PUT /tasks/{id}", h.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)
	return mux
}

func parseID(r *http.Request) (uint64, error) {
	return strconv.ParseUint(r.PathValue("id"), 10, 64)
}

func toResponse(t *domaintask.Task) taskResponse {
	return taskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
	}
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domaintask.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domaintask.ErrInvalidTitle), errors.Is(err, domaintask.ErrInvalidStatus):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
