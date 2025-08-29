package userRepository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/theguarantors/tiger/internal/entities"
	userRepository "github.com/theguarantors/tiger/internal/user/repository"
)

func setupTest(t *testing.T) (*userRepository.UserRepository, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := userRepository.NewUserRepository(sqlxDB)

	return repo, mock
}

func TestCreate(t *testing.T) {
	repo, mock := setupTest(t)
	ctx := context.Background()

	testCases := []struct {
		name          string
		input         *entities.User
		mockBehavior  func()
		expectedError error
	}{
		{
			name: "Success",
			input: &entities.User{
				ID:    "123",
				Name:  "Test User",
				Email: "test@example.com",
			},
			mockBehavior: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("123", "Test User", "test@example.com").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name: "Database Error",
			input: &entities.User{
				ID:    "123",
				Name:  "Test User",
				Email: "test@example.com",
			},
			mockBehavior: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("123", "Test User", "test@example.com").
					WillReturnError(sqlmock.ErrCancelled)
			},
			expectedError: sqlmock.ErrCancelled,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			err := repo.Create(ctx, tc.input)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGet(t *testing.T) {
	repo, mock := setupTest(t)
	ctx := context.Background()

	testCases := []struct {
		name          string
		userID        string
		mockBehavior  func()
		expectedUser  *entities.User
		expectedError error
	}{
		{
			name:   "Success",
			userID: "123",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email"}).
					AddRow("123", "Test User", "test@example.com")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("123").
					WillReturnRows(rows)
			},
			expectedUser: &entities.User{
				ID:    "123",
				Name:  "Test User",
				Email: "test@example.com",
			},
			expectedError: nil,
		},
		{
			name:   "Not Found",
			userID: "456",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("456").
					WillReturnError(errors.New("user not found"))
			},
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			user, err := repo.Get(ctx, tc.userID)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUser, user)
		})
	}
}
