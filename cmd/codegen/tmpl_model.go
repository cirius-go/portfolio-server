package main

import "github.com/cirius-go/codegen"

func mkModelGT() codegen.GroupTemplate {
	modelGT := codegen.GroupTemplate{
		Path:         "internal/repo/model{{ .entity | gopkg | lslash }}.go",
		Name:         "Model",
		Description:  "Group of templates",
		RequiredArgs: []string{"entity"},
		Templates: []codegen.TemplateDefinition{
			{
				Description: "New Model definition File",
				Path:        "internal/repo/model{{ .entity | gopkg | lslash }}.go",
				Name:        "InitModelFile",
				Rule: codegen.TemplateDefinitionRule{
					OnFileNotExists: codegen.RuleOnFileNotExistsCreate,
					OnFileExist:     codegen.RuleOnFileExistsIgnore,
					AppendContentAt: codegen.RuleAppendContentAtInit,
				},
				ContentTmpl: `package model
// {{ .entity | sCamel }} model.
type {{ .entity | sCamel }} struct {
	Model {{ mkTags "gorm:\"embedded\"" }}
}`,
			},
		},
	}
	return modelGT
}
