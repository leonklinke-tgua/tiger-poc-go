package policyService

import (
	"context"
	"errors"
	"net/http"

	structs "github.com/theguarantors/tiger/api/structs"
	"github.com/theguarantors/tiger/utils"
)

//go:generate mockery --name PolicyApp --with-expecter --srcpkg github.com/theguarantors/tiger/api/policy/service
type PolicyApp interface {
	GetPolicy(ctx context.Context, id string) (*structs.Policy, error)
}

type PolicyService struct {
	userApp PolicyApp
}

func NewPolicyService(userApp PolicyApp) *PolicyService {
	return &PolicyService{
		userApp: userApp,
	}
}

func (s *PolicyService) GetPolicy(ctx context.Context, request *http.Request) *http.Response {
	id := utils.GetPathParam(request, "id")
	if id == "" {
		return utils.ServerResponse(nil, errors.New("id is required"))
	}

	return utils.ServerResponse(s.userApp.GetPolicy(ctx, id))
}
