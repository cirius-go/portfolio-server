package servicecms

import (
	"github.com/cirius-go/portfolio-server/internal/service"
	"github.com/cirius-go/portfolio-server/internal/uow"
)

// Article is a service struct that encapsulates business logic.
type Article struct {
	service.Service
	uow uow.UnitOfWork
}

// NewArticle creates a new instance of Article service.
func NewArticle(uow uow.UnitOfWork) *Article {
	s := &Article{
		uow: uow,
	}
	return s
}
