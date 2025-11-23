package apiHelper

import (
	"onboarding/pkg/token"

	"github.com/gin-gonic/gin"
)

func HandleWithClaim(ctx *gin.Context, action func(*token.CustomClaims), fallback func()) {
	claim, exists := ctx.Get(token.JWTClaim)

	if !exists {
		fallback()
		return
	}

	if claim, ok := claim.(*token.CustomClaims); ok {
		action(claim)
		return
	}

	fallback()
}
