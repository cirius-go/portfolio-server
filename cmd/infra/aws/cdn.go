package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type CDNConfig struct {
	Name                    pulumi.String `pulumi:"name"`
	CachePolicyID           pulumi.String `pulumi:"cachePolicyID"`
	PriceClass              pulumi.String `pulumi:"priceClass"`
	ACMCertARN              pulumi.String `pulumi:"acmCertARN"`
	ResponseHeadersPolicyID pulumi.String `pulumi:"responseHeadersPolicyID"`
	TempDir                 pulumi.String `pulumi:"tempDir"`
	TempExpDays             pulumi.Int    `pulumi:"tempExpDays"`
}

type CDNResourceConfig struct {
	Assets CDNConfig `pulumi:"assets"`
}

type CDNRoute53Config struct {
	HostedZoneID pulumi.String `pulumi:"hostedZoneID"`
}

type CDNResources struct {
	AssetBucket                  *s3.BucketV2                       `pulumi:"assetBucket"`
	AssetBucketPolicy            *s3.BucketPolicy                   `pulumi:"assetBucketPolicy"`
	AssetBucketOwnership         *s3.BucketOwnershipControls        `pulumi:"assetBucketOwnership"`
	AssetBucketAccessBlock       *s3.BucketPublicAccessBlock        `pulumi:"assetBucketAccessBlock"`
	AssetBucketTempFileLifeCycle *s3.BucketLifecycleConfigurationV2 `pulumi:"assetBucketTempFileLifeCycle"`
	AssetCDN                     *cloudfront.Distribution           `pulumi:"assetCDN"`
	AssetCDNOAC                  *cloudfront.OriginAccessControl    `pulumi:"assetCDNOAC"`
	AssetRoute53Record           *route53.Record                    `pulumi:"assetRoute53Record"`
}

func createCDNResources(ctx *pulumi.Context, fnRes *LambdaResources) (res *CDNResources, err error) {
	var (
		cfg   = config.New(ctx, "")
		cCfg  = &CDNResourceConfig{}
		nwCfg = &NetworkingResourceConfig{}
	)
	res = &CDNResources{}

	if err = cfg.GetObject("cdn", cCfg); err != nil {
		return nil, err
	}
	if err = cfg.GetObject("networking", nwCfg); err != nil {
		return nil, err
	}

	bucketName := fmt.Sprintf("%s-%s-bucket", namespace, cCfg.Assets.Name)
	res.AssetBucket, err = s3.NewBucketV2(ctx, bucketName, &s3.BucketV2Args{
		Bucket: pulumi.String(bucketName),
	})
	if err != nil {
		return nil, err
	}

	res.AssetBucketTempFileLifeCycle, err = s3.NewBucketLifecycleConfigurationV2(ctx, fmt.Sprintf("%s-%s-temp-file-lifecycle", namespace, cCfg.Assets.Name), &s3.BucketLifecycleConfigurationV2Args{
		Bucket: res.AssetBucket.ID(),
		Rules: s3.BucketLifecycleConfigurationV2RuleArray{
			s3.BucketLifecycleConfigurationV2RuleArgs{
				Id:                             pulumi.Sprintf("%s-%s-temp-lifecycle-rule", namespace, cCfg.Assets.Name),
				AbortIncompleteMultipartUpload: nil,
				Expiration: s3.BucketLifecycleConfigurationV2RuleExpirationArgs{
					Date: nil,
					Days: cCfg.Assets.TempExpDays,
				},
				Filter: s3.BucketLifecycleConfigurationV2RuleFilterArgs{
					Prefix: cCfg.Assets.TempDir,
				},
				Status: pulumi.String("Enabled"),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	res.AssetBucketAccessBlock, err = s3.NewBucketPublicAccessBlock(ctx, fmt.Sprintf("%s-%s-access-block", namespace, cCfg.Assets.Name), &s3.BucketPublicAccessBlockArgs{
		BlockPublicAcls:       pulumi.Bool(true),
		BlockPublicPolicy:     pulumi.Bool(true),
		Bucket:                res.AssetBucket.ID(),
		IgnorePublicAcls:      pulumi.Bool(true),
		RestrictPublicBuckets: pulumi.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	bucketOwnershipName := fmt.Sprintf("%s-%s-ownership", namespace, cCfg.Assets.Name)
	res.AssetBucketOwnership, err = s3.NewBucketOwnershipControls(ctx, bucketOwnershipName, &s3.BucketOwnershipControlsArgs{
		Bucket: res.AssetBucket.ID(),
		Rule: s3.BucketOwnershipControlsRuleArgs{
			ObjectOwnership: pulumi.String("BucketOwnerPreferred"),
		},
	})

	oacName := fmt.Sprintf("%s-%s-oac", namespace, cCfg.Assets.Name)
	res.AssetCDNOAC, err = cloudfront.NewOriginAccessControl(ctx, oacName, &cloudfront.OriginAccessControlArgs{
		Name:                          pulumi.String(oacName),
		OriginAccessControlOriginType: pulumi.String("s3"),
		SigningBehavior:               pulumi.String("always"),
		SigningProtocol:               pulumi.String("sigv4"),
	})
	if err != nil {
		return nil, err
	}

	var (
		cdnName    = fmt.Sprintf("%s-%s-cdn", namespace, cCfg.Assets.Name)
		subDomain  = fmt.Sprintf("%s-%s-%s", ctx.Project(), cCfg.Assets.Name, ctx.Stack())
		domainName = fmt.Sprintf("%s.%s", subDomain, nwCfg.DomainName.Domain)
	)
	res.AssetCDN, err = cloudfront.NewDistribution(ctx, cdnName, &cloudfront.DistributionArgs{
		Aliases: pulumi.StringArray{
			pulumi.String(domainName),
		},
		DefaultCacheBehavior: cloudfront.DistributionDefaultCacheBehaviorArgs{
			AllowedMethods:          pulumi.ToStringArray([]string{"GET", "HEAD"}),
			CachePolicyId:           cCfg.Assets.CachePolicyID,
			CachedMethods:           pulumi.ToStringArray([]string{"GET", "HEAD"}),
			ResponseHeadersPolicyId: cCfg.Assets.ResponseHeadersPolicyID,
			TargetOriginId:          res.AssetBucket.ID(),
			ViewerProtocolPolicy:    pulumi.String("redirect-to-https"),
		},
		Enabled: pulumi.Bool(true),
		Origins: cloudfront.DistributionOriginArray{
			cloudfront.DistributionOriginArgs{
				DomainName:            res.AssetBucket.BucketRegionalDomainName,
				OriginAccessControlId: res.AssetCDNOAC.ID(),
				OriginId:              res.AssetBucket.ID(),
			},
		},
		Restrictions: cloudfront.DistributionRestrictionsArgs{
			GeoRestriction: cloudfront.DistributionRestrictionsGeoRestrictionArgs{
				RestrictionType: pulumi.String("none"),
			},
		},
		PriceClass: cCfg.Assets.PriceClass,
		ViewerCertificate: cloudfront.DistributionViewerCertificateArgs{
			AcmCertificateArn:      cCfg.Assets.ACMCertARN,
			MinimumProtocolVersion: pulumi.String("TLSv1.2_2021"),
			SslSupportMethod:       pulumi.String("sni-only"),
		},
	})
	if err != nil {
		return nil, err
	}

	policyContent := (pulumi.All(res.AssetBucket.Arn, res.AssetCDN.Arn, fnRes.ExecutionRole.Arn).ApplyT(func(args []any) (pulumi.String, error) {
		var (
			bucketArn = args[0].(string)
			cdnArn    = args[1].(string)
			fnExecArn = args[2].(string)
		)

		p := map[string]any{
			"Version": "2012-10-17",
			"Statement": []map[string]any{
				{
					"Principal": map[string]any{
						"Service": "cloudfront.amazonaws.com",
					},
					"Effect":   "Allow",
					"Action":   "s3:GetObject",
					"Resource": fmt.Sprintf("%s/*", bucketArn),
					"Condition": map[string]any{
						"StringEquals": map[string]any{
							"AWS:SourceArn": cdnArn,
						},
					},
				},
				{
					"Principal": map[string]any{
						"AWS": fnExecArn,
					},
					"Effect":   "Allow",
					"Action":   []string{"s3:PutObject", "s3:GetObject"},
					"Resource": filepath.Join(bucketArn, string(cCfg.Assets.TempDir), "/*"),
				},
			},
		}
		b, err := json.Marshal(&p)
		if err != nil {
			return "", err
		}
		return pulumi.String(b), nil
	})).(pulumi.StringOutput)

	bucketPolicyName := fmt.Sprintf("%s-%s-policy", namespace, cCfg.Assets.Name)
	res.AssetBucketPolicy, err = s3.NewBucketPolicy(ctx, bucketPolicyName, &s3.BucketPolicyArgs{
		Bucket: res.AssetBucket.ID(),
		Policy: policyContent,
	})
	if err != nil {
		return nil, err
	}

	route53RecordName := fmt.Sprintf("%s-%s-record", namespace, cCfg.Assets.Name)
	res.AssetRoute53Record, err = route53.NewRecord(ctx, route53RecordName, &route53.RecordArgs{
		Aliases: route53.RecordAliasArray{
			route53.RecordAliasArgs{
				EvaluateTargetHealth: pulumi.Bool(true),
				Name:                 res.AssetCDN.DomainName,
				ZoneId:               res.AssetCDN.HostedZoneId,
			},
		},
		Name:   pulumi.String(domainName),
		Type:   pulumi.String("A"),
		ZoneId: nwCfg.Route53.HostedZoneId,
	})
	if err != nil {
		return nil, err
	}

	return
}
