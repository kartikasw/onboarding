package token

import (
	"encoding/base64"
	"onboarding/common"
	"onboarding/pkg/config"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	_ "github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestExpiredJWTToken(t *testing.T) {
	private, public := common.GenerateRSAKey(t)

	cfg := config.Token{
		AccessTokenDuration: -time.Minute,
		PrivateKey:          private,
		PublicKey:           public,
	}

	jwtImpl, err := NewJWT(cfg)
	require.NoError(t, err)

	usrUUID := uuid.New()
	token, err := jwtImpl.CreateAccessToken(usrUUID)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claim, err := jwtImpl.VerifyToken(token.SignedToken, AccessTokenExpectation())
	require.Error(t, err)
	require.EqualError(t, err, JWTExpirationError.Error())
	require.Nil(t, claim)
}

func TestInvalidJWTToken(t *testing.T) {
	testCases := []struct {
		name        string
		token       func(privateKey string) string
		checkResult func(claim *CustomClaims, err error)
	}{
		{
			name: "NoSigningMethod",
			token: func(privateKey string) string {
				// Create a JWT without a signature (header.alg = none)
				hdr := `{"alg":"none"}`
				payload := `{"user_id":"00000000-0000-0000-0000-000000000000"}`
				return base64.RawURLEncoding.EncodeToString([]byte(hdr)) + "." +
					base64.RawURLEncoding.EncodeToString([]byte(payload)) + "."
			},
			checkResult: func(claim *CustomClaims, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), `unexpected signature algorithm "none"`)
				require.Nil(t, claim)
			},
		},
		{
			name: "InvalidIDFormat",
			token: func(privateKey string) string {
				privKey, _ := ParseRSAPrivateKeyFromPEM(privateKey)

				signer, _ := jose.NewSigner(
					jose.SigningKey{Algorithm: jose.RS256, Key: privKey},
					nil,
				)

				claims := CustomClaims{
					UserID: uuid.New(),
					Scope:  "access",
				}

				tok, err := jwt.Signed(signer).Claims(claims).Serialize()
				require.NoError(t, err)
				return tok
			},
			checkResult: func(claim *CustomClaims, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "Missing token_id")
				require.Nil(t, claim)
			},
		},
		{
			name: "InvalidScope",
			token: func(privatekey string) string {
				privKey, _ := ParseRSAPrivateKeyFromPEM(privatekey)

				signer, _ := jose.NewSigner(
					jose.SigningKey{Algorithm: jose.RS256, Key: privKey},
					nil,
				)

				claims := CustomClaims{
					UserID:  uuid.New(),
					TokenID: uuid.New(),
					Scope:   "refresh",
					Claims: jwt.Claims{
						Expiry: jwt.NewNumericDate(time.Now().Add(time.Hour)),
					},
				}

				tok, err := jwt.Signed(signer).Claims(claims).Serialize()
				require.NoError(t, err)
				return tok
			},
			checkResult: func(claim *CustomClaims, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "Failed expectation")
				require.Nil(t, claim)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			private, public := common.GenerateRSAKey(t)

			cfg := config.Token{
				AccessTokenDuration: -time.Minute,
				PrivateKey:          private,
				PublicKey:           public,
			}

			jwtImpl, err := NewJWT(cfg)
			require.NoError(t, err)

			claim, err := jwtImpl.VerifyToken(tc.token(cfg.PrivateKey), AccessTokenExpectation())
			tc.checkResult(claim, err)
		})
	}
}
