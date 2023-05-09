package main

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

type mocks struct{}

func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return args.Name, args.Inputs, nil
}

func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

func Test_createBucket(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		bucket, err := createBucket(ctx, "test-account-id")
		assert.NoError(t, err)

		pulumi.All(
			bucket.ID(),
			bucket.URN(),
		).ApplyT(func(args []interface{}) error {
			assert.Equal(t, "my-bucket", string(args[0].(pulumi.ID)))
			assert.Equal(t,
				"urn:pulumi:stack::project::aws:s3/bucketV2:BucketV2::my-bucket",
				string(args[1].(pulumi.URN)),
			)
			return nil
		})
		return nil
	}, pulumi.WithMocks("project", "stack", mocks{}))
	assert.NoError(t, err)
}
