package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	policyService "github.com/theguarantors/tiger/internal/policy/service"
)

func GetPolicyHandler(c echo.Context, policyService *policyService.PolicyService) error {
	fmt.Println("GetPolicy")
	resp := policyService.GetPolicy(c.Request().Context(), c.Request())
	defer resp.Body.Close()
	return c.JSON(resp.StatusCode, resp.Body)
}
