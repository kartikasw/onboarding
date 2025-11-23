package api

import (
	"errors"
	"net/http"
	"onboarding/api/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	contentType           = "Content-Type"
	multipartFormDataType = "multipart/form-data"
)

func ContentTypeValidation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ct := ctx.Request.Header.Get(contentType)

		if !strings.HasPrefix(ct, multipartFormDataType) {
			err := errors.New("Content-Type is not supported")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(err))
			return
		}

		ctx.Next()
	}
}
