package main

import "github.com/cirius-go/codegen"

func mkRepoGT() codegen.GroupTemplate {
	return codegen.GroupTemplate{
		Path:         "internal/repo{{ .entity | gopkg | lslash }}.go",
		Name:         "Repo",
		Description:  "Group of templates",
		RequiredArgs: []string{"entity"},
		Templates: []codegen.TemplateDefinition{
			{
				Description: "New Repo definition File",
				Path:        "internal/repo{{ .entity | gopkg | lslash }}.go",
				Name:        "InitRepoFile",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists: codegen.RuleOnFileNotExistsCreate,
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtInit,
				},
				ContentTmpl: `package repo
{{- $repoName := .entity | piCamel -}}
{{- $modelName := .entity | siCamel }}

// {{ $repoName }} Repo.
type {{ $repoName }} struct {
	db *gorm.DB
	*Common[model.{{ $modelName }}]
}

// New{{ $repoName }} Repository.
func New{{ $repoName }}(db *gorm.DB) *{{ $repoName }} {
	return &{{ $repoName }}{ db, NewCommon[model.{{ $modelName }}](db), }
}`,
				Output: []codegen.SimpleTemplateOutput{
					{
						Path: "internal/uow/interface.go",
						Name: "InitRepoInterface",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtEnd,
						},
						ContentTmpl: `{{- $repoName := .entity | piCamel -}}
{{- $modelName := .entity | siCamel }}

// {{ $repoName }} repo as a unit.
type {{ $repoName }} interface {
	Common[model.{{ $modelName }}]
}`,
					},
					{
						Path: "internal/uow/interface.go",
						Name: "InitUOWInterfaceHandler",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtPlaceholder,
							Placeholder:     "DefineUOWHandler",
						},
						ContentTmpl: `{{- $repoName := .entity | piCamel }}{{- $repoName := .entity | piCamel }}
{{ $repoName }}() {{ $repoName }}`,
					},
					{
						Path: "internal/uow/uow.go",
						Name: "ImplementUnit",
						Rule: codegen.TemplateDefinitionRule{
							AppendContentAt: codegen.RuleAppendContentAtEnd,
						},
						ContentTmpl: `{{- $repoName := .entity | piCamel }}

// {{ $repoName }} retrieve cached unit or init a new one.
func (u *uow) {{ $repoName }}() {{ $repoName }} {
	return lazyCache(u, "{{ $repoName }}", repo.New{{ $repoName }})
}`,
					},
				},
			},
		},
	}
}
