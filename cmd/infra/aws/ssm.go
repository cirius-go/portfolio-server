package main

import (
	"net/url"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func createParameterStoreSecrets(ctx *pulumi.Context) error {
	cfg := config.New(ctx, "")
	nestedSecretCfg := map[string]any{}
	if _, err := cfg.GetSecretObject("secrets", &nestedSecretCfg); err != nil {
		return err
	}
	secretCfg := FlatMapConfig(nestedSecretCfg)

	for k, v := range secretCfg {
		path, _ := url.JoinPath(ctx.Project(), ctx.Stack(), k)
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		if _, err := ssm.NewParameter(ctx, path, &ssm.ParameterArgs{
			Type:  pulumi.String("SecureString"),
			Name:  pulumi.String(path),
			Value: v.ToStringPtrOutput(),
		}); err != nil {
			return err
		}
	}

	return nil
}
