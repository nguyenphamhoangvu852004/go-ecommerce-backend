package router

import (
	"go-ecommerce-backend-api/internal/controller/account"
	"go-ecommerce-backend-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (userRouter *UserRouter) InitUserRouter(router *gin.RouterGroup) {

	//public router
	userPublicRouter := router.Group("/auth")
	{
		userPublicRouter.POST("/register", account.Login.Register)
		userPublicRouter.POST("/verify_account", account.Login.VerifyOTP)
		userPublicRouter.POST("/update_password_register", account.Login.UpdatePasswordRegistration)
		userPublicRouter.POST("/login", account.Login.Login)
	}

	// //private router
	userPrivateRouter := router.Group("/auth")
	userPrivateRouter.Use(middleware.AuthMiddleware())
	// userPrivateRouter.Use(middleware.Limiter())
	// userPrivateRouter.Use(middleware.userMiddleware())
	// userPrivateRouter.Use(middleware.PermissionMiddleware())
	{
		userPrivateRouter.GET("/getInfo/:id")
		userPrivateRouter.POST("/two_factor/setup", account.TwoFactor.SetupTwoFactor)
		userPrivateRouter.POST("/two_factor/verify", account.TwoFactor.VerifyTwoFactor)
	}

}
