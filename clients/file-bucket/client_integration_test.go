package fileBucketClient

import (
	"context"
	"flag"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go-v2/config"
)

var runIntegrationTests = flag.Bool("integration", false, "Run integration tests")

func setupTestEnvironment(ctx context.Context, t *testing.T) (*s3.Client, *s3.PresignClient) {
	envLoadErr := godotenv.Load("../../.env-integration-test")
	if envLoadErr != nil {
		t.Fatalf("Error loading .env-integration-test file: %v", envLoadErr)
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithSharedConfigProfile("default-dev"),
	)
	if err != nil {
		t.Fatalf("failed to load s3 config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	return s3Client, presignClient
}
