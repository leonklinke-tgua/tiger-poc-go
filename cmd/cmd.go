package cmd

import (
	logger "github.com/TheGuarantors/tg-logger/pkg"
	"github.com/theguarantors/tiger/config"
	policyApp "github.com/theguarantors/tiger/internal/policy/app"
	policyRepository "github.com/theguarantors/tiger/internal/policy/repository"
	policyService "github.com/theguarantors/tiger/internal/policy/service"
	userApp "github.com/theguarantors/tiger/internal/user/app"
	userRepository "github.com/theguarantors/tiger/internal/user/repository"
	userService "github.com/theguarantors/tiger/internal/user/service"
	"github.com/theguarantors/tiger/routes"
)

// Injectors from wire.go:

func InitializeService() (*routes.ServerHTTP, error) {
	logger := logger.New()

	db, err := config.NewDBConnection()
	if err != nil {
		return nil, err
	}
	// other clients would be initialized here

	userRepositoryUserRepository := userRepository.NewUserRepository(db)
	userAppUserApp := userApp.NewUserApp(userRepositoryUserRepository)
	userServiceUserService := userService.NewUserService(userAppUserApp, logger)

	policyRepositoryPolicyRepository := policyRepository.NewPolicyRepository(db)
	policyAppPolicyApp := policyApp.NewPolicyApp(policyRepositoryPolicyRepository, userServiceUserService)
	policyServicePolicyService := policyService.NewPolicyService(policyAppPolicyApp, logger)

	// other services would be initialized here

	serverHTTP := newServer(userServiceUserService, policyServicePolicyService)
	return serverHTTP, nil
}

func newServer(
	userService2 *userService.UserService,
	policyService2 *policyService.PolicyService,
) *routes.ServerHTTP {
	return routes.SetupRoutes(userService2, policyService2)
}
