package main

import (
	"fmt"
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	accountID := os.Getenv("CF_ACCOUNT_ID")
	pulumi.Run(func(ctx *pulumi.Context) error {
		p, err := aws.NewProvider(ctx, "aws.cloudflare_r2", &aws.ProviderArgs{
			Profile:                   pulumi.String("pulumir2"),
			Region:                    pulumi.String("auto"),
			SkipCredentialsValidation: pulumi.Bool(true),
			SkipRegionValidation:      pulumi.Bool(true),
			SkipRequestingAccountId:   pulumi.Bool(true),
			SkipMetadataApiCheck:      pulumi.Bool(true),
			Endpoints: aws.ProviderEndpointArray{
				aws.ProviderEndpointArgs{
					S3: pulumi.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)),
				},
			},
		})
		if err != nil {
			return err
		}
		// Note: Use NewBucketV2 instead of NewBucket to create a bucket.
		// because NewBucket return error:
		// 	`S3 Bucket acceleration configuration: NotImplemented: GetBucketAccelerateConfiguration not implemented`
		// See: https://github.com/pulumi/pulumi-aws/pull/1859
		_, err = s3.NewBucketV2(ctx, "my-bucket", nil, pulumi.Provider(p))
		if err != nil {
			return err
		}
		return nil
	})
}
