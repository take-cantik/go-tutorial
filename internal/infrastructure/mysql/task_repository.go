package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domaintask "github.com/take_cantik/go-tutorial/internal/domain/task"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *domaintask.Task) error {
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO tasks (title, description, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)`,
		task.Title,
		task.Description,
		string(task.Status),
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}
	task.ID = uint64(id)
	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id uint64) (*domaintask.Task, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, title, description, status, created_at, updated_at
		 FROM tasks WHERE id = ?`,
		id,
	)

	task, err := scanTask(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domaintask.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find task: %w", err)
	}
	return task, nil
}

func (r *TaskRepository) FindAll(ctx context.Context) ([]*domaintask.Task, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, title, description, status, created_at, updated_at
		 FROM tasks ORDER BY id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]*domaintask.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *domaintask.Task) error {
	res, err := r.db.ExecContext(
		ctx,
		`UPDATE tasks
		 SET title = ?, description = ?, status = ?, updated_at = ?
		 WHERE id = ?`,
		task.Title,
		task.Description,
		string(task.Status),
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domaintask.ErrNotFound
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domaintask.ErrNotFound
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTask(s scanner) (*domaintask.Task, error) {
	var (
		task   domaintask.Task
		status string
	)
	if err := s.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}
	task.Status = domaintask.Status(status)
	return &task, nil
}
