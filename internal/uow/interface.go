package uow

import (
	"context"

	"github.com/cirius-go/portfolio-server/internal/repo/model"
)

// TxHandler represents the transaction handler.
type TxHandler func(ctx context.Context, tx UnitOfWork) error

// UnitOfWork represents the unit of work.
type UnitOfWork interface {
	Transaction(ctx context.Context, txHandler func(ctx context.Context, tx UnitOfWork) error) error
}

// Common represents the common repository.
type Common[T any] interface {
	// Fetch all records at once.
	All(ctx context.Context) ([]*T, error)
	Create(ctx context.Context, m *T) error
	Get(ctx context.Context, m *T) error
	GetByID(ctx context.Context, id string) (*T, error)
	Update(ctx context.Context, id string, data any) error
	DeleteByID(ctx context.Context, id string) error
	HardDeleteByID(ctx context.Context, id string) error
}

// Users repo as a unit.
type Users interface {
	Common[model.User]
}

// Projects repo as a unit.
type Projects interface {
	Common[model.Project]
}

// Articles repo as a unit.
type Articles interface {
	Common[model.Article]
}
