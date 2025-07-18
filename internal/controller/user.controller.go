package controller

import (
	"fmt"
	"go-ecommerce-backend-api/internal/service"
	"go-ecommerce-backend-api/internal/vo"
	"go-ecommerce-backend-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.IUserService
}

func (uc *UserController) Register(c *gin.Context) {
	var params = vo.UserRegistrationRequest{}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ErrorReponse(c, response.ErrorParameterInvalidCode, err.Error())
		return
	}
	fmt.Println("Registering user with email:", params.Email)
	result := uc.userService.Register(params.Email, params.Purpose)
	response.SuccessReponse(c, response.RegisterSuccessCode, result)
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}
