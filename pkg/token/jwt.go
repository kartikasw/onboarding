package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"onboarding/pkg/config"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
)

const (
	AuthorizationHeader = "authorization"
	BearerScheme        = "bearer"
	JWTClaim            = "jwt_claim"
)

type JWT interface {
	CreateAccessToken(usrUUID uuid.UUID) (*JWTToken, error)
	// TODO: IMPLEMENT
	CreateRefreshToken(usrUUID uuid.UUID) (*JWTToken, error)
	VerifyToken(token string, expectation ...Expectation) (*CustomClaims, error)
}

type IJWT struct {
	cfg        config.Token
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewJWT(cfg config.Token) (JWT, error) {
	j := &IJWT{cfg: cfg}
	var err error

	j.publicKey, err = ParseRSAPublicKeyFromPEM(cfg.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse public key: %w", err)
	}

	j.privateKey, err = ParseRSAPrivateKeyFromPEM(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse private key: %w", err)
	}

	return j, nil
}

func (j *IJWT) signer() (jose.Signer, error) {
	opts := (&jose.SignerOptions{}).WithType("JWT")
	return jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.RS256,
			Key:       j.privateKey,
		},
		opts,
	)
}

func (j *IJWT) createJWTToken(claim CustomClaims, duration time.Duration) (*JWTToken, error) {
	now := time.Now()
	exp := now.Add(duration)

	claim.IssuedAt = jwt.NewNumericDate(now)
	claim.NotBefore = jwt.NewNumericDate(now)
	claim.Expiry = jwt.NewNumericDate(exp)

	jti, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("jti uuid error: %w", err)
	}
	claim.TokenID = jti

	signer, err := j.signer()
	if err != nil {
		return nil, fmt.Errorf("Signer error: %w", err)
	}

	token, err := jwt.Signed(signer).Claims(claim).Serialize()
	if err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}

	return &JWTToken{
		SignedToken: token,
		Claims:      claim,
		ExpireAt:    exp,
		Scheme:      BearerScheme,
	}, nil
}

func AccessTokenExpectation() Expectation {
	return Expectation(func(parsed CustomClaims) error {
		if parsed.Scope != Scope(ScopeAccess) {
			return fmt.Errorf("Scope %s to have %s", ScopeAccess, parsed.Scope)
		}
		return nil
	})
}

func (j *IJWT) CreateAccessToken(usrUUID uuid.UUID) (*JWTToken, error) {
	claim := CustomClaims{
		UserID: usrUUID,
		Scope:  Scope(ScopeAccess),
	}

	return j.createJWTToken(claim, j.cfg.AccessTokenDuration)
}

func RefreshTokenExpectation() Expectation {
	return Expectation(func(parsed CustomClaims) error {
		if parsed.Scope != ScopeRefresh {
			return fmt.Errorf("Scope %s to have %s", ScopeRefresh, parsed.Scope)
		}
		return nil
	})
}

func (j *IJWT) CreateRefreshToken(usrUUID uuid.UUID) (*JWTToken, error) {
	claim := CustomClaims{
		UserID: usrUUID,
		Scope:  Scope(ScopeRefresh),
	}

	return j.createJWTToken(claim, j.cfg.RefreshTokenDuration)
}

func (j *IJWT) VerifyToken(token string, expectations ...Expectation) (*CustomClaims, error) {
	parsed, err := jose.ParseSigned(token, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse signed token: %w", err)
	}

	rawPayload, err := parsed.Verify(j.publicKey)
	if err != nil {
		return nil, fmt.Errorf("verification: %w", err)
	}

	var c CustomClaims
	if err := json.Unmarshal(rawPayload, &c); err != nil {
		return nil, fmt.Errorf("Couldn't unmarshal claims: %w", err)
	}

	if err := c.Validate(jwt.Expected{Time: time.Now().UTC()}); err != nil {
		if err.Error() == jwt.ErrExpired.Error() {
			return nil, JWTExpirationError
		}
		return nil, err
	}

	if c.UserID == uuid.Nil {
		return nil, errors.New("Missing user_id")
	}

	if c.UserID == (uuid.UUID{}) {
		return nil, fmt.Errorf("Empty user_id claim")
	}

	if c.TokenID == uuid.Nil {
		return nil, errors.New("Missing token_id")
	}

	if c.TokenID == (uuid.UUID{}) {
		return nil, fmt.Errorf("Empty token_id claim")
	}

	for _, e := range expectations {
		err := e(c)
		if err != nil {
			return nil, fmt.Errorf("Failed expectation: %w", err)
		}
	}

	return &c, nil
}

func ParseRSAPrivateKeyFromPEM(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("Invalid pem")
	}

	if k, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return k, nil
	}

	ifc, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	k, ok := ifc.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("Not RSA private key")
	}

	return k, nil
}

func ParseRSAPublicKeyFromPEM(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("Invalid pem")
	}

	ifc, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		k, ok := ifc.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("Not RSA public key")
		}
		return k, nil
	}

	return x509.ParsePKCS1PublicKey(block.Bytes)
}
