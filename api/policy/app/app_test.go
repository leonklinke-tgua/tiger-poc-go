package policyApp_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	policyApp "github.com/theguarantors/tiger/api/policy/app"
	"github.com/theguarantors/tiger/api/policy/app/mocks"
	"github.com/theguarantors/tiger/api/structs"
	"github.com/theguarantors/tiger/utils"
)

func setupTest(t *testing.T) (*policyApp.PolicyApp, *mocks.PolicyRepository, *mocks.UserService) {
	mockPolicyRepo := mocks.NewPolicyRepository(t)
	mockUserService := mocks.NewUserService(t)
	app := policyApp.NewPolicyApp(mockPolicyRepo, mockUserService)
	return app, mockPolicyRepo, mockUserService
}

func TestPolicyApp_GetPolicy(t *testing.T) {
	app, mockPolicyRepo, mockUserService := setupTest(t)

	testCases := []struct {
		name           string
		policyID       string
		setupMocks     func()
		expectedError  error
		expectedPolicy *structs.Policy
	}{
		{
			name:     "successfully gets policy with user",
			policyID: "test-policy-id",
			setupMocks: func() {
				policy := &structs.Policy{
					ID:                "test-policy-id",
					UserID:            "test-user-id",
					PolicyType:        structs.LG,
					PolicyAmountCents: 100000,
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				}

				user := &structs.User{
					ID: "test-user-id",
				}
				req := &http.Request{URL: &url.URL{Path: "/users?id=" + policy.UserID}}

				mockPolicyRepo.EXPECT().
					Get(mock.Anything, "test-policy-id").
					Return(policy, nil).
					Once()

				mockUserService.EXPECT().
					GetUser(mock.Anything, req).
					Return(utils.ServerResponse(user, nil)).
					Once()
			},
			expectedError: nil,
			expectedPolicy: &structs.Policy{
				ID:                "test-policy-id",
				UserID:            "test-user-id",
				User:              &structs.User{ID: "test-user-id"},
				PolicyType:        structs.LG,
				PolicyAmountCents: 100000,
			},
		},
		{
			name:     "policy not found",
			policyID: "non-existent-id",
			setupMocks: func() {
				mockPolicyRepo.EXPECT().
					Get(mock.Anything, "non-existent-id").
					Return(nil, errors.New("policy not found")).
					Once()
			},
			expectedError:  errors.New("policy not found"),
			expectedPolicy: nil,
		},
		{
			name:     "user not found",
			policyID: "test-policy-id",
			setupMocks: func() {
				policy := &structs.Policy{
					ID:     "test-policy-id",
					UserID: "non-existent-user",
				}

				mockPolicyRepo.EXPECT().
					Get(mock.Anything, "test-policy-id").
					Return(policy, nil).
					Once()

				req := &http.Request{URL: &url.URL{Path: "/users?id=" + policy.UserID}}
				mockUserService.EXPECT().
					GetUser(mock.Anything, req).
					Return(utils.ServerResponse(nil, errors.New("user not found"))).
					Once()
			},
			expectedError:  errors.New("user not found"),
			expectedPolicy: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			// Execute
			policy, err := app.GetPolicy(context.Background(), tc.policyID)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)

				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, policy)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedPolicy.ID, policy.ID)
			assert.Equal(t, tc.expectedPolicy.UserID, policy.UserID)
			assert.Equal(t, tc.expectedPolicy.User, policy.User)
			assert.Equal(t, tc.expectedPolicy.PolicyType, policy.PolicyType)
			assert.Equal(t, tc.expectedPolicy.PolicyAmountCents, policy.PolicyAmountCents)

		})
	}
}
