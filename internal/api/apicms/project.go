package apicms

import "github.com/labstack/echo/v4"

// Project API controller.
type Project struct {
	svc ProjectService
}

// NewProject creates a new Project controller.
func NewProject(svc ProjectService) *Project {
	return &Project{
		svc: svc,
	}
}

// RegisterHTTP register HTTP handlers based on actions for the service.
func (s *Project) RegisterHTTP(r *echo.Group) {
	//+codegen=BindingApiHandler
}
