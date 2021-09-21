package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/kataras/iris/v12/context"
	"github.com/ujum/dictap/internal/config"
	"net/http"
)

const (
	modelPath     = "/rbac/model.conf"
	policyPath    = "/rbac/policy.csv"
	anonymousRole = "anonymous"
)

type RBACService interface {
	Handler() context.Handler
}

type RBACServiceCached struct {
	Enforcer *casbin.CachedEnforcer
}

func newRBACService(cfg *config.Config) (*RBACServiceCached, error) {
	authEnforcer, err := casbin.NewCachedEnforcer(cfg.ConfigDir+modelPath, cfg.ConfigDir+policyPath)
	if err != nil {
		return nil, err
	}
	return &RBACServiceCached{
		Enforcer: authEnforcer,
	}, nil
}

func (rbac *RBACServiceCached) Handler() context.Handler {
	return func(ctx *context.Context) {
		user := ctx.User()
		var roles []string
		if user != nil {
			var err error
			roles, err = user.GetRoles()
			if err != nil {
				ctx.StopWithJSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
				return
			}
		}
		if !grantAccess(ctx, roles, rbac.Enforcer) {
			return
		}
		ctx.Next()
	}
}

func grantAccess(ctx *context.Context, roles []string, enf *casbin.CachedEnforcer) bool {
	if len(roles) == 0 {
		roles = append(roles, anonymousRole)
	}
	for _, role := range roles {
		res, err := enf.Enforce(role, ctx.Request().URL.Path, ctx.Request().Method)
		if err != nil {
			ctx.StopWithJSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
			return false
		}
		if !res {
			ctx.StopWithStatus(http.StatusForbidden)
			return false
		}
	}
	return true
}
