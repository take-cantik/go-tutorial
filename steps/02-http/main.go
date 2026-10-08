package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Step 2: 標準ライブラリだけで HTTP API を作る
// 実行: make step2  または  go run ./steps/02-http
// 確認: curl http://localhost:8081/tasks

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type store struct {
	mu    sync.Mutex
	seq   int
	tasks []Task
}

func main() {
	s := &store{tasks: []Task{}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		writeJSON(w, http.StatusOK, s.tasks)
	})
	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		s.mu.Lock()
		s.seq++
		task := Task{ID: s.seq, Title: body.Title}
		s.tasks = append(s.tasks, task)
		s.mu.Unlock()

		writeJSON(w, http.StatusCreated, task)
	})

	addr := ":8081"
	fmt.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
