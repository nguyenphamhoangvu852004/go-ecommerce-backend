//go:build wireinject

package wire

import (
	"go-ecommerce-backend-api/internal/controller"
	"go-ecommerce-backend-api/internal/repository"
	"go-ecommerce-backend-api/internal/service"

	"github.com/google/wire"
)

func InitUserModule() (*controller.UserController, error) {
	wire.Build(
		controller.NewUserController,
		service.NewUserService,
		repository.NewUserRepository,
		repository.NewUserAuthRepository,
	)
	return new(controller.UserController),nil
}
