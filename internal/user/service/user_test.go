package userService_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	logger "github.com/TheGuarantors/tg-logger/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/theguarantors/tiger/entities"
	userService "github.com/theguarantors/tiger/internal/user/service"
	"github.com/theguarantors/tiger/internal/user/service/mocks"
	utils "github.com/theguarantors/tiger/utils"
)

func setupTest(t *testing.T, logger *logger.Logger) (*userService.UserService, *mocks.UserApp) {
	mockApp := mocks.NewUserApp(t)
	service := userService.NewUserService(mockApp, logger)
	return service, mockApp
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	logger := logger.New()
	service, mockApp := setupTest(t, logger)
	testCases := []struct {
		name             string
		requestBody      string
		expectedResponse *http.Response
		mocks            func()
	}{
		{
			name: "successful user creation",
			requestBody: `{
				"email": "test@example.com",
				"name": "Test User"
			}`,
			expectedResponse: utils.ServerResponse(ctx, &entities.User{
				ID:    "123",
				Email: "test@example.com",
				Name:  "Test User",
			}, nil, logger),
			mocks: func() {
				mockApp.EXPECT().
					CreateUser(mock.Anything, &entities.User{
						Email: "test@example.com",
						Name:  "Test User",
					}).
					Return(&entities.User{
						ID:    "123",
						Email: "test@example.com",
						Name:  "Test User",
					}, nil).Once()
			},
		},
		{
			name: "missing email",
			requestBody: `{
				"name": "Test User"
			}`,
			expectedResponse: utils.ServerResponse(ctx, nil, errors.New("email is required"), logger),
		},
		{
			name: "missing name",
			requestBody: `{
				"email": "test@example.com"
			}`,
			expectedResponse: utils.ServerResponse(ctx, nil, errors.New("name is required"), logger),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Setup mock if needed
			if testCase.mocks != nil {
				testCase.mocks()
			}

			// Generate request
			req, err := http.NewRequest("POST", "/users", strings.NewReader(testCase.requestBody))
			assert.NoError(t, err)

			// Call service
			response := service.CreateUser(context.Background(), req)

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response)
		})
	}
}
