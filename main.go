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
		if _, err := createBucket(ctx, accountID); err != nil {
			return err
		}
		return nil
	})
}

func createBucket(ctx *pulumi.Context, accountID string) (*s3.BucketV2, error) {
	// This configures the AWS provider to use the Cloudflare R2 endpoint.
	// See: https://developers.cloudflare.com/r2/examples/terraform/
	p, err := aws.NewProvider(ctx, "aws.cloudflare_r2", &aws.ProviderArgs{
		Profile:                   pulumi.String("pulumir2"), // your profile name
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
		return nil, err
	}
	// Note: Use NewBucketV2 instead of NewBucket to create a bucket.
	// because NewBucket return error:
	// 	`S3 Bucket acceleration configuration: NotImplemented: GetBucketAccelerateConfiguration not implemented`
	// See: https://github.com/pulumi/pulumi-aws/pull/1859
	return s3.NewBucketV2(ctx, "my-bucket", nil, pulumi.Provider(p))
}
