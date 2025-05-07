package servicecms

import (
	"github.com/cirius-go/portfolio-server/internal/service"
	"github.com/cirius-go/portfolio-server/internal/uow"
)

// User is a service struct that encapsulates business logic.
type User struct {
	service.Service
	uow uow.UnitOfWork
	enf RBACEnforcer
}

// NewUser creates a new instance of User service.
func NewUser(uow uow.UnitOfWork, enf RBACEnforcer) *User {
	s := &User{
		uow: uow,
		enf: enf,
	}
	return s
}
