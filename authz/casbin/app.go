package authz_casbin

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/configuration"
)

type CasbinAuthorization struct {
	*application.EmptyApplication
	Enforcer *casbin.Enforcer

	config  CasbinAuthorizationConfig
	Adapter persist.Adapter `inject:",optional"`
	Model   model.Model     `inject:",optional"`
}

func NewCasbinAuthorization() *CasbinAuthorization {
	return &CasbinAuthorization{
		EmptyApplication: application.NewEmptyApplication("CasbinAuthorization"),
	}
}

func (c *CasbinAuthorization) Configuration() configuration.Configuration {
	return &c.config
}

func (c *CasbinAuthorization) Setup(ctx context.Context) {
	var _model model.Model = c.Model
	var err error
	if _model == nil {
		_model, err = model.NewModelFromFile(c.config.ModelConfPath)
		if err != nil {
			panic("failed to load Casbin model from file: " + err.Error())
		}
	}

	var _adapter persist.Adapter = c.Adapter
	if _adapter == nil {
		_adapter = fileadapter.NewAdapter(c.config.PolicyConfPath)
	}

	enforcer, err := casbin.NewEnforcer(_model, _adapter)
	if err != nil {
		panic("failed to create Casbin enforcer: " + err.Error())
	}
	enforcer.SetLogger(newCasbinZapLogger(c.Logger))

	c.Enforcer = enforcer
}

type CasbinAuthorizationDump struct {
	Policy [][]string
	Role   []string
}

func (c *CasbinAuthorization) Dump() (*CasbinAuthorizationDump, error) {
	p, err := c.Enforcer.GetPolicy()
	if err != nil {
		return nil, err
	}
	mockLogger := &mockLogger{}
	rm := c.Enforcer.GetRoleManager()
	rm.SetLogger(mockLogger)

	rm.PrintRoles()

	return &CasbinAuthorizationDump{
		Policy: p,
		Role:   mockLogger.Roles,
	}, nil
}

func (c *CasbinAuthorization) Enforce(user string, obj string, act string) (bool, error) {
	return c.Enforcer.Enforce(user, obj, act)
}

func (c *CasbinAuthorization) AddPolicy(user string, obj string, act string) (bool, error) {
	return c.Enforcer.AddPolicy(user, obj, act)
}

func (c *CasbinAuthorization) GetRole(user string) ([]string, error) {
	return c.Enforcer.GetRolesForUser(user)
}

func (c *CasbinAuthorization) SetRole(user string, role []string) error {
	if _, err := c.Enforcer.DeleteRolesForUser(user); err != nil {
		return err
	}
	_, err := c.Enforcer.AddRolesForUser(user, role)
	return err
}
