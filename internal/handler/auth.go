package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	apiHelper "onboarding/api/helper"
	"onboarding/api/request"
	"onboarding/api/response"
	"onboarding/common"
	"onboarding/internal/service"
	"onboarding/pkg/token"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.RegisterRequest
		if err := ctx.ShouldBind(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		result, err := h.authService.Register(c, req.Email, req.Password)

		if err != nil {
			err := err
			if common.ErrorCode(err) == common.ErrUniqueViolation {
				err = errors.New("E-mail is already registered.")
			}
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusCreated,
			Message:    "Registration completed successfully.",
			Data:       response.NewUserResponse(result),
		}
	})
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.LoginRequest
		if err := ctx.ShouldBind(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		token, err := h.authService.Login(c, req.Email, req.Password)
		if err != nil {
			var statusCode = http.StatusInternalServerError
			if errors.Is(err, common.ErrRecordNotFound) || common.ErrorCode(err) == fmt.Sprint(common.ErrCredentiials) {
				err = errors.New("E-mail or Password is incorrect")
				statusCode = http.StatusUnauthorized
			}
			resChan <- apiHelper.ResponseData{
				StatusCode: statusCode,
				Error:      err,
			}
		}
		if token == nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusInternalServerError,
				Error:      errors.New("Failed to generate token."),
			}
		}

		setAuthCookies(ctx, token)

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusOK,
			Message:    "Login successful.",
			Data:       response.NewLoginResponse(token.SignedToken),
		}
	})
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		setAuthCookies(ctx, nil)

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusOK,
			Message:    "Logout successful.",
		}
	})
}

func setAuthCookies(ctx *gin.Context, token *token.JWTToken) {
	var (
		value  string
		maxAge int
	)

	if token == nil {
		value = ""
		maxAge = -1
	} else {
		value = token.SignedToken
		maxAge = int(time.Until(token.ExpireAt).Seconds())
		if maxAge < 0 {
			maxAge = -1
		}
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
	})
}
