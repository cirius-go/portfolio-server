package servicecms

import (
	"github.com/cirius-go/portfolio-server/internal/service"
	"github.com/cirius-go/portfolio-server/internal/uow"
)

// User is a service struct that encapsulates business logic.
type User struct {
	service.Service
	uow uow.UnitOfWork
}

// NewUser creates a new instance of User service.
func NewUser(uow uow.UnitOfWork) *User {
	s := &User{
		uow: uow,
	}
	return s
}
