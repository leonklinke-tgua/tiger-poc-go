package userService_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	userStructs "github.com/theguarantors/tiger/api/structs"
	userService "github.com/theguarantors/tiger/api/user/service"
	"github.com/theguarantors/tiger/api/user/service/mocks"
	utils "github.com/theguarantors/tiger/utils"
)

func setupTest(t *testing.T) (*userService.UserService, *mocks.UserApp) {
	mockApp := mocks.NewUserApp(t)
	service := userService.NewUserService(mockApp)
	return service, mockApp
}

func TestCreateUser(t *testing.T) {
	service, mockApp := setupTest(t)
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
			expectedResponse: utils.ServerResponse(&userStructs.User{
				ID:    "123",
				Email: "test@example.com",
				Name:  "Test User",
			}, nil),
			mocks: func() {
				mockApp.EXPECT().
					CreateUser(mock.Anything, &userStructs.User{
						Email: "test@example.com",
						Name:  "Test User",
					}).
					Return(&userStructs.User{
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
			expectedResponse: utils.ServerResponse(nil, errors.New("email is required")),
		},
		{
			name: "missing name",
			requestBody: `{
				"email": "test@example.com"
			}`,
			expectedResponse: utils.ServerResponse(nil, errors.New("name is required")),
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
