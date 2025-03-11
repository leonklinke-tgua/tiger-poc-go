package routes

import (
	"fmt"

	"github.com/labstack/echo/v4"
	policyService "github.com/theguarantors/tiger/api/policy/service"
	userService "github.com/theguarantors/tiger/api/user/service"
)

type ServerHTTP struct {
	echo *echo.Echo
}

func SetupRoutes(
	userService *userService.UserService,
	policyService *policyService.PolicyService,
) *ServerHTTP {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		resp := userService.GetUser(c.Request().Context(), c.Request())
		defer resp.Body.Close()
		return c.JSON(resp.StatusCode, resp.Body)
	})

	e.GET("/policies/:id", func(c echo.Context) error {
		fmt.Println("GetPolicy")
		resp := policyService.GetPolicy(c.Request().Context(), c.Request())
		defer resp.Body.Close()
		return c.JSON(resp.StatusCode, resp.Body)
	})

	return &ServerHTTP{echo: e}
}

func (s *ServerHTTP) Start() error {
	return s.echo.Start(":8080")
}
