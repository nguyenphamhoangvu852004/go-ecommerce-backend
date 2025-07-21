package account

import (
	"go-ecommerce-backend-api/internal/dto"
	"go-ecommerce-backend-api/internal/service"
	"go-ecommerce-backend-api/internal/utils/context"
	"go-ecommerce-backend-api/pkg/response"

	"github.com/gin-gonic/gin"
)

var TwoFactor = new(sUser2FA)

type sUser2FA struct{}

// Setup two factor  authentication
// @Summary       Setup two factor  authentication
// @Description  Setup two factor  authentication
// @Tags         account 2fa
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Authorization token"
// @Param        payload body dto.SetupTwoFactorAuthInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router      /auth/two_factor/setup [post]
func (s *sUser2FA) SetupTwoFactor(ctx *gin.Context) {
	var params dto.SetupTwoFactorAuthInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorReponse(ctx, response.ErrorParameterInvalidCode, err.Error())
		return
	}
	//get userid from uuid token
	uuid, error := context.GetUserIDFromUUID(ctx.Request.Context())
	if error != nil {
		response.ErrorReponse(ctx, response.ErrorParameterInvalidCode, error.Error())
		return
	}
	params.UserId = uint32(uuid)

	code, err := service.UserLogin().SetupTwoFactorAuth(ctx, &params)
	if err != nil {
		response.ErrorReponse(ctx, code, err.Error())
		return
	}
	response.SuccessReponse(ctx, code, err)
}
