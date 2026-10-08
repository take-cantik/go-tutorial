package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/take_cantik/go-tutorial/internal/config"
	"github.com/take_cantik/go-tutorial/internal/infrastructure/mysql"
	httpapi "github.com/take_cantik/go-tutorial/internal/interface/http"
	taskuc "github.com/take_cantik/go-tutorial/internal/usecase/task"
)

func main() {
	cfg := config.Load()

	db, err := mysql.Open(cfg.DSN())
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer db.Close()

	repo := mysql.NewTaskRepository(db)
	uc := taskuc.NewUseCase(repo)
	handler := httpapi.NewHandler(uc)
	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s (env=%s)", cfg.HTTPPort, cfg.AppEnv)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	fmt.Println("server stopped")
}
