package router

import (
	"go-ecommerce-backend-api/internal/controller/account"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (userRouter *UserRouter) InitUserRouter(router *gin.RouterGroup) {
	//public router
	// userController, _ := wire.InitUserModule()
	userPublicRouter := router.Group("/auth")
	{
		// userPublicRouter.POST("/register", userController.Register)
		userPublicRouter.POST("/register", account.Login.Register)
		userPublicRouter.POST("/verify_account", account.Login.VerifyOTP)
		userPublicRouter.POST("/update_password_register", account.Login.UpdatePasswordRegistration)
		userPublicRouter.POST("/login", account.Login.Login)

		// userPublicRouter.POST("/verifyOtp", authController.VerifyOTP)
		// userPublicRouter.POST("/register", authController.Register)
		// userPublicRouter.POST("/login", authController.Login)
		// userPublicRouter.PUT("/resetPassword", authController.ResetPassword)
	}

	// //private router
	// userPrivateRouter := router.Group("/user")
	// // userPrivateRouter.Use(middleware.Limiter())
	// // userPrivateRouter.Use(middleware.userMiddleware())
	// // userPrivateRouter.Use(middleware.PermissionMiddleware())
	// {
	// 	userPrivateRouter.GET("/getInfo/:id")
	// }

}
