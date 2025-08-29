package fileBucketClient

import (
	"context"
	"errors"
	"testing"

	logger "github.com/TheGuarantors/tg-logger/pkg"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/theguarantors/tiger/clients/file-bucket/mocks"
)

func initTest(t *testing.T) *logger.Logger {
	t.Setenv("FILE_BUCKET_NAME", "test-bucket")
	return logger.New()

}

func setupTest(t *testing.T) *mocks.S3Client {
	t.Helper()
	return mocks.NewS3Client(t)
}

func TestClient_SaveFile(t *testing.T) {
	logger := initTest(t)

	tests := []struct {
		name          string
		fileName      string
		content       []byte
		mockSetup     func(*mocks.S3Client)
		expectedError error
	}{
		{
			name:     "successful file save",
			fileName: "test.pdf",
			content:  []byte("test content"),
			mockSetup: func(m *mocks.S3Client) {
				m.EXPECT().
					Options().
					Return(s3.Options{
						Region: "us-east-1",
					}).Once()

				m.EXPECT().
					PutObject(mock.Anything, mock.MatchedBy(func(input *s3.PutObjectInput) bool {
						return *input.Key == "test.pdf" && *input.Bucket == "test-bucket"
					}), mock.Anything).
					Return(&s3.PutObjectOutput{}, nil).Once()
			},
			expectedError: nil,
		},
		{
			name:     "error saving file",
			fileName: "error.pdf",
			content:  []byte("test content"),
			mockSetup: func(m *mocks.S3Client) {
				m.EXPECT().
					PutObject(mock.Anything, mock.Anything).
					Return(nil, errors.New("s3 error")).Once()
			},
			expectedError: errors.New("s3 error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s3Client := setupTest(t)
			tt.mockSetup(s3Client)

			client := NewClient(logger, s3Client)

			_, err := client.SaveFile(context.Background(), tt.fileName, tt.content)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}
