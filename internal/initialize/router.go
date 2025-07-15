package initialize

import (
	"go-ecommerce-backend-api/global"
	"go-ecommerce-backend-api/internal/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	var r *gin.Engine
	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	// middleware

	// r.Use() //loggin

	// CORS config
	config := cors.Config{
		AllowOrigins:     []string{global.Config.Cors.Url},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		// MaxAge:           12 * time.Hour,
	}
	// r.Use() //corss
	r.Use(cors.New(config))

	// Routers
	mainGroup := r.Group("/api/v1")
	userRouter := router.UserRouter{}

	mainGroup.GET("/checkStatus", func(c *gin.Context) { c.JSON(200, gin.H{"message": "ok lam t da thay doi nhaaa"}) }) // tracking monitor
	{
		// banRouter.InitBanRouter(mainGroup)
		userRouter.InitUserRouter(mainGroup)
	}

	return r
}
