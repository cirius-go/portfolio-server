package main

import (
	"fmt"

	"github.com/cirius-go/codegen"
)

func mkApiGT() codegen.GroupTemplate {
	apiGT := codegen.GroupTemplate{
		Path:         "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
		Name:         "API",
		Description:  "Group of templates",
		RequiredArgs: []string{"entity"},
		Templates: []codegen.TemplateDefinition{
			{
				Description: "New API Service Interface",
				Path:        "internal/api/api{{ .subdomain | gopkg }}/interface.go",
				Name:        "InitAPIServiceInterfaceFile",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists:              codegen.RuleOnFileNotExistsCreate,
					OnFileExist:                  codegen.RuleOnFileExistsIgnore,
					AppendContentAt:              codegen.RuleAppendContentAtInit,
					AutoApplyOnValidationSuccess: true,
				},
				ContentTmpl: `package api{{ .subdomain | gopkg }}`,
			},
			{
				Description: "New API file",
				Path:        "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitAPIFile",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists:              codegen.RuleOnFileNotExistsCreate,
					OnFileExist:                  codegen.RuleOnFileExistsIgnore,
					AppendContentAt:              codegen.RuleAppendContentAtInit,
					AutoApplyOnValidationSuccess: true,
				},
				ContentTmpl: `package api{{ .subdomain | gopkg }}

import "github.com/labstack/echo/v4"

{{- $ident := .entity | sCamel }}

// {{ $ident }} API controller.
type {{ $ident }} struct {
	svc {{ $ident }}Service
}

// New{{ $ident }} creates a new {{ $ident }} controller.
func New{{ $ident }}(svc {{ $ident }}Service) *{{ $ident }} {
	return &{{ $ident }}{
		svc: svc,
	}
}

// RegisterHTTP register HTTP handlers based on actions for the service.
func (s *{{ $ident }}) RegisterHTTP(r *echo.Group) {
	//+codegen=BindingApiHandler
}`,
				Output: []codegen.SimpleTemplateOutput{
					{
						Path: "internal/api/api{{ .subdomain | gopkg }}/interface.go",
						Name: "ApiServiceInterface",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt:              codegen.RuleAppendContentAtEnd,
							AutoApplyOnValidationSuccess: true,
						},
						ContentTmpl: `{{- $ident := .entity | sCamel }}
// {{ $ident }}Service represents the service handler for {{ $ident }}.
type {{ $ident }}Service interface {
	//+codegen={{ $ident }}ServiceHandler
}`,
					},
					{
						Path: "cmd/api/main.go",
						Name: "DeclApi",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
							Placeholder:     "Define{{ .subdomain | siCamel }}APIs",
						},
						ContentTmpl: `{{- $api_name := .entity | siCamel -}}
					api{{ .subdomain | gopkg }}.New{{ $api_name }}({{ .entity | sLowerCamel }}Svc),`,
					},
				},
			},
			{
				Description: "New API Handler",
				Path:        "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitAPIHandler_POST",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists: codegen.RuleOnFileNotExistsError,
					MkFileNotExistsError: func(path string) error {
						return fmt.Errorf("file not exists: %s", path)
					},
					OnFileExist:                  codegen.RuleOnFileExistsIgnore,
					AppendContentAt:              codegen.RuleAppendContentAtEnd,
					AutoApplyOnValidationSuccess: true,
				},
				ContentTmpl: `{{- $ident := .entity | sCamel }}
{{- $actionIdent := .action | sCamel }}
{{- $subdomain := .subdomain | gopkg }}
{{- $route := .route | lslash -}}

// {{ $actionIdent }}
//
//	@id {{ $subdomain }}-{{ $ident | pKebab }}-{{ $actionIdent | sKebab | lower }}
//	@Summary {{ $actionIdent }}
//	@Description {{ $actionIdent }}
//	@Tags {{ $subdomain }}/{{ $ident | pKebab | lower }}
//	@Accept json
//	@Produce json
//	@Security BearerAuth
//	@Param Payload body dto{{ $subdomain }}.{{ $actionIdent }}{{ $ident }}Req true "JSON Request Payload"
//	@Success 200 {object} dto{{ $subdomain }}.{{ $actionIdent }}{{ $ident }}Res "JSON Response Payload"
//	@Failure 400 {object} dto.ErrorRes "JSON Response Payload"
//	@Failure 500 {object} dto.ErrorRes "JSON Response Payload"
//	@Router /{{ $subdomain }}/{{ $ident | pKebab | lower }}{{ $route }} [POST]
func (s *{{ $ident }}) {{ $actionIdent }}(c echo.Context) error {
	return api.MakeJSONHandler(c, s.svc.{{ $actionIdent }})
}`,
				Output: []codegen.SimpleTemplateOutput{
					{
						Path:     "internal/api/api{{ .subdomain | gopkg }}/interface.go",
						Name:     "NewServiceInterfaceHandler",
						NameTmpl: "{{ .entity | siCamel }}ServiceHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt:              codegen.RuleAppendContentAtPlaceholder,
							AutoApplyOnValidationSuccess: true,
						},
						ContentTmpl: `{{- $ident := .entity | sCamel -}}
{{- $actionIdent := .action | sCamel -}}
{{- $dtoIdent := printf "%s%s" $actionIdent $ident -}}
{{- $subdomain := .subdomain | gopkg }}
{{ $actionIdent }}(ctx context.Context, req *dto{{ $subdomain }}.{{ $dtoIdent }}Req) (res *dto{{ $subdomain }}.{{ $dtoIdent }}Res, err error)`,
					},
					{
						Path: "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
						Name: "BindingApiHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
						},
						ContentTmpl: `r.POST("{{ .route | lslash | parseSwagRoute }}", s.{{ .action | sCamel }})`,
					},
				},
			},
			{
				Description: "New API Handler",
				Path:        "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitAPIHandler_PUT",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists: codegen.RuleOnFileNotExistsError,
					MkFileNotExistsError: func(path string) error {
						return fmt.Errorf("file not exists: %s", path)
					},
					OnFileExist:                  codegen.RuleOnFileExistsIgnore,
					AppendContentAt:              codegen.RuleAppendContentAtEnd,
					AutoApplyOnValidationSuccess: true,
				},
				ContentTmpl: `{{- $ident := .entity | sCamel }}
{{- $actionIdent := .action | sCamel }}
{{- $subdomain := .subdomain | gopkg }}
{{- $route := .route | lslash -}}

// {{ $actionIdent }}
//
//	@id {{ $subdomain }}-{{ $ident | pKebab }}-{{ $actionIdent | sKebab | lower }}
//	@Summary {{ $actionIdent }}
//	@Description {{ $actionIdent }}
//	@Tags {{ $subdomain }}/{{ $ident | pKebab | lower }}
//	@Accept json
//	@Produce json
//	@Security BearerAuth
//	@Param ID path string true "ID"
//	@Param Payload body dto{{ $subdomain }}.{{ $actionIdent }}{{ $ident }}Req true "JSON Request Payload"
//	@Success 200 {object} dto{{ $subdomain }}.{{ $actionIdent }}{{ $ident }}Res "JSON Response Payload"
//	@Failure 400 {object} dto.ErrorRes "JSON Response Payload"
//	@Failure 500 {object} dto.ErrorRes "JSON Response Payload"
//	@Router /{{ $subdomain }}/{{ $ident | pKebab | lower }}{{ $route }} [PUT]
func (s *{{ $ident }}) {{ $actionIdent }}(c echo.Context) error {
	return api.MakeJSONHandler(c, s.svc.{{ $actionIdent }})
}`,
				Output: []codegen.SimpleTemplateOutput{
					{
						Path:     "internal/api/api{{ .subdomain | gopkg }}/interface.go",
						Name:     "NewServiceInterfaceHandler",
						NameTmpl: "{{ .entity | siCamel }}ServiceHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt:              codegen.RuleAppendContentAtPlaceholder,
							AutoApplyOnValidationSuccess: true,
						},
						ContentTmpl: `{{- $ident := .entity | sCamel -}}
{{- $actionIdent := .action | sCamel -}}
{{- $dtoIdent := printf "%s%s" $actionIdent $ident -}}
{{- $subdomain := .subdomain | gopkg }}
{{ $actionIdent }}(ctx context.Context, req *dto{{ $subdomain }}.{{ $dtoIdent }}Req) (res *dto{{ $subdomain }}.{{ $dtoIdent }}Res, err error)`,
					},
					{
						Path: "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
						Name: "BindingApiHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
						},
						ContentTmpl: `r.PUT("{{ .route | lslash | parseSwagRoute }}", s.{{ .action | sCamel }})`,
					},
				},
			},
			{
				Description: "New API Handler",
				Path:        "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitAPIHandler_DELETE",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists: codegen.RuleOnFileNotExistsError,
					MkFileNotExistsError: func(path string) error {
						return fmt.Errorf("file not exists: %s", path)
					},
					OnFileExist:                  codegen.RuleOnFileExistsIgnore,
					AppendContentAt:              codegen.RuleAppendContentAtEnd,
					AutoApplyOnValidationSuccess: true,
				},
				ContentTmpl: `{{- $ident := .entity | sCamel }}
{{- $actionIdent := .action | sCamel }}
{{- $subdomain := .subdomain | gopkg }}
{{- $route := .route | lslash -}}

// {{ $actionIdent }}
//
//	@id {{ $subdomain }}-{{ $ident | pKebab }}-{{ $actionIdent | sKebab | lower }}
//	@Summary {{ $actionIdent }}
//	@Description {{ $actionIdent }}
//	@Tags {{ $subdomain }}/{{ $ident | pKebab | lower }}
//	@Accept json
//	@Produce json
//	@Security BearerAuth
//	@Param ID path string true "ID"
//	@Param Payload body dto{{ $subdomain }}.{{ $actionIdent }}{{ $ident }}Req true "JSON Request Payload"
//	@Success 200 {object} dto{{ $subdomain }}.{{ $actionIdent }}{{ $ident }}Res "JSON Response Payload"
//	@Failure 400 {object} dto.ErrorRes "JSON Response Payload"
//	@Failure 500 {object} dto.ErrorRes "JSON Response Payload"
//	@Router /{{ $subdomain }}/{{ $ident | pKebab | lower }}{{ $route }} [DELETE]
func (s *{{ $ident }}) {{ $actionIdent }}(c echo.Context) error {
	return api.MakeJSONHandler(c, s.svc.{{ $actionIdent }})
}`,
				Output: []codegen.SimpleTemplateOutput{
					{
						Path:     "internal/api/api{{ .subdomain | gopkg }}/interface.go",
						Name:     "NewServiceInterfaceHandler",
						NameTmpl: "{{ .entity | siCamel }}ServiceHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
						},
						ContentTmpl: `{{- $ident := .entity | sCamel -}}
{{- $actionIdent := .action | sCamel -}}
{{- $dtoIdent := printf "%s%s" $actionIdent $ident -}}
{{- $subdomain := .subdomain | gopkg }}
{{ $actionIdent }}(ctx context.Context, req *dto{{ $subdomain }}.{{ $dtoIdent }}Req) (res *dto{{ $subdomain }}.{{ $dtoIdent }}Res, err error)`,
					},
					{
						Path: "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
						Name: "BindingApiHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
						},
						ContentTmpl: `r.DELETE("{{ .route | lslash | parseSwagRoute }}", s.{{ .action | sCamel }})`,
					},
				},
			},
			{
				Description: "New API Handler",
				Path:        "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitAPIHandler_GET",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists: codegen.RuleOnFileNotExistsError,
					MkFileNotExistsError: func(path string) error {
						return fmt.Errorf("file not exists: %s", path)
					},
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtEnd,
				},
				ContentTmpl: `{{- $ident := .entity | sCamel }}
{{- $actionIdent := .action | sCamel }}
{{- $subdomain := .subdomain | gopkg }}
{{- $route := .route | lslash -}}

// {{ $actionIdent }}
//
//	@id {{ $subdomain }}-{{ $ident | pKebab }}-{{ $actionIdent | sKebab | lower }}
//	@Summary {{ $actionIdent }}
//	@Description {{ $actionIdent }}
//	@Tags {{ $subdomain }}/{{ $ident | pKebab | lower }}
//	@Accept json
//	@Produce json
//	@Security BearerAuth
//	@Param ID path string true "ID"
//	@Success 200 {object} dto{{ $subdomain }}.{{ $actionIdent }}{{ $ident }}Res "JSON Response Payload"
//	@Failure 400 {object} dto.ErrorRes "JSON Response Payload"
//	@Failure 500 {object} dto.ErrorRes "JSON Response Payload"
//	@Router /{{ $subdomain }}/{{ $ident | pKebab | lower }}{{ $route }} [GET]
func (s *{{ $ident }}) {{ $actionIdent }}(c echo.Context) error {
	return api.MakeJSONHandler(c, s.svc.{{ $actionIdent }})
}`,
				Output: []codegen.SimpleTemplateOutput{
					{
						Path:     "internal/api/api{{ .subdomain | gopkg }}/interface.go",
						Name:     "NewServiceInterfaceHandler",
						NameTmpl: "{{ .entity | siCamel }}ServiceHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
						},
						ContentTmpl: `{{- $ident := .entity | sCamel -}}
{{- $actionIdent := .action | sCamel -}}
{{- $dtoIdent := printf "%s%s" $actionIdent $ident -}}
{{- $subdomain := .subdomain | gopkg }}
{{ $actionIdent }}(ctx context.Context, req *dto{{ $subdomain }}.{{ $dtoIdent }}Req) (res *dto{{ $subdomain }}.{{ $dtoIdent }}Res, err error)`,
					},
					{
						Path: "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
						Name: "BindingApiHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
						},
						ContentTmpl: `r.GET("{{ .route | lslash | parseSwagRoute }}", s.{{ .action | sCamel }})`,
					},
				},
			},
			{
				Description: "New API Handler",
				Path:        "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitAPIHandler_PATCH",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists: codegen.RuleOnFileNotExistsError,
					MkFileNotExistsError: func(path string) error {
						return fmt.Errorf("file not exists: %s", path)
					},
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtEnd,
				},
				ContentTmpl: `{{- $ident := .entity | sCamel }}
{{- $actionIdent := .action | sCamel }}
{{- $subdomain := .subdomain | gopkg }}
{{- $route := .route | lslash -}}

// {{ $actionIdent }}
//
//	@id {{ $subdomain }}-{{ $ident | pKebab }}-{{ $actionIdent | sKebab | lower }}
//	@Summary {{ $actionIdent }}
//	@Description {{ $actionIdent }}
//	@Tags {{ $subdomain }}/{{ $ident | pKebab | lower }}
//	@Accept json
//	@Produce json
//	@Security BearerAuth
//	@Param ID path string true "ID"
//	@Param Payload body dto{{ $subdomain }}.{{ $actionIdent }}{{ $ident }}Req true "JSON Request Payload"
//	@Success 200 {object} dto{{ $subdomain }}.{{ $actionIdent }}{{ $ident }}Res "JSON Response Payload"
//	@Failure 400 {object} dto.ErrorRes "JSON Response Payload"
//	@Failure 500 {object} dto.ErrorRes "JSON Response Payload"
//	@Router /{{ $subdomain }}/{{ $ident | pKebab | lower }}{{ $route }} [PATCH]
func (s *{{ $ident }}) {{ $actionIdent }}(c echo.Context) error {
	return api.MakeJSONHandler(c, s.svc.{{ $actionIdent }})
}`,
				Output: []codegen.SimpleTemplateOutput{
					{
						Path:     "internal/api/api{{ .subdomain | gopkg }}/interface.go",
						Name:     "NewServiceInterfaceHandler",
						NameTmpl: "{{ .entity | siCamel }}ServiceHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
						},
						ContentTmpl: `{{- $ident := .entity | sCamel -}}
{{- $actionIdent := .action | sCamel -}}
{{- $dtoIdent := printf "%s%s" $actionIdent $ident -}}
{{- $subdomain := .subdomain | gopkg }}
{{ $actionIdent }}(ctx context.Context, req *dto{{ $subdomain }}.{{ $dtoIdent }}Req) (res *dto{{ $subdomain }}.{{ $dtoIdent }}Res, err error)`,
					},
					{
						Path: "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
						Name: "BindingApiHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
						},
						ContentTmpl: `r.PATCH("{{ .route | lslash | parseSwagRoute }}", s.{{ .action | sCamel }})`,
					},
				},
			},
			{
				Description: "New API Handler",
				Path:        "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitAPIHandler_LIST",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists: codegen.RuleOnFileNotExistsError,
					MkFileNotExistsError: func(path string) error {
						return fmt.Errorf("file not exists: %s", path)
					},
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtEnd,
				},
				ContentTmpl: `{{- $ident := .entity | sCamel }}
{{- $actionIdent := .action | sCamel }}
{{- $subdomain := .subdomain | gopkg }}
{{- $route := .route | lslash -}}

// {{ $actionIdent }}
//
//	@id {{ $subdomain }}-{{ $ident | pKebab }}-{{ $actionIdent | sKebab | lower }}
//	@Summary {{ $actionIdent }}
//	@Description {{ $actionIdent }}
//	@Tags {{ $subdomain }}/{{ $ident | pKebab | lower }}
//	@Accept json
//	@Produce json
//	@Security BearerAuth
//	@Success 200 {object} dto{{ $subdomain }}.{{ $actionIdent }}{{ $ident }}Res "JSON Response Payload"
//	@Failure 400 {object} dto.ErrorRes "JSON Response Payload"
//	@Failure 500 {object} dto.ErrorRes "JSON Response Payload"
//	@Router /{{ $subdomain }}/{{ $ident | pKebab | lower }}{{ $route }} [GET]
func (s *{{ $ident }}) {{ $actionIdent }}(c echo.Context) error {
	return api.MakeJSONHandler(c, s.svc.{{ $actionIdent }})
}`,
				Output: []codegen.SimpleTemplateOutput{
					{
						Path:     "internal/api/api{{ .subdomain | gopkg }}/interface.go",
						Name:     "NewServiceInterfaceHandler",
						NameTmpl: "{{ .entity | siCamel }}ServiceHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
						},
						ContentTmpl: `{{- $ident := .entity | sCamel -}}
{{- $actionIdent := .action | sCamel -}}
{{- $dtoIdent := printf "%s%s" $actionIdent $ident -}}
{{- $subdomain := .subdomain | gopkg }}
{{ $actionIdent }}(ctx context.Context, req *dto{{ $subdomain }}.{{ $dtoIdent }}Req) (res *dto{{ $subdomain }}.{{ $dtoIdent }}Res, err error)`,
					},
					{
						Path: "internal/api/api{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
						Name: "BindingApiHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
						},
						ContentTmpl: `r.GET("{{ .route | lslash | parseSwagRoute }}", s.{{ .action | sCamel }})`,
					},
				},
			},
		},
	}

	return apiGT
}
