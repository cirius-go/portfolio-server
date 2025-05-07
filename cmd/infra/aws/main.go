package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	// namespace format: %s-%s
	namespace = ""
)

func main() {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		namespace = fmt.Sprintf("%s-%s", ctx.Project(), ctx.Stack())
		if err := createParameterStoreSecrets(ctx); err != nil {
			return err
		}

		nwRes, err := createNetworkingServices(ctx)
		if err != nil {
			return err
		}

		fnRes, err := createLambdaResources(ctx)
		if err != nil {
			return err
		}

		_, err = integrateLambdaWithNetworking(ctx, nwRes, fnRes)
		if err != nil {
			return err
		}

		_, err = createCDNResources(ctx, fnRes)
		if err != nil {
			return err
		}

		return nil
	})
	panicIf(err)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
