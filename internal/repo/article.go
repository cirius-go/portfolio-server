package repo

import (
	"github.com/cirius-go/portfolio-server/internal/repo/model"
	"gorm.io/gorm"
)

// Articles Repo.
type Articles struct {
	db *gorm.DB
	*Common[model.Article]
}

// NewArticles Repository.
func NewArticles(db *gorm.DB) *Articles {
	return &Articles{db, NewCommon[model.Article](db)}
}
