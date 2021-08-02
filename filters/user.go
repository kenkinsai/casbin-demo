package filters

import (
	"casbin-demo/models"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/casbin/casbin/v2"
)

// BasicAuthorizer ...
type BasicAuthorizer struct {
	enforcer *casbin.Enforcer
}

// GetUserRole ...
func (a *BasicAuthorizer) GetUserRole(input *context.BeegoInput) string {
	user, ok := input.Session("user").(*models.User)
	// Determine whether the user information is successfully obtained through Session
	if !ok || user.Role.Name == "" {
		// If unsuccessful, return to anonymous directly
		return models.GetRoleString(models.RoleAnonymous)
	}
	return user.Role.Name
}

// NewAuthorizer ...
func NewAuthorizer(e *casbin.Enforcer) beego.FilterFunc {
	return func(ctx *context.Context) {
		// Store the Enforcer by creating a structure
		a := &BasicAuthorizer{enforcer: e}
		// Get user role
		userRole := a.GetUserRole(ctx.Input)
		// Get the access path
		method := strings.ToLower(ctx.Request.Method)
		// Get access method
		path := strings.ToLower(ctx.Request.URL.Path)
		// Verify-return 401 on failure
		if status, err := a.enforcer.Enforce(userRole, path, method); err != nil || !status {
			ctx.Output.Status = 401
			_ = ctx.Output.JSON(map[string]string{"msg": "Insufficient user rights"}, beego.BConfig.RunMode != "prod", false)
		}
	}
}
