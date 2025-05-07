package service

import (
	"context"

	"github.com/cirius-go/portfolio-server/pkg/errors"
	"github.com/cirius-go/portfolio-server/pkg/validator"
)

// Common errors.
var (
	ErrForbiddenAction = errors.NewForbidden(nil, "You don't have permission to perform this action")
)

// Service contains common service logics.
type Service struct {
}

// Validate validates the input struct.
func (s *Service) Validate(ctx context.Context, i any) error {
	if err := validator.Instance().StructCtx(ctx, i); err != nil {
		return errors.NewInvalidRequest(err, "invalid request")
	}
	return nil
}
