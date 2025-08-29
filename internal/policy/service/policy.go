package policyService

import (
	"context"
	"errors"
	"net/http"

	logger "github.com/TheGuarantors/tg-logger/pkg"
	entities "github.com/theguarantors/tiger/entities"
	"github.com/theguarantors/tiger/utils"
)

//go:generate mockery --name PolicyApp --with-expecter --srcpkg github.com/theguarantors/tiger/internal/policy/service
type PolicyApp interface {
	GetPolicy(ctx context.Context, id string) (*entities.Policy, error)
}

type PolicyService struct {
	userApp PolicyApp
	logger  *logger.Logger
}

func NewPolicyService(userApp PolicyApp, logger *logger.Logger) *PolicyService {
	return &PolicyService{
		userApp: userApp,
		logger:  logger,
	}
}

func (s *PolicyService) GetPolicy(ctx context.Context, request *http.Request) *http.Response {
	id := utils.GetPathParam(request, "id")
	if id == "" {
		return utils.ServerResponse(ctx, nil, errors.New("id is required"), s.logger)
	}

	policy, err := s.userApp.GetPolicy(ctx, id)

	return utils.ServerResponse(ctx, policy, err, s.logger)
}
