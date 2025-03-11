//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
	policyApp "github.com/theguarantors/tiger/api/policy/app"
	policyRepo "github.com/theguarantors/tiger/api/policy/repository"
	policyService "github.com/theguarantors/tiger/api/policy/service"
	userApp "github.com/theguarantors/tiger/api/user/app"
	userRepo "github.com/theguarantors/tiger/api/user/repository"
	userService "github.com/theguarantors/tiger/api/user/service"
	"github.com/theguarantors/tiger/config"
	"github.com/theguarantors/tiger/routes"
)

func InitializeService() (*routes.ServerHTTP, error) {

	wire.Build(
		config.NewDBConnection,

		// User Service
		userRepo.NewUserRepository,
		wire.Bind(new(userApp.UserRepository), new(*userRepo.UserRepository)),
		userApp.NewUserApp,
		wire.Bind(new(userService.UserApp), new(*userApp.UserApp)),
		userService.NewUserService,

		// Policy Service
		policyRepo.NewPolicyRepository,
		wire.Bind(new(policyApp.PolicyRepository), new(*policyRepo.PolicyRepository)),
		wire.Bind(new(policyApp.UserService), new(*userService.UserService)),
		policyApp.NewPolicyApp,
		wire.Bind(new(policyService.PolicyApp), new(*policyApp.PolicyApp)),
		policyService.NewPolicyService,

		newServer,
	)
	return &routes.ServerHTTP{}, nil
}

func newServer(
	userService *userService.UserService,
	policyService *policyService.PolicyService,
) *routes.ServerHTTP {

	server := routes.SetupRoutes(userService, policyService)

	return server
}
