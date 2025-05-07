package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// NetworkingResourceConfig defines the configuration for networking resources.
type NetworkingResourceConfig struct {
	Cors struct {
		MaxAge           pulumi.Int  `pulumi:"maxAge"`
		AllowCredentials pulumi.Bool `pulumi:"allowCredentials"`
		AllowHeaders     []string    `pulumi:"allowHeaders"`
		AllowMethods     []string    `pulumi:"allowMethods"`
		AllowOrigins     []string    `pulumi:"allowOrigins"`
	} `pulumi:"cors"`
	DomainName struct {
		Domain         pulumi.String `pulumi:"domain"`
		CertArn        pulumi.String `pulumi:"certArn"`
		EndpointType   pulumi.String `pulumi:"endpointType"`
		SecurityPolicy pulumi.String `pulumi:"securityPolicy"`
	} `pulumi:"domainName"`
	Route53 struct {
		HostedZoneId pulumi.String `pulumi:"hostedZoneId"`
	} `pulumi:"route53"`
}

// NetworkingResources defines the netwoking resources required for the application.
type NetworkingResources struct {
	ApiGw              *apigatewayv2.Api
	ApiGwStage         *apigatewayv2.Stage
	ApiGwDomainName    *apigatewayv2.DomainName
	Route53ApiGwRecord *route53.Record
}

// CreateNetworkingServices creates the networking resources required for the application.
func createNetworkingServices(ctx *pulumi.Context) (*NetworkingResources, error) {
	var (
		cfg   = config.New(ctx, "")
		nwCfg = &NetworkingResourceConfig{}
	)

	if err := cfg.GetObject("networking", nwCfg); err != nil {
		return nil, err
	}
	var (
		res       = &NetworkingResources{}
		err       error
		corsCfg   = nwCfg.Cors
		domainCfg = nwCfg.DomainName
	)

	apiGwName := fmt.Sprintf("%s-apigw", namespace)
	res.ApiGw, err = apigatewayv2.NewApi(ctx, apiGwName, &apigatewayv2.ApiArgs{
		Name:         pulumi.String(apiGwName),
		ProtocolType: pulumi.String("HTTP"),
		CorsConfiguration: &apigatewayv2.ApiCorsConfigurationArgs{
			AllowCredentials: corsCfg.AllowCredentials,
			MaxAge:           corsCfg.MaxAge,
			AllowHeaders:     pulumi.ToStringArray(corsCfg.AllowHeaders),
			AllowMethods:     pulumi.ToStringArray(corsCfg.AllowMethods),
			AllowOrigins:     pulumi.ToStringArray(corsCfg.AllowOrigins),
		},
	})
	if err != nil {
		return nil, err
	}

	res.ApiGwStage, err = apigatewayv2.NewStage(ctx, fmt.Sprintf("%s-apigw-default-stage", namespace), &apigatewayv2.StageArgs{
		ApiId:      res.ApiGw.ID(),
		Name:       pulumi.String("$default"),
		AutoDeploy: pulumi.BoolPtr(true),
	})
	if err != nil {
		return nil, err
	}

	subDomain := fmt.Sprintf("%s-%s", ctx.Project(), ctx.Stack())
	domainName := pulumi.Sprintf("%s.%s", subDomain, nwCfg.DomainName.Domain)
	res.ApiGwDomainName, err = apigatewayv2.NewDomainName(ctx, fmt.Sprintf("%s-apigw-domain-name", namespace), &apigatewayv2.DomainNameArgs{
		DomainName: domainName,
		DomainNameConfiguration: &apigatewayv2.DomainNameDomainNameConfigurationArgs{
			CertificateArn: domainCfg.CertArn,
			EndpointType:   domainCfg.EndpointType,
			SecurityPolicy: domainCfg.SecurityPolicy,
		},
	})
	if err != nil {
		return nil, err
	}

	recordAliasArgs := (res.ApiGwDomainName.DomainNameConfiguration.ApplyT(func(c apigatewayv2.DomainNameDomainNameConfiguration) (route53.RecordAliasArray, error) {
		if c.TargetDomainName == nil || c.HostedZoneId == nil {
			return nil, fmt.Errorf("no target domain name or hosted zone id found")
		}

		return route53.RecordAliasArray{
			route53.RecordAliasArgs{
				EvaluateTargetHealth: pulumi.Bool(true),
				Name:                 pulumi.String(*c.TargetDomainName),
				ZoneId:               pulumi.String(*c.HostedZoneId),
			},
		}, nil
	})).(route53.RecordAliasArrayInput)
	res.Route53ApiGwRecord, err = route53.NewRecord(ctx, fmt.Sprintf("%s-route53-record", namespace), &route53.RecordArgs{
		ZoneId:  nwCfg.Route53.HostedZoneId,
		Name:    domainName,
		Type:    pulumi.String("A"),
		Aliases: recordAliasArgs,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
