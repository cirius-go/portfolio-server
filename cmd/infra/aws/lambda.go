package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/cirius-go/portfolio-server/util"
)

type LambdaConfig struct {
	Runtime         pulumi.String `pulumi:"runtime"`
	HandlerName     pulumi.String `pulumi:"handlerName"`
	Arch            pulumi.String `pulumi:"arch"`
	MemorySize      pulumi.Int    `pulumi:"memorySize"`
	ExecRoleAssumer pulumi.String `pulumi:"execRoleAssumer"`
}

type LambdaResourceConfig struct {
	LambdaConfig
	AppEnv             map[string]any                  `pulumi:"appEnv"`
	BuildWorkerPath    string                          `pulumi:"buildWorkerPath"`
	CustomApiLambda    *LambdaConfig                   `pulumi:"customApi"`
	CustomWorkerLambda map[pulumi.String]*LambdaConfig `pulumi:"customWorker"`
}

type LambdaResources struct {
	ExecutionRole     *iam.Role                          `pulumi:"executionRole"`
	AssumerRolePolicy *iam.RolePolicy                    `pulumi:"assumerRolePolicy"`
	LoggingPolicy     *iam.RolePolicy                    `pulumi:"loggingPolicy"`
	QuerySSMPolicy    *iam.RolePolicy                    `pulumi:"querySSMPolicy"`
	ApiLambda         *lambda.Function                   `pulumi:"apiLambda"`
	WorkerLambdas     map[pulumi.String]*lambda.Function `pulumi:"workerLambdas"`
}

func createLambdaResources(ctx *pulumi.Context) (*LambdaResources, error) {
	var (
		cfg  = config.New(ctx, "")
		lCfg = &LambdaResourceConfig{}
		res  = &LambdaResources{
			WorkerLambdas: map[pulumi.String]*lambda.Function{},
		}
		err error
	)

	if err := cfg.GetObject("lambda", lCfg); err != nil {
		return nil, err
	}

	executionRoleName := fmt.Sprintf("%s-lambda-execution-role", namespace)
	res.ExecutionRole, err = iam.NewRole(ctx, executionRoleName, &iam.RoleArgs{
		Name: pulumi.String(executionRoleName),
		AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Sid": "",
						"Effect": "Allow",
						"Principal": {
							"Service": "lambda.amazonaws.com"
						},
						"Action": "sts:AssumeRole"
					}
				]
			}`),
	})
	if err != nil {
		return nil, err
	}

	// in order to assume the execution role for an IAM role. then IAM role
	// will have permission to trigger some aws resources assigned to execution
	// role.
	if ctx.Stack() == "dev" {
		pulumi.All(res.ExecutionRole.Arn).ApplyT(func(args []any) error {
			var execRoleArn = args[0].(string)
			assumerRolePolicyName := fmt.Sprintf("%s-iam-assumer-lambda-execution-role", namespace)
			statement := map[string]any{
				"Version": "2012-10-17",
				"Statement": []map[string]any{
					{
						"Effect": "Allow",
						"Action": "sts:AssumeRole",
						"Resource": []string{
							execRoleArn,
						},
					},
				},
			}
			statementBytes, err := json.Marshal(&statement)
			if err != nil {
				return err
			}
			res.AssumerRolePolicy, err = iam.NewRolePolicy(ctx, assumerRolePolicyName, &iam.RolePolicyArgs{
				Name:   pulumi.String(assumerRolePolicyName),
				Role:   pulumi.String(lCfg.ExecRoleAssumer),
				Policy: pulumi.String(string(statementBytes)),
			})

			return err
		})

	}

	logPolicyName := fmt.Sprintf("%s-lambda-log-policy", namespace)
	res.LoggingPolicy, err = iam.NewRolePolicy(ctx, logPolicyName, &iam.RolePolicyArgs{
		Name: pulumi.String(logPolicyName),
		Role: res.ExecutionRole.Name,
		Policy: pulumi.String(`{
	               "Version": "2012-10-17",
	               "Statement": [{
	                   "Effect": "Allow",
	                   "Action": [
	                       "logs:CreateLogGroup",
	                       "logs:CreateLogStream",
	                       "logs:PutLogEvents"
	                   ],
	                   "Resource": "arn:aws:logs:*:*:*"
	               }]
	           }`),
	})
	if err != nil {
		return nil, err
	}

	querySSMPolicyName := fmt.Sprintf("%s-lambda-query-ssm-policy", namespace)
	res.QuerySSMPolicy, err = iam.NewRolePolicy(ctx, querySSMPolicyName, &iam.RolePolicyArgs{
		Name: pulumi.String(querySSMPolicyName),
		Role: res.ExecutionRole,
		Policy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Action": "ssm:GetParametersByPath",
					"Resource": "arn:aws:ssm:*:*:parameter/portfolio-server/*"
				}
			]
		}`),
	})

	var (
		runtime    = util.IfZero(lambda.RuntimeCustomAL2, lambda.Runtime(lCfg.Runtime))
		handler    = util.IfZero(pulumi.String("bootstrap"), lCfg.HandlerName)
		arch       = util.IfZero(pulumi.String("arm64"), lCfg.Arch)
		memorySize = util.IfZero(pulumi.Int(128), lCfg.MemorySize)
		envs       = lambda.FunctionEnvironmentArgs{
			Variables: FlatMapConfig(lCfg.AppEnv),
		}
	)

	apiFnName := fmt.Sprintf("%s-api", namespace)
	apiFnArgs := &lambda.FunctionArgs{
		Name:          pulumi.String(apiFnName),
		Runtime:       runtime,
		Architectures: pulumi.StringArray{arch},
		Handler:       handler,
		MemorySize:    memorySize,
		Environment:   envs,
		Role:          res.ExecutionRole.Arn,
		Code: pulumi.NewAssetArchive(map[string]any{
			".": pulumi.NewFileArchive("./.build/lambda/api.zip"),
		}),
	}
	mergeLambdaArgsWithConfig(apiFnArgs, lCfg.CustomApiLambda)
	res.ApiLambda, err = lambda.NewFunction(ctx, apiFnName, apiFnArgs, pulumi.DependsOn([]pulumi.Resource{
		res.LoggingPolicy, res.QuerySSMPolicy,
	}))
	if err != nil {
		return nil, err
	}

	workerNames, err := getListWorkerNames(lCfg.BuildWorkerPath)
	if err != nil {
		return nil, err
	}
	for _, workerName := range workerNames {
		var (
			workerFnName = fmt.Sprintf("%s-worker-%s", namespace, workerName)
			fnArgs       = &lambda.FunctionArgs{
				Name:          pulumi.String(workerFnName),
				Runtime:       runtime,
				Architectures: pulumi.StringArray{arch},
				Handler:       handler,
				Role:          res.ExecutionRole.Arn,
				MemorySize:    memorySize,
				Environment:   envs,
				Code: pulumi.NewAssetArchive(map[string]any{
					".": pulumi.NewFileArchive(fmt.Sprintf("./.build/lambda/%s-worker.zip", workerName)),
				}),
			}
			dependedOnPolicies = []pulumi.Resource{res.LoggingPolicy, res.QuerySSMPolicy}
		)

		mergeLambdaArgsWithConfig(fnArgs, lCfg.CustomWorkerLambda[pulumi.String(workerName)])
		res.WorkerLambdas[workerName], err = lambda.NewFunction(ctx, workerFnName, fnArgs, pulumi.DependsOn(dependedOnPolicies))
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func mergeLambdaArgsWithConfig(fnArgs *lambda.FunctionArgs, cfg *LambdaConfig) {
	if fnArgs == nil || cfg == nil {
		return
	}
	fnArgs.Runtime = util.IfZero[pulumi.StringPtrInput](fnArgs.Runtime, lambda.Runtime(cfg.Runtime))
	fnArgs.Handler = util.IfZero[pulumi.StringPtrInput](fnArgs.Handler, cfg.HandlerName)
	fnArgs.MemorySize = util.IfZero[pulumi.IntPtrInput](fnArgs.MemorySize, cfg.MemorySize)
	fnArgs.Architectures = util.IfZero[pulumi.StringArrayInput](fnArgs.Architectures, pulumi.StringArray{cfg.Arch})
}

func getListWorkerNames(buildPath string) ([]pulumi.String, error) {
	buildPath = util.IfZero("./.build/lambda", buildPath)

	workerNames := []pulumi.String{}
	if err := filepath.WalkDir(buildPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		workerName := filepath.Base(path)
		if workerName == "." {
			return nil
		}

		if strings.HasSuffix(workerName, "-worker.zip") {
			workerName = strings.TrimSuffix(workerName, "-worker.zip")
			workerNames = append(workerNames, pulumi.String(workerName))
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return workerNames, nil
}
