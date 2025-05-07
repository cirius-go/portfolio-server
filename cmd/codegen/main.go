package main

import (
	"strings"

	"github.com/cirius-go/codegen"
	"github.com/cirius-go/codegen/pipeline"
	"github.com/spf13/cobra"
)

func main() {
	cg := initCodegen()

	rootCMD := cobra.Command{
		Use:               "codegen",
		Short:             "A code generator for Go",
		Args:              cobra.MaximumNArgs(1),
		DisableAutoGenTag: true,
	}

	newAPIModuleCMD := &cobra.Command{
		Use:  "api-module",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cg.SetArgs(codegen.Args{
				"subdomain": args[0],
				"entity":    args[1],
			})
			seq := NewAPIModule()
			if err := cg.BuildSeq(seq); err != nil {
				return err
			}
			if err := cg.Apply(); err != nil {
				return err
			}
			return nil
		},
	}

	newInternalModuleCMD := &cobra.Command{
		Use:  "internal-module",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cg.SetArgs(codegen.Args{
				"subdomain": args[0],
				"entity":    args[1],
			})
			seq := NewInternalModule()
			if err := cg.BuildSeq(seq); err != nil {
				return err
			}
			if err := cg.Apply(); err != nil {
				return err
			}
			return nil
		},
	}

	newAPIInternalModuleCMD := &cobra.Command{
		Use:  "api-internal-module",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cg.SetArgs(codegen.Args{
				"subdomain": args[0],
				"entity":    args[1],
			})
			seq := NewAPIInternalModule()
			if err := cg.BuildSeq(seq); err != nil {
				return err
			}
			if err := cg.Apply(); err != nil {
				return err
			}
			return nil
		},
	}

	newSelectAPICMD := &cobra.Command{
		Use:  "api-method",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cg.SetArgs(codegen.Args{
				"subdomain": args[0],
				"entity":    args[1],
			})
			seq := NewSelectAPIHandler()
			if err := cg.BuildSeq(seq); err != nil {
				return err
			}
			if err := cg.Apply(); err != nil {
				return err
			}
			return nil
		},
	}

	rootCMD.AddCommand(newAPIModuleCMD)
	rootCMD.AddCommand(newSelectAPICMD)
	rootCMD.AddCommand(newInternalModuleCMD)
	rootCMD.AddCommand(newAPIInternalModuleCMD)

	if err := rootCMD.Execute(); err != nil {
		panic(err)
	}
}

func initCodegen() *codegen.Codegen {
	c := codegen.NewConfig().
		SetStateDir(".codegen/state").
		AfterParsedTmplHook(formatCode).
		BeforeSaveHook(formatCodeWithImport).
		SetValidateContentHandler(ValidateContent)
	cg, err := codegen.NewWithConfig(c)
	if err != nil {
		panic(err)
	}

	{
		cs := pipeline.GetCasing()
		pl := pipeline.GetPluralize()
		pl.AddUncountableRule("cms")
		pl.AddUncountableRule("CMS")
		pl.AddIrregularRule("staff", "staffs")
		pl.AddIrregularRule("Staff", "Staffs")

		pipelines := pipeline.NewGoCollection()
		pipelines["gopkg"] = func(entity string) string {
			v := pl.Singular(entity)
			v = cs.Snake(entity)
			return strings.ReplaceAll(v, "_", "")
		}
		pipelines["parseSwagRoute"] = parseSwagRoute
		pipelines["lslash"] = func(v string) string {
			if v == "" {
				return v
			}

			if strings.HasPrefix(v, "/") {
				return v
			}
			return "/" + v
		}
		pipelines["mkTags"] = mkTags
		cg.UsePipelines(pipelines)
	}

	cg.RegisterTemplates(
		mkApiGT(),
		mkModelGT(),
		mkDTOGT(),
		mkDTOImplGT(),
		mkServiceGT(),
		mkServiceImplGT(),
		mkRepoGT(),
	)
	return cg
}
