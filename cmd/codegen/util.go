package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os/exec"
	"regexp"
	"strings"

	"github.com/cirius-go/codegen"
)

// dummyImporter implements types.Importer and returns an empty package for any import.
type dummyImporter struct{}

func (d dummyImporter) Import(path string) (*types.Package, error) {
	return types.NewPackage(path, path), nil
}

// checkRedeclareError validates the Go source text and returns an error only if a
// redeclaration error is found.
func ValidateContent(src string) error {
	fset := token.NewFileSet()
	// Parse the source text.
	file, err := parser.ParseFile(fset, "src.go", src, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	pkgName := file.Name.Name

	// Slice to collect type checking errors.
	var typeErrors []error

	// Configure type checking with our dummy importer.
	conf := types.Config{
		Importer: dummyImporter{},
		Error: func(err error) {
			typeErrors = append(typeErrors, err)
		},
	}

	// Run type checking.
	_, err = conf.Check(pkgName, fset, []*ast.File{file}, nil)
	if err != nil {
		// Optionally, you might want to include the error from conf.Check.
		typeErrors = append(typeErrors, err)
	}

	// Check for redeclaration errors.
	for _, v := range typeErrors {
		errMsg := v.Error()
		if strings.Contains(errMsg, "redeclared in this block") {
			return err
		}
		if regexp.MustCompile(`method\s+(\S+)\s+already\s+declared\s+at\s+(\S+)`).MatchString(errMsg) {
			return err
		}
		// INFO: used to inspect kind of error
		// fmt.Println(errMsg)
	}

	return nil
}

func formatCode(c codegen.Context, w string) (string, error) {
	cmd := exec.Command("goimports", "-format-only")
	cmd.Stdin = strings.NewReader(w)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return w, nil
	}
	return out.String(), nil
}

func formatCodeWithImport(c codegen.Context, w string) (string, error) {
	cmd := exec.Command("goimports")
	cmd.Stdin = strings.NewReader(w)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return w, nil
	}
	return out.String(), nil
}

func mkTags(tags ...string) string {
	tagStrs := strings.Join(tags, " ")
	return fmt.Sprintf("`%s`", tagStrs)
}

func forwardAllArgs(a codegen.Args) codegen.Args {
	return a
}

func parseSwagRoute(v string) string {
	re := regexp.MustCompile(`{([^}]+)}`)
	return re.ReplaceAllStringFunc(v, func(match string) string {
		paramName := strings.Trim(match, "{}")
		return ":" + paramName
	})
}
