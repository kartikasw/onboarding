package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	apiHelper "onboarding/api/helper"
	"onboarding/api/request"
	"onboarding/common"
	"onboarding/internal/service"

	"github.com/gin-gonic/gin"
)

type ForgotPasswordHandler struct {
	userService service.UserService
	otpService  service.OtpService
}

func NewForgotPasswordHandler(
	otpService service.OtpService,
	userService service.UserService,
) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		otpService:  otpService,
		userService: userService,
	}
}

func (h *ForgotPasswordHandler) RequestResetPassword(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.ResetPasswordRequest
		if err := ctx.ShouldBind(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		err := h.otpService.SendOtpForgotPassword(c, req.Email)
		if err != nil {
			log.Printf("OTP email send failed for %s: %v", req.Email, err)
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusCreated,
			Message:    "If your email exists, an OTP has been sent.",
		}
	})
}

func (h *ForgotPasswordHandler) ResetPassword(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		var req request.ResetPassword
		if err := ctx.ShouldBind(&req); err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      common.ErrorValidation(err),
			}
		}

		err := h.otpService.VerifyOtpForgotPassword(c, req.Email, req.OTP)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			}
		}

		if req.NewPassword != req.VerifyPassword {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      errors.New("Verify password doesn't match."),
			}
		}

		err = h.userService.ChangeUserPassword(ctx, req.Email, req.NewPassword)
		if err != nil {
			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			}
		}

		resChan <- apiHelper.ResponseData{
			StatusCode: http.StatusCreated,
			Message:    "Password reset successfully.",
		}
	})
}
