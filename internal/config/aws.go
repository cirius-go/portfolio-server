package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/cirius-go/portfolio-server/pkg/util"
)

var (
	lambdaFnName     = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	lambdaFnVer      = os.Getenv("AWS_LAMBDA_FUNCTION_VERSION")
	executionRoleArn = os.Getenv("AWS_EXEC_ROLE_ARN")
)

var (
	awsCfgOnce = sync.Once{}
	awsCfg     *aws.Config
)

// GetAWSConfig get the current aws configuration.
func GetAWSConfig(ctx context.Context) aws.Config {
	if awsCfg == nil {
		awsCfgOnce.Do(func() {
			v := util.MustE(config.LoadDefaultConfig(ctx))
			awsCfg = &v
		})
	}
	return *awsCfg
}

// IsInAWSLambda check if current env is aws lambda function.
func IsInAWSLambda() bool {
	return (lambdaFnName != "") && (lambdaFnVer != "")
}

// GetAWSConfigWithExecRole returns AWS config with execution role.
func GetAWSConfigWithExecRole(ctx context.Context) (aws.Config, error) {
	cfg := GetAWSConfig(ctx)
	if executionRoleArn == "" {
		return cfg, errors.New("AWS_EXEC_ROLE_ARN is not set")
	}
	stsClient := sts.NewFromConfig(cfg)
	creds := stscreds.NewAssumeRoleProvider(stsClient, executionRoleArn)
	cfg.Credentials = creds
	return cfg, nil
}

// LoadFromAPS loads the configuration from AWS parameter store.
// path: app/env
func LoadFromAPS(path string, vars map[string]string, nextToken *string) error {
	if path == "" {
		return fmt.Errorf("path prefix is required")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}
	ssmClient := ssm.NewFromConfig(cfg)
	resp, err := ssmClient.GetParametersByPath(context.Background(), &ssm.GetParametersByPathInput{
		NextToken:      nextToken,
		Path:           aws.String(path),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	for _, param := range resp.Parameters {
		paramName := strings.Replace(*param.Name, path, "", 1)
		vars[strings.ToUpper(paramName)] = *param.Value
	}

	if resp.NextToken != nil {
		return LoadFromAPS(path, vars, resp.NextToken)
	}

	return nil
}
