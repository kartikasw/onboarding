package token

import (
	"time"

	_ "github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
)

type JWTError string

func (e JWTError) Error() string {
	return string(e)
}

const JWTExpirationError = JWTError("JWT token is expired")

type CustomClaims struct {
	TokenID uuid.UUID `json:"token_id"`
	UserID  uuid.UUID `json:"user_id"`
	Scope   Scope     `json:"scope"`
	jwt.Claims
}

type Scope string

const (
	ScopeAccess  = Scope("access")
	ScopeRefresh = Scope("refresh")
)

type Expectation func(parsed CustomClaims) error

type JWTToken struct {
	SignedToken string
	Claims      CustomClaims
	ExpireAt    time.Time
	Scheme      string
}
