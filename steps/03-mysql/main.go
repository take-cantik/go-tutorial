package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Step 3: MySQL と直接やりとりする CRUD
// 前提: make up で DB が起動していること
// 実行: make step3

type Task struct {
	ID     uint64 `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT id, title, status FROM tasks ORDER BY id`)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		tasks := []Task{}
		for rows.Next() {
			var t Task
			if err := rows.Scan(&t.ID, &t.Title, &t.Status); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			tasks = append(tasks, t)
		}
		writeJSON(w, 200, tasks)
	})

	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" {
			http.Error(w, "invalid", 400)
			return
		}

		res, err := db.Exec(`INSERT INTO tasks (title, description, status) VALUES (?, '', 'todo')`, body.Title)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		id, _ := res.LastInsertId()
		writeJSON(w, 201, Task{ID: uint64(id), Title: body.Title, Status: "todo"})
	})

	fmt.Println("step3 listening on :8082 (MySQL direct)")
	log.Fatal(http.ListenAndServe(":8082", mux))
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
