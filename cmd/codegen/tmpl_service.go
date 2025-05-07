package main

import "github.com/cirius-go/codegen"

func mkServiceImplGT() codegen.GroupTemplate {
	serviceGT := codegen.GroupTemplate{
		Path:         "internal/service/service{{ .subdomain |gopkg }}{{ .entity | gopkg | lslash }}.go",
		Name:         "ServiceImpl",
		Description:  "Group of templates",
		RequiredArgs: []string{"entity"},
		Templates: []codegen.TemplateDefinition{
			{
				Description: "New Service definition File",
				Path:        "internal/service/service{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitServiceFile",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists: codegen.RuleOnFileNotExistsCreate,
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtInit,
				},
				ContentTmpl: `package service{{ .subdomain | gopkg }}
{{- $service_name := .entity | siCamel }}
// {{ $service_name }} is a service struct that encapsulates business logic.
type {{ $service_name }} struct {
	service.Service
	uow uow.UnitOfWork
}

// New{{ $service_name }} creates a new instance of {{ $service_name }} service.
func New{{ $service_name }}(uow uow.UnitOfWork) *{{ $service_name }} {
	s := &{{ $service_name }}{
		uow: uow,
	}
	return s
}`,
			},
			{
				Description: "New Service Request & Response Definition",
				Path:        "internal/service/service{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitServiceHandler",
				Rule: codegen.TemplateDefinitionRule{
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtEnd,
				},
				ContentTmpl: `{{- $entityName := .entity | sCamel -}}
{{- $actionName := .action | sCamel -}}
{{- $dtoName := $actionName }}
{{- $dtoPackage := printf "dto%s" (.subdomain | gopkg) }}
{{- $apiPackage := printf "api%s" (.subdomain | gopkg) }}
{{- $requestType := printf "%s.%s%s" $dtoPackage $dtoName "Req" }}
{{- $responseType := printf "%s.%s%s" $dtoPackage $dtoName "Res" }}

// {{ $actionName }} implements {{ $apiPackage }}.{{ $entityName }}Service.
func (s *{{ $entityName }}) {{ $actionName }}(ctx context.Context, req *{{ $requestType }}) (*{{ $responseType }}, error) {
	// TODO: Implement {{ $actionName }} method
	panic("not implemented")
}`,
			},
		},
	}
	return serviceGT
}

func mkServiceGT() codegen.GroupTemplate {
	serviceGT := codegen.GroupTemplate{
		Path:         "internal/service/service{{ .subdomain |gopkg }}{{ .entity | gopkg | lslash }}.go",
		Name:         "Service",
		Description:  "Group of templates",
		RequiredArgs: []string{"entity"},
		Templates: []codegen.TemplateDefinition{
			{
				Description: "New Service definition File",
				Path:        "internal/service/service{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitServiceFile",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists: codegen.RuleOnFileNotExistsCreate,
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtInit,
				},
				ContentTmpl: `package service{{ .subdomain | gopkg }}
{{- $service_name := .entity | siCamel }}
// {{ $service_name }} is a service struct that encapsulates business logic.
type {{ $service_name }} struct {
	service.Service
	uow uow.UnitOfWork
}

// New{{ $service_name }} creates a new instance of {{ $service_name }} service.
func New{{ $service_name }}(uow uow.UnitOfWork) *{{ $service_name }} {
	s := &{{ $service_name }}{
		uow: uow,
	}
	return s
}`,
				Output: []codegen.SimpleTemplateOutput{
					{
						Path: "cmd/api/main.go",
						Name: "DeclSvc",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
							Placeholder:     "Define{{ .subdomain | siCamel }}Services",
						},
						ContentTmpl: `{{- $service_name := .entity | siCamel -}}
				{{ .entity | sLowerCamel }}Svc = service{{ .subdomain | gopkg }}.New{{ $service_name }}(unitOfWork)`,
					},
				},
			},
			{
				Description: "New Service Request & Response Definition",
				Path:        "internal/service/service{{ .subdomain | gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitServiceHandler",
				Rule: codegen.TemplateDefinitionRule{
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtEnd,
				},
				ContentTmpl: `{{- $entityName := .entity | sCamel -}}
{{- $actionName := .action | sCamel -}}
{{- $dtoName := printf "%s%s" $actionName $entityName }}
{{- $dtoPackage := printf "dto%s" (.subdomain | gopkg) }}
{{- $apiPackage := printf "api%s" (.subdomain | gopkg) }}
{{- $requestType := printf "%s.%s%s" $dtoPackage $dtoName "Req" }}
{{- $responseType := printf "%s.%s%s" $dtoPackage $dtoName "Res" }}

// {{ $actionName }} implements {{ $apiPackage }}.{{ $entityName }}Service.
func (s *{{ $entityName }}) {{ $actionName }}(ctx context.Context, req *{{ $requestType }}) (*{{ $responseType }}, error) {
	// TODO: Implement {{ $actionName }} method
	panic("not implemented")
}`,
			},
		},
	}
	return serviceGT
}
