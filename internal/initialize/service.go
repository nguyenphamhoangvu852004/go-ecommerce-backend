package initialize

import (
	"go-ecommerce-backend-api/global"
	"go-ecommerce-backend-api/internal/database"
	"go-ecommerce-backend-api/internal/service"
	"go-ecommerce-backend-api/internal/service/impl"
)

func InitServiceInterface() {
	queries := database.New(global.Mdbc)
	// User Service Interface
	service.InitUserLogin(impl.NewUserLogin(queries))
	// ...
}
