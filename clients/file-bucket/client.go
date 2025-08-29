package fileBucketClient

import (
	"bytes"
	"context"
	"fmt"
	"os"

	logger "github.com/TheGuarantors/tg-logger/pkg"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

//go:generate mockery --name=S3Client --output=./mocks --outpkg=mocks --with-expecter
type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	Options() s3.Options
}

type Client struct {
	logger     *logger.Logger
	s3Client   S3Client
	bucketName string
}

func NewClient(logger *logger.Logger, s3Client S3Client) *Client {
	return &Client{
		logger:     logger,
		s3Client:   s3Client,
		bucketName: os.Getenv("FILE_BUCKET_NAME"),
	}
}

func (c *Client) SaveFile(ctx context.Context, fileName string, content []byte) (string, error) {
	if _, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(content),
	}); err != nil {
		c.logger.Error().Err(err).Msg("error occurred while saving file to S3")
		return "", errors.WithStack(err)
	}

	region := c.s3Client.Options().Region

	// Construct the S3 URL
	fileURL := fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", region, c.bucketName, fileName)
	return fileURL, nil
}
