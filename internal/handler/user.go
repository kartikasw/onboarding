package handler

import (
	"context"
	"net/http"
	apiHelper "onboarding/api/helper"
	"onboarding/api/request"
	"onboarding/common"
	"onboarding/internal/service"
	"onboarding/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(ctx *gin.Context) {
	apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
		getUser := func(id uuid.UUID) {
			result, err := h.userService.GetUser(c, id)
			if err != nil {
				resChan <- apiHelper.ResponseData{
					StatusCode: http.StatusInternalServerError,
					Error:      err,
				}
				return
			}

			resChan <- apiHelper.ResponseData{
				StatusCode: http.StatusOK,
				Message:    "User retrieved successfully.",
				Data:       result,
			}
		}

		apiHelper.HandleWithClaim(
			ctx,
			func(claim *token.CustomClaims) {
				getUser(claim.UserID)
			},
			func() {
				var req request.GetDataByUUIDRequest
				if err := ctx.ShouldBindUri(&req); err != nil {
					resChan <- apiHelper.ResponseData{
						StatusCode: http.StatusBadRequest,
						Error:      common.ErrorValidation(err),
					}
				}

				uuid, err := uuid.Parse(req.UUID)
				if err != nil {
					resChan <- apiHelper.ResponseData{
						StatusCode: http.StatusBadRequest,
						Error:      err,
					}
				}

				getUser(uuid)
			},
		)

	})
}
