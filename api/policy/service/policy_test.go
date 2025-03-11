package policyService_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	policyService "github.com/theguarantors/tiger/api/policy/service"
	"github.com/theguarantors/tiger/api/policy/service/mocks"
	structs "github.com/theguarantors/tiger/api/structs"
	"github.com/theguarantors/tiger/utils"
)

func setupTest(t *testing.T) (*policyService.PolicyService, *mocks.PolicyApp) {
	mockPolicyApp := mocks.NewPolicyApp(t)
	service := policyService.NewPolicyService(mockPolicyApp)
	return service, mockPolicyApp
}

func TestGetPolicy(t *testing.T) {
	service, mockPolicyApp := setupTest(t)
	now := time.Now()

	testCases := []struct {
		name           string
		idRequest      string
		setupMock      func()
		expectedError  error
		expectedPolicy *structs.Policy
	}{
		{
			name:      "Success",
			idRequest: "123",
			setupMock: func() {
				expectedPolicy := &structs.Policy{
					ID:     "123",
					UserID: "123",
					User: &structs.User{
						ID:    "123",
						Name:  "Test User",
						Email: "test@example.com",
					},
					PolicyType:        structs.LG,
					PolicyAmountCents: 100000,
					CreatedAt:         now,
					UpdatedAt:         now,
				}
				mockPolicyApp.EXPECT().
					GetPolicy(mock.Anything, "123").
					Return(expectedPolicy, nil).Once()
			},
			expectedPolicy: &structs.Policy{
				ID:     "123",
				UserID: "123",
				User: &structs.User{
					ID:    "123",
					Name:  "Test User",
					Email: "test@example.com",
				},
				PolicyType:        structs.LG,
				PolicyAmountCents: 100000,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		},
		{
			name:           "Missing ID",
			idRequest:      "",
			setupMock:      func() {},
			expectedError:  errors.New("id is required"),
			expectedPolicy: nil,
		},
		{
			name:      "Policy App Error",
			idRequest: "123",
			setupMock: func() {
				mockPolicyApp.EXPECT().
					GetPolicy(mock.Anything, "123").
					Return(nil, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedPolicy: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.setupMock != nil {
				testCase.setupMock()
			}

			// Generate request
			req, err := http.NewRequest("GET", fmt.Sprintf("/policy?id=%s", testCase.idRequest), strings.NewReader(""))
			assert.NoError(t, err)

			response := service.GetPolicy(context.Background(), req)

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), utils.GetErrorFromResponse(response).Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, utils.ServerResponse(testCase.expectedPolicy, nil), response)
		})
	}
}
