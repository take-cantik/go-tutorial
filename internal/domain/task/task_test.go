package task_test

import (
	"testing"

	"github.com/take_cantik/go-tutorial/internal/domain/task"
)

func TestNewTask(t *testing.T) {
	got, err := task.New("学ぶ", "Goを学ぶ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Status != task.StatusTodo {
		t.Fatalf("expected todo, got %s", got.Status)
	}
}

func TestNewTaskEmptyTitle(t *testing.T) {
	_, err := task.New("  ", "")
	if err != task.ErrInvalidTitle {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestUpdateStatus(t *testing.T) {
	got, err := task.New("学ぶ", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := got.Update("学ぶ", "進行中", task.StatusDoing); err != nil {
		t.Fatal(err)
	}
	if got.Status != task.StatusDoing {
		t.Fatalf("expected doing, got %s", got.Status)
	}
}
