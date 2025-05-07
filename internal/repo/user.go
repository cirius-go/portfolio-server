package repo

import (
	"github.com/cirius-go/portfolio-server/internal/repo/model"
	"gorm.io/gorm"
)

// Users Repo.
type Users struct {
	db *gorm.DB
	*Common[model.User]
}

// NewUsers Repository.
func NewUsers(db *gorm.DB) *Users {
	return &Users{db, NewCommon[model.User](db)}
}
