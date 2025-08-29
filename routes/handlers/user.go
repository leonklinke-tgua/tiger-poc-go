package handlers

import (
	"github.com/labstack/echo/v4"
	userService "github.com/theguarantors/tiger/internal/user/service"
)

func GetUserHandler(c echo.Context, userService *userService.UserService) error {
	resp := userService.GetUser(c.Request().Context(), c.Request())
	defer resp.Body.Close()
	return c.JSON(resp.StatusCode, resp.Body)
}
