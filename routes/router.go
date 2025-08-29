package routes

import (
	"github.com/labstack/echo/v4"
	policyService "github.com/theguarantors/tiger/internal/policy/service"
	userService "github.com/theguarantors/tiger/internal/user/service"
)

type ServerHTTP struct {
	echo *echo.Echo
}

func SetupRoutes(
	userService *userService.UserService,
	policyService *policyService.PolicyService,
) *ServerHTTP {
	e := echo.New()

	routes := getRoutes(
		userService,
		policyService,
	)

	for _, route := range routes {
		switch route.Method {
		case echo.GET:
			e.GET(route.Path, route.Handler)
		case echo.POST:
			e.POST(route.Path, route.Handler)
		case echo.PUT:
			e.PUT(route.Path, route.Handler)
		case echo.DELETE:
			e.DELETE(route.Path, route.Handler)
		case echo.PATCH:
			e.PATCH(route.Path, route.Handler)
		}
	}

	return &ServerHTTP{echo: e}
}

func (s *ServerHTTP) Start() error {
	return s.echo.Start(":8080")
}
