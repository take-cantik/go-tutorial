package task

import (
	"context"

	domaintask "github.com/take_cantik/go-tutorial/internal/domain/task"
)

type UseCase struct {
	repo domaintask.Repository
}

func NewUseCase(repo domaintask.Repository) *UseCase {
	return &UseCase{repo: repo}
}

type CreateInput struct {
	Title       string
	Description string
}

func (u *UseCase) Create(ctx context.Context, in CreateInput) (*domaintask.Task, error) {
	task, err := domaintask.New(in.Title, in.Description)
	if err != nil {
		return nil, err
	}
	if err := u.repo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (u *UseCase) Get(ctx context.Context, id uint64) (*domaintask.Task, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *UseCase) List(ctx context.Context) ([]*domaintask.Task, error) {
	return u.repo.FindAll(ctx)
}

type UpdateInput struct {
	ID          uint64
	Title       string
	Description string
	Status      string
}

func (u *UseCase) Update(ctx context.Context, in UpdateInput) (*domaintask.Task, error) {
	task, err := u.repo.FindByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	status, err := domaintask.ParseStatus(in.Status)
	if err != nil {
		return nil, err
	}

	if err := task.Update(in.Title, in.Description, status); err != nil {
		return nil, err
	}
	if err := u.repo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (u *UseCase) Delete(ctx context.Context, id uint64) error {
	if _, err := u.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return u.repo.Delete(ctx, id)
}
