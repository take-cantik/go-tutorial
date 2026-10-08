package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Step 4: ハンドラ / ユースケース / リポジトリに分離（クリーンアーキテクチャへの橋渡し）
// 実行: make step4

type Task struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

type TaskRepository interface {
	List(ctx context.Context) ([]Task, error)
	Create(ctx context.Context, title string) (Task, error)
}

type MySQLTaskRepository struct {
	db *sql.DB
}

func (r *MySQLTaskRepository) List(ctx context.Context) ([]Task, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title FROM tasks ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *MySQLTaskRepository) Create(ctx context.Context, title string) (Task, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO tasks (title, description, status) VALUES (?, '', 'todo')`, title)
	if err != nil {
		return Task{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Task{}, err
	}
	return Task{ID: uint64(id), Title: title}, nil
}

type TaskService struct {
	repo TaskRepository
}

func (s *TaskService) List(ctx context.Context) ([]Task, error) {
	return s.repo.List(ctx)
}

func (s *TaskService) Create(ctx context.Context, title string) (Task, error) {
	if title == "" {
		return Task{}, fmt.Errorf("title is required")
	}
	return s.repo.Create(ctx, title)
}

type Handler struct {
	service *TaskService
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, 200, tasks)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	task, err := h.service.Create(r.Context(), body.Title)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, 201, task)
}

func main() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		env("DB_USER", "tutorial"),
		env("DB_PASSWORD", "tutorial"),
		env("DB_HOST", "127.0.0.1"),
		env("DB_PORT", "3306"),
		env("DB_NAME", "go_tutorial"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	repo := &MySQLTaskRepository{db: db}
	service := &TaskService{repo: repo}
	handler := &Handler{service: service}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", handler.List)
	mux.HandleFunc("POST /tasks", handler.Create)

	fmt.Println("step4 listening on :8083 (layered)")
	log.Fatal(http.ListenAndServe(":8083", mux))
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
