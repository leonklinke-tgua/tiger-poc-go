package config

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

func S3BucketConfig(ctx context.Context) (*s3.Client, error) {
	options := []func(*config.LoadOptions) error{}
	if os.Getenv("IS_LOCAL") == "true" {
		options = append(options, config.WithSharedConfigProfile("default-dev"))
	}

	cfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return s3.NewFromConfig(cfg), nil
}
