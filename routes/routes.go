package routes

import (
	"github.com/labstack/echo/v4"
	policyService "github.com/theguarantors/tiger/internal/policy/service"
	userService "github.com/theguarantors/tiger/internal/user/service"
	"github.com/theguarantors/tiger/routes/handlers"
)

type Route struct {
	Method  string
	Path    string
	Handler func(c echo.Context) error
}

func getRoutes(
	userService *userService.UserService,
	policyService *policyService.PolicyService,
) []Route {
	return []Route{
		{
			Method:  echo.GET,
			Path:    "/",
			Handler: func(c echo.Context) error { return c.JSON(200, "T800 - I'm alive") },
		},
		{
			Method:  echo.GET,
			Path:    "/users/:id",
			Handler: func(c echo.Context) error { return handlers.GetUserHandler(c, userService) },
		},
		{
			Method:  echo.GET,
			Path:    "/policies/:id",
			Handler: func(c echo.Context) error { return handlers.GetPolicyHandler(c, policyService) },
		},
		// Add more routes here as needed
	}
}
