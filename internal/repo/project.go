package repo

import (
	"github.com/cirius-go/portfolio-server/internal/repo/model"
	"gorm.io/gorm"
)

// Projects Repo.
type Projects struct {
	db *gorm.DB
	*Common[model.Project]
}

// NewProjects Repository.
func NewProjects(db *gorm.DB) *Projects {
	return &Projects{db, NewCommon[model.Project](db)}
}
