package apicms

import "github.com/labstack/echo/v4"

// User API controller.
type User struct {
	svc UserService
}

// NewUser creates a new User controller.
func NewUser(svc UserService) *User {
	return &User{
		svc: svc,
	}
}

// RegisterHTTP register HTTP handlers based on actions for the service.
func (s *User) RegisterHTTP(r *echo.Group) {
	//+codegen=BindingApiHandler
}
