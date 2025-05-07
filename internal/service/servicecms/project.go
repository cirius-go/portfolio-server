package servicecms

import (
	"github.com/cirius-go/portfolio-server/internal/service"
	"github.com/cirius-go/portfolio-server/internal/uow"
)

// Project is a service struct that encapsulates business logic.
type Project struct {
	service.Service
	uow uow.UnitOfWork
	enf RBACEnforcer
}

// NewProject creates a new instance of Project service.
func NewProject(uow uow.UnitOfWork, enf RBACEnforcer) *Project {
	s := &Project{
		uow: uow,
		enf: enf,
	}
	return s
}
