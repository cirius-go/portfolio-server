package main

import (
	"github.com/cirius-go/codegen"
)

func mkDTOImplGT() codegen.GroupTemplate {
	dtoGT := codegen.GroupTemplate{
		Path:         "internal/dto/dto{{ .subdomain |gopkg  }}/impl.go",
		Name:         "DTOImpl",
		Description:  "Group of impl templates",
		RequiredArgs: []string{"entity"},
		Templates: []codegen.TemplateDefinition{
			{
				Description: "New DTO definition File",
				Path:        "internal/dto/dto{{ .subdomain |gopkg }}/impl.go",
				Name:        "InitDTOFile",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists: codegen.RuleOnFileNotExistsCreate,
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtInit,
				},
				ContentTmpl: `package dto{{ .subdomain | gopkg }}`,
			},
			{
				Description: "New DTO Request & Response Definition",
				Path:        "internal/dto/dto{{ .subdomain |gopkg }}/impl.go",
				Name:        "InitRqRp",
				Rule: codegen.TemplateDefinitionRule{
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtEnd,
				},
				ContentTmpl: `{{- $ident := .entity | sCamel -}}
{{- $dtoIdent := .action | sCamel -}}
{{- $req := printf "%s%s" $dtoIdent "Req" }}
{{- $res := printf "%s%s" $dtoIdent "Res" }}

type (
	// {{ $req }} is the request data of Impl.{{ $dtoIdent }}.
	{{ $req }} struct {}

	// {{ $res }} is the response data of Impl.{{ $dtoIdent }}.
	{{ $res }} struct {}
)`,
			},
		},
	}
	return dtoGT
}

func mkDTOGT() codegen.GroupTemplate {
	dtoGT := codegen.GroupTemplate{
		Path:         "internal/dto/dto{{ .subdomain |gopkg | lslash }}{{ .entity | gopkg | lslash }}.go",
		Name:         "DTO",
		Description:  "Group of templates",
		RequiredArgs: []string{"entity"},
		Templates: []codegen.TemplateDefinition{
			{
				Description: "New DTO definition File",
				Path:        "internal/dto/dto{{ .subdomain |gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitDTOFile",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists: codegen.RuleOnFileNotExistsCreate,
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtInit,
				},
				ContentTmpl: `package dto{{ .subdomain | gopkg }}`,
			},
			{
				Description: "New DTO Request & Response Definition",
				Path:        "internal/dto/dto{{ .subdomain |gopkg }}{{ .entity | gopkg | lslash }}.go",
				Name:        "InitRqRp",
				Rule: codegen.TemplateDefinitionRule{
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtEnd,
				},
				ContentTmpl: `{{- $ident := .entity | sCamel -}}
{{- $actionIdent := .action | sCamel -}}
{{- $dtoIdent := printf "%s%s" $actionIdent $ident -}}
{{- $req := printf "%s%s" $dtoIdent "Req" }}
{{- $res := printf "%s%s" $dtoIdent "Res" }}

type (
	// {{ $req }} is the request data of {{ $ident }}.{{ $actionIdent }}.
	{{- if eq $actionIdent "List" }}
	{{ $req }} struct {
		dto.ListingReq
		model.Filter{{ $ident }}RecIn{{ .subdomain | siCamel }}
	}
	{{- else if eq $actionIdent "Update" }}
	{{ $req }} struct {
		ID int {{ mkTags "param:\"id\"" }}
		model.Update{{ $ident }}DataIn{{ .subdomain | siCamel }}
	}
	{{- else if or  (eq $actionIdent "Delete") (eq $actionIdent "Get") }}
	{{ $req }} struct {
		ID int {{ mkTags "param:\"id\"" }}
	}
	{{- else }}
	{{ $req }} struct {}
	{{- end }}

	// {{ $res }} is the response data of {{ $ident }}.{{ $actionIdent }}.
	{{- if eq $actionIdent "List" }}
	{{ $res }} = dto.ListingRes[model.List{{ $ident }}RecIn{{ .subdomain | siCamel }}]
	{{- else if eq $actionIdent "Update" }}
	{{ $res }} = Get{{ $ident }}Res
	{{- else if eq $actionIdent "Create" }}
	{{ $res }} = Get{{ $ident }}Res
	{{- else }}
	{{ $res }} struct {}
	{{- end }}
)`,
				Output: []codegen.SimpleTemplateOutput{
					{
						Name: "InjectRequiredModel",
						Path: "internal/repo/model{{ .entity | gopkg | lslash }}.go",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtEnd,
						},
						ContentTmpl: `{{- $ident := .entity | sCamel -}}
{{- $pident := .entity | pSnake -}}
{{- $actionIdent := .action | sCamel -}}

{{- if eq $actionIdent "Update" }}
// Update{{ $ident }}DataIn{{ .subdomain | siCamel }} is used to update {{ $ident }} data.
type Update{{ $ident }}DataIn{{ .subdomain | siCamel }} struct {}
{{- else if eq $actionIdent "List" }}
// List{{ $ident }}RecIn{{ .subdomain | siCamel }} is used to list {{ $ident }} records.
type List{{ $ident }}RecIn{{ .subdomain | siCamel }} struct {}

func (*List{{ $ident }}RecIn{{ .subdomain | siCamel }}) TableName() string {
	return "{{ $pident }}"
}

// Filter{{ $ident }}RecIn{{ .subdomain | siCamel }} is used to filter {{ $ident }} records.
type Filter{{ $ident }}RecIn{{ .subdomain | siCamel }} struct {}
{{- end }}`,
					},
				},
			},
		},
	}
	return dtoGT
}
