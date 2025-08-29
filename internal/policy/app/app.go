package policyApp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/theguarantors/tiger/entities"
	"github.com/theguarantors/tiger/utils"
)

//go:generate mockery --name PolicyRepository --with-expecter --srcpkg github.com/theguarantors/tiger/internal/policy/app
type PolicyRepository interface {
	Get(ctx context.Context, id string) (*entities.Policy, error)
}

//go:generate mockery --name UserService --with-expecter --srcpkg github.com/theguarantors/tiger/internal/policy/app
type UserService interface {
	GetUser(ctx context.Context, request *http.Request) *http.Response
}

type PolicyApp struct {
	policyRepo  PolicyRepository
	userService UserService
}

func NewPolicyApp(
	policyRepo PolicyRepository,
	userService UserService,
) *PolicyApp {
	return &PolicyApp{
		policyRepo:  policyRepo,
		userService: userService,
	}
}

func (u *PolicyApp) GetPolicy(ctx context.Context, id string) (*entities.Policy, error) {
	policy, err := u.policyRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := u.fetchUser(ctx, policy.UserID)
	if err != nil {
		return nil, err
	}

	policy.User = user

	return policy, nil
}

// fetchUser fetches a user from the user service
func (u *PolicyApp) fetchUser(ctx context.Context, id string) (*entities.User, error) {
	req := &http.Request{URL: &url.URL{Path: "/users?id=" + id}}
	response := u.userService.GetUser(ctx, req)

	if utils.ResponseFailed(response) {
		return nil, utils.GetErrorFromResponse(response)
	}

	userResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var user entities.User
	err = json.Unmarshal(userResponse, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
