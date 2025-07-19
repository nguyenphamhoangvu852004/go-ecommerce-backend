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

// UpdatePasswordRegistration
// @Summary       UpdatePasswordRegistration
// @Description  UpdatePasswordRegistration
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload body dto.UpdateUserPasswordInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router      /auth/update_password_register [post]
func (c *cUserLogin) UpdatePasswordRegistration(ctx *gin.Context) {
	var paramas dto.UpdateUserPasswordInput

	if err := ctx.ShouldBindJSON(&paramas); err != nil {
		response.ErrorReponse(ctx, response.ErrorParameterInvalidCode, err.Error())
		return
	}
	result, err := service.UserLogin().UpdatePasswordRegister(ctx, paramas.UserToken, paramas.UserPassword)
	if err != nil {
		response.ErrorReponse(ctx, response.ErrorInValidOTP, err.Error())
		return
	}
	response.SuccessReponse(ctx, response.VerifyOTPSuccess, result)
}

// Verify OTP Login By User
// @Summary       Verify OTP Login by User
// @Description  Verify OTP Login by User
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload body dto.VerifyInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router      /auth/verify_account [post]
func (c *cUserLogin) VerifyOTP(ctx *gin.Context) {
	var paramas dto.VerifyInput

	if err := ctx.ShouldBindJSON(&paramas); err != nil {
		response.ErrorReponse(ctx, response.ErrorParameterInvalidCode, err.Error())
		return
	}
	result, err := service.UserLogin().VerifyOTP(ctx, &paramas)
	if err != nil {
		response.ErrorReponse(ctx, response.ErrorInValidOTP, err.Error())
		return
	}
	response.SuccessReponse(ctx, response.VerifyOTPSuccess, result)
}

// User Login
// @Summary       User Login
// @Description  User Login
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload body dto.LoginUserInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router      /auth/login [post]
func (c *cUserLogin) Login(ctx *gin.Context) {
	var params dto.LoginUserInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorReponse(ctx, response.ErrorParameterInvalidCode, err.Error())
		return
	}
	codeRs, dataRS, err := service.UserLogin().Login(ctx, &params)
	if err != nil {
		response.ErrorReponse(ctx, codeRs, err.Error())
		return
	}
	response.SuccessReponse(ctx, codeRs, dataRS)
}

// User Registration documentation
// @Summary      User Register
// @Description  when user register, sent otp to email
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload body dto.RegisterInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router      /auth/register [post]
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
	}
	response.SuccessReponse(ctx, codeStatus, nil)
}
