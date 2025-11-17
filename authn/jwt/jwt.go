package authn_jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/configuration"
)

type JwtAuthentication struct {
	*application.EmptyApplication
	config JwtAuthenticationConfig
}

func NewJwtAuthentication() *JwtAuthentication {
	return &JwtAuthentication{
		EmptyApplication: application.NewEmptyApplication("JwtAuthentication"),
	}
}

func (j *JwtAuthentication) Configuration() configuration.Configuration {
	return &j.config
}

// CreateJWT creates a signed JWT with given subject/email/roles.
func (j *JwtAuthentication) CreateJWT(userID, email string, roles []string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.config.Ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.Secret))
}

// ParseJWT parses and validates JWT string.
func (j *JwtAuthentication) ParseJWT(tokenStr string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return j.config.Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func AuthGetterFromAuthorizationHeader(prefix ...[]string) AuthGetter {
	return func(c echo.Context) (bool, string) {
		authHeader := c.Request().Header.Get("Authorization")
		prefixStr := "Bearer "
		if len(prefix) > 0 && len(prefix[0]) > 0 {
			prefixStr = prefix[0][0] + " "
		}
		if authHeader == "" || !strings.HasPrefix(authHeader, prefixStr) {
			return false, ""
		}
		return true, strings.TrimPrefix(authHeader, prefixStr)
	}
}

func AuthGetterFromHeader(headerKey string) AuthGetter {
	return func(c echo.Context) (bool, string) {
		headerValue := c.Request().Header.Get(headerKey)
		if headerValue == "" {
			return false, ""
		}
		return true, headerValue
	}
}

func AuthGetterFromCookie(cookieName string) AuthGetter {
	return func(c echo.Context) (bool, string) {
		cookie, err := c.Cookie(cookieName)
		if err != nil {
			return false, ""
		}
		return true, cookie.Value
	}
}
