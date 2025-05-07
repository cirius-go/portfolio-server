package apicms

import "github.com/labstack/echo/v4"

// Article API controller.
type Article struct {
	svc ArticleService
}

// NewArticle creates a new Article controller.
func NewArticle(svc ArticleService) *Article {
	return &Article{
		svc: svc,
	}
}

// RegisterHTTP register HTTP handlers based on actions for the service.
func (s *Article) RegisterHTTP(r *echo.Group) {
	//+codegen=BindingApiHandler
}
