package uow

import (
	"context"
	"sync"

	"github.com/cirius-go/portfolio-server/internal/repo"
	"gorm.io/gorm"
)

// uow represents the Unit of Work.
type uow struct {
	db *gorm.DB

	mu     *sync.Mutex
	caches map[string]any
}

// New creates a new Unit of Work.
func New(db *gorm.DB) *uow {
	return &uow{
		db:     db,
		mu:     &sync.Mutex{},
		caches: make(map[string]any),
	}
}

// Transaction implements service.UnitOfWork.
func (u *uow) Transaction(ctx context.Context, txHandler func(ctx context.Context, uow UnitOfWork) error) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		return txHandler(ctx, New(tx))
	})
}

// Users retrieve cached unit or init a new one.
func (u *uow) Users() Users {
	return lazyCache(u, "Users", repo.NewUsers)
}

// Projects retrieve cached unit or init a new one.
func (u *uow) Projects() Projects {
	return lazyCache(u, "Projects", repo.NewProjects)
}

// Articles retrieve cached unit or init a new one.
func (u *uow) Articles() Articles {
	return lazyCache(u, "Articles", repo.NewArticles)
}
