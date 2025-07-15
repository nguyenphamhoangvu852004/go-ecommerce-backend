package account

import (
	"go-ecommerce-backend-api/global"
	"go-ecommerce-backend-api/internal/dto"
	"go-ecommerce-backend-api/internal/service"
	"go-ecommerce-backend-api/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Login = &cUserLogin{}

type cUserLogin struct{}

func (c *cUserLogin) Login(ctx *gin.Context) {
	err := service.UserLogin().Login(ctx)
	if err != nil {
		response.ErrorReponse(ctx, response.ErrorSendEmailOTPCode, err.Error())
	}
	response.SuccessReponse(ctx, response.ErrorSendEmailOTPCode, "Success")
}

func (c *cUserLogin) Register(ctx *gin.Context) {
	var params dto.RegisterInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorReponse(ctx, response.ErrorParameterInvalidCode, err.Error())
		return
	}
	codeStatus, err := service.UserLogin().Register(ctx, &params)
	if err != nil {
		global.Logger.Error("Error registration user OTP", zap.Error(err))
		response.ErrorReponse(ctx, codeStatus, err.Error())
		return
	}
	response.SuccessReponse(ctx, codeStatus, nil)
	return
}
