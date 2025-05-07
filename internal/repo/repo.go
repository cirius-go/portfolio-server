package repo

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/cirius-go/portfolio-server/internal/repo/model"
)

// Common represents the repository.
type Common[Model any] struct {
	db *gorm.DB
}

func (r *Common[Model]) QuoteCol(name string) string {
	return QuoteCol(r.db, name)
}

// withCtx sets the context.
func (r *Common[Model]) withCtx(ctx context.Context) *gorm.DB {
	debug, ok := ctx.Value(model.ContextKeyDebug).(bool)
	if !ok || !debug {
		return r.db.WithContext(ctx)
	}

	logger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})

	return r.db.Session(&gorm.Session{
		Logger: logger,
	}).WithContext(ctx)
}

// Create creates a new record.
func (r *Common[Model]) Create(ctx context.Context, m *Model) error {
	return r.withCtx(ctx).Create(m).Error
}

// Get gets the record.
// get the first record that matches the given conditions, or return error if
// no record was found.
//
// Example: Get(ctx, &Model{ID: "foo"})
//
// Example 2: Get(ctx, &Model{Username: "cirius"})
func (r *Common[Model]) Get(ctx context.Context, m *Model) error {
	return r.withCtx(ctx).First(m).Error
}

// GetByID gets the record by ID.
//
// Example: GetByID(ctx, "foo")
func (r *Common[Model]) GetByID(ctx context.Context, id string) (*Model, error) {
	m := new(Model)
	err := r.withCtx(ctx).First(m, id).Error
	return m, err
}

// Update updates the record through struct.
func (r *Common[Model]) Update(ctx context.Context, id string, data any) error {
	m := new(Model)

	return r.withCtx(ctx).Model(m).Where("id = ?", id).Updates(data).Error
}

// Delete deletes the record.
func (r *Common[Model]) Delete(ctx context.Context, m *Model) error {
	return r.withCtx(ctx).Delete(m).Error
}

// DeleteByID deletes the record by ID.
func (r *Common[Model]) DeleteByID(ctx context.Context, id string) error {
	m := new(Model)
	return r.withCtx(ctx).Delete(m, id).Error
}

// HardDelete hard deletes the record.
func (r *Common[Model]) HardDelete(ctx context.Context, m *Model) error {
	return r.withCtx(ctx).Unscoped().Delete(m).Error
}

// HardDeleteByID hard deletes the record by ID.
func (r *Common[Model]) HardDeleteByID(ctx context.Context, id string) error {
	m := new(Model)
	return r.withCtx(ctx).Unscoped().Delete(m, id).Error
}

// Count counts the records.
func (r *Common[Model]) Count(ctx context.Context, conditions any, args ...any) (int64, error) {
	count := int64(0)
	if err := r.withCtx(ctx).Model(new(Model)).Where(conditions, args...).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// All returns all records.
func (r *Common[Model]) All(ctx context.Context) ([]*Model, error) {
	res := make([]*Model, 0)
	if err := r.withCtx(ctx).Model(new(Model)).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Common[Model]) AllWith(ctx context.Context, res any) error {
	if err := r.withCtx(ctx).Find(&res).Error; err != nil {
		return err
	}
	return nil
}

// NewCommon creates a new repository.
func NewCommon[Model any](db *gorm.DB) *Common[Model] {
	return &Common[Model]{db: db}
}

// ListingRequest represents the request to list the records.
type ListingRequest[T any] struct {
	Page    int
	PerPage int
	Filter  T
	Count   bool
	Sort    string
}
