package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type IntegrateLambdaWithNetworkingResources struct {
	ApiGwIntegration    *apigatewayv2.Integration
	ApiGwRouter         *apigatewayv2.Route
	ApiGwMapping        *apigatewayv2.ApiMapping
	ApiLambdaPermission *lambda.Permission
}

func integrateLambdaWithNetworking(ctx *pulumi.Context, nwResources *NetworkingResources, lambdaResources *LambdaResources) (*IntegrateLambdaWithNetworkingResources, error) {
	var (
		res = &IntegrateLambdaWithNetworkingResources{}
		err error
	)

	apiGwIntegrationName := fmt.Sprintf("%s-apigw-integration", namespace)
	res.ApiGwIntegration, err = apigatewayv2.NewIntegration(ctx, apiGwIntegrationName, &apigatewayv2.IntegrationArgs{
		ApiId:                nwResources.ApiGw.ID(),
		IntegrationUri:       lambdaResources.ApiLambda.InvokeArn,
		IntegrationType:      pulumi.String("AWS_PROXY"),
		PayloadFormatVersion: pulumi.String("2.0"),
	})
	if err != nil {
		return nil, err
	}
	ctx.Export("apiGwIntegrationId", res.ApiGwIntegration.ID())

	res.ApiGwRouter, err = apigatewayv2.NewRoute(ctx, fmt.Sprintf("%s-apigw-router", namespace), &apigatewayv2.RouteArgs{
		ApiId:    nwResources.ApiGw.ID(),
		RouteKey: pulumi.String("ANY /{proxy+}"),
		Target:   pulumi.Sprintf("integrations/%s", res.ApiGwIntegration.ID()),
	})
	if err != nil {
		return nil, err
	}
	ctx.Export("apiGwRouterId", res.ApiGwRouter.ID())

	res.ApiLambdaPermission, err = lambda.NewPermission(ctx, fmt.Sprintf("%s-apigw-invoke-lambda-permission", namespace), &lambda.PermissionArgs{
		Function:  lambdaResources.ApiLambda.Name,
		SourceArn: pulumi.Sprintf("%s/*/*", nwResources.ApiGw.ExecutionArn),
		Action:    pulumi.String("lambda:InvokeFunction"),
		Principal: pulumi.String("apigateway.amazonaws.com"),
	})
	if err != nil {
		return nil, err
	}

	// INFO: if status != 200 -> check sourceArn is correct or not.
	ctx.Export("testApiGwLambdaIntegration", pulumi.Sprintf("curl -X get '%s'", nwResources.ApiGw.ApiEndpoint))

	res.ApiGwMapping, err = apigatewayv2.NewApiMapping(ctx, fmt.Sprintf("%s-apigw-mapping", namespace), &apigatewayv2.ApiMappingArgs{
		ApiId:         nwResources.ApiGw.ID(),
		Stage:         nwResources.ApiGwStage.ID(),
		DomainName:    nwResources.ApiGwDomainName.DomainName,
		ApiMappingKey: pulumi.String(""),
	}, pulumi.DependsOn([]pulumi.Resource{res.ApiGwRouter}))
	if err != nil {
		return nil, err
	}

	return res, nil
}
