package repository

import (
	"go-ecommerce-backend-api/global"
	"go-ecommerce-backend-api/internal/database"
)

type IUserRepository interface {
	FindById(id uint)
	FindByEmail(email string) bool
}

type userRepository struct {
	sqlc *database.Queries
}

// FindByEmail implements IUserRepository.
func (u *userRepository) FindByEmail(email string) bool {
	// rs, err := u.sqlc.FindByEmail(ctx, email)
	// if err != nil {
	// 	return false
	// }
	return true
}

// FindById implements IUserRepository.
func (u *userRepository) FindById(id uint) {
	panic("unimplemented")
}

func NewUserRepository() IUserRepository {
	return &userRepository{
		sqlc: database.New(global.Mdbc),
	}
}
