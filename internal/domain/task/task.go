package task

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrNotFound      = errors.New("task not found")
	ErrInvalidTitle  = errors.New("title is required")
	ErrInvalidStatus = errors.New("invalid status")
)

// Status はタスクの状態を表す値オブジェクト相当の型です。
type Status string

const (
	StatusTodo  Status = "todo"
	StatusDoing Status = "doing"
	StatusDone  Status = "done"
)

func ParseStatus(s string) (Status, error) {
	switch Status(s) {
	case StatusTodo, StatusDoing, StatusDone:
		return Status(s), nil
	default:
		return "", ErrInvalidStatus
	}
}

// Task はドメインの中心となるエンティティです。
// DB や HTTP の都合ではなく「業務上のタスク」を表現します。
type Task struct {
	ID          uint64
	Title       string
	Description string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func New(title, description string) (*Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrInvalidTitle
	}

	now := time.Now()
	return &Task{
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (t *Task) Update(title, description string, status Status) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return ErrInvalidTitle
	}
	if _, err := ParseStatus(string(status)); err != nil {
		return err
	}

	t.Title = title
	t.Description = description
	t.Status = status
	t.UpdatedAt = time.Now()
	return nil
}
