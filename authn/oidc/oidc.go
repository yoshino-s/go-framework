package authn_oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"slices"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-errors/errors"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/configuration"
	"golang.org/x/oauth2"
)

// OidcAuthentication encapsulates OIDC provider interaction.
type OidcAuthentication struct {
	*application.EmptyApplication
	config      OidcAuthenticationConfig
	Provider    *oidc.Provider
	Verifier    *oidc.IDTokenVerifier
	OAuthConfig *oauth2.Config
}

// NewOIDCAuth initializes OIDC provider.
func NewOidcAuthentication() *OidcAuthentication {
	return &OidcAuthentication{
		EmptyApplication: application.NewEmptyApplication("OidcAuthentication"),
	}
}

func (o *OidcAuthentication) Configuration() configuration.Configuration {
	return &o.config
}

func (o *OidcAuthentication) Initialize(ctx context.Context) {
	provider, err := oidc.NewProvider(ctx, o.config.IssuerURL)
	if err != nil {
		panic(errors.Errorf("failed to create OIDC provider: %w", err))
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: o.config.ClientID})
	oc := &oauth2.Config{
		ClientID:     o.config.ClientID,
		ClientSecret: o.config.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  o.config.RedirectURL,
		Scopes:       o.config.Scopes,
	}
	o.Provider = provider
	o.Verifier = verifier
	o.OAuthConfig = oc
}

// authCodeURL builds authorization URL (real mode).
func (o *OidcAuthentication) authCodeURL(state string) string {
	return o.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// exchangeAndVerify exchanges code for token & verifies ID token.
func (o *OidcAuthentication) exchangeAndVerify(ctx context.Context, code string) (*oidc.IDToken, map[string]any, error) {
	tok, err := o.OAuthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, nil, fmt.Errorf("exchange: %w", err)
	}
	rawIDToken, ok := tok.Extra("id_token").(string)
	if !ok {
		return nil, nil, fmt.Errorf("id_token missing in token response")
	}
	idTok, err := o.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, nil, fmt.Errorf("verify id token: %w", err)
	}
	claims := map[string]any{}
	if err := idTok.Claims(&claims); err != nil {
		return nil, nil, fmt.Errorf("parse claims: %w", err)
	}
	// Provide fallback roles if provider doesn't include them.
	if _, ok := claims["roles"]; !ok {
		claims["roles"] = []string{"user"}
	}
	return idTok, claims, nil
}

// tokenUser extracts basic user info from claims.
func (o *OidcAuthentication) tokenUser(claims map[string]any) (userID, email string, roles []string) {
	if v, ok := claims["sub"].(string); ok {
		userID = v
	}
	if v, ok := claims["email"].(string); ok {
		email = v
	}
	if rs, ok := claims["roles"].([]any); ok {
		for _, r := range rs {
			if s, ok := r.(string); ok {
				roles = append(roles, s)
			}
		}
	}
	if rs2, ok := claims["roles"].([]string); ok {
		roles = rs2
	}
	if o.config.GroupRoleMap != nil {
		if gs, ok := claims["groups"].([]any); ok {
			for _, g := range gs {
				if s, ok := g.(string); ok {
					s = strings.ToLower(s)
					if r, found := o.config.GroupRoleMap[s]; found {
						roles = append(roles, r)
					}
				}
			}
		}
	}
	if o.config.DefaultRole != "" && !slices.Contains(roles, o.config.DefaultRole) {
		roles = append(roles, o.config.DefaultRole)
	}
	return
}

// randomState returns a cryptographically random string for OAuth state.
func randomState(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
