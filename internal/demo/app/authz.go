package app

import (
	"context"

	authz_casbin "github.com/yoshino-s/go-framework/authz/casbin"
)

type SelfCasbinAuthorization struct {
	*authz_casbin.CasbinAuthorization `nested-inject:""`
}

func NewSelfCasbinAuthorization() *SelfCasbinAuthorization {
	return &SelfCasbinAuthorization{
		CasbinAuthorization: authz_casbin.NewCasbinAuthorization(),
	}
}

func (s *SelfCasbinAuthorization) Setup(ctx context.Context) {
	s.CasbinAuthorization.Setup(ctx)
	s.SetRole("user0", []string{"admin"})
	s.SetRole("user1", []string{"admin"})
	s.AddPolicy("admin", "/obj0/:id", "")
	s.AddPolicy("user0", "/obj1/1", "GET")
	s.Enforcer.SavePolicy()
}
