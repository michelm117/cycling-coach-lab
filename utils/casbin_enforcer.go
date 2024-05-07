package utils

import (
	"github.com/casbin/casbin"
	"github.com/labstack/echo/v4"
)

type Enforcer struct {
	enforcer *casbin.Enforcer
}
type CasbinEnforcer interface {
	Enforce(role string, resource string, method string) *echo.HTTPError
}

func NewCasbinEnforcer(modelPath string, policyPath string) CasbinEnforcer {
	enforcer := casbin.NewEnforcer(modelPath, policyPath)
	return Enforcer{enforcer: enforcer}
}

func (e Enforcer) Enforce(role string, resource string, method string) *echo.HTTPError {
	if !e.enforcer.Enforce(role, resource, method) {
		return echo.ErrForbidden
	}
	return nil
}
