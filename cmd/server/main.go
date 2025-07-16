package main

import (
	"go-ecommerce-backend-api/internal/initialize"

	_ "go-ecommerce-backend-api/cmd/swag/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           API Documentation Go Ecommerce Backend SHOPDEVGO
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  github.com/nguyenphamhoangvu852004/go-ecommerce-backend

// @contact.name   Team Vu
// @contact.url     github.com/nguyenphamhoangvu852004/go-ecommerce-backend
// @contact.email   nguyenphamhoangvu852004@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8001
// @BasePath   /api/v1
// @schema	http
func main() {
	r := initialize.Run()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8001")
}
