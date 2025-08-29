package userApp_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/theguarantors/tiger/internal/entities"
	userApp "github.com/theguarantors/tiger/internal/user/app"
	"github.com/theguarantors/tiger/internal/user/app/mocks"
)

func setupTest(t *testing.T) (*userApp.UserApp, *mocks.UserRepository) {
	mockRepo := mocks.NewUserRepository(t)
	app := userApp.NewUserApp(mockRepo)
	return app, mockRepo
}

func TestCreateUser(t *testing.T) {
	userApp, mockRepo := setupTest(t)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         *entities.User
		mocks         func()
		expectedError error
		expectedUser  *entities.User
	}{
		{
			name: "Success",
			input: &entities.User{
				Name:  "Test User",
				Email: "test@example.com",
			},
			mocks: func() {
				mockRepo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(u *entities.User) bool {
						return u.Email == "test@example.com" && u.Name == "Test User"
					})).
					Return(nil).Once()
			},
			expectedError: nil,
			expectedUser: &entities.User{
				Email: "test@example.com",
				Name:  "Test User",
			},
		},
		{
			name: "Repository Error",
			input: &entities.User{
				Name:  "Test User",
				Email: "test@example.com",
			},
			mocks: func() {
				mockRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(errors.New("db error")).Once()
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			if tc.mocks != nil {
				tc.mocks()
			}

			// Execute
			user, err := userApp.CreateUser(ctx, tc.input)

			// Assert
			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}

			assert.Equal(t, tc.expectedUser.Email, user.Email)
			assert.Equal(t, tc.expectedUser.Name, user.Name)
		})
	}
}

func TestGetUser(t *testing.T) {
	// Setup
	ctx := context.Background()
	userApp, mockRepo := setupTest(t)

	testCases := []struct {
		name          string
		userID        string
		expectedUser  *entities.User
		expectedError error
		mocks         func()
	}{
		{
			name:   "successful get user",
			userID: "123",
			expectedUser: &entities.User{
				ID:    "123",
				Email: "test@example.com",
				Name:  "Test User",
			},
			mocks: func() {
				mockRepo.EXPECT().
					Get(mock.Anything, "123").
					Return(&entities.User{
						ID:    "123",
						Email: "test@example.com",
						Name:  "Test User",
					}, nil).Once()
			},
		},
		{
			name:          "user not found",
			userID:        "456",
			expectedError: errors.New("user not found"),
			mocks: func() {
				mockRepo.EXPECT().
					Get(mock.Anything, "456").
					Return(nil, errors.New("user not found")).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			if tc.mocks != nil {
				tc.mocks()
			}

			// Execute
			user, err := userApp.GetUser(ctx, tc.userID)

			// Assert
			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}

			assert.Equal(t, tc.expectedUser, user)
		})
	}
}
