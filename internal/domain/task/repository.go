package task

import "context"

// Repository は永続化の「契約」だけを定義します。
// 実装（MySQL など）はインフラ層に置き、ドメインは DB を知りません。
type Repository interface {
	Create(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id uint64) (*Task, error)
	FindAll(ctx context.Context) ([]*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id uint64) error
}
