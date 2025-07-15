package repository

import (
	"context"
	"fmt"
	"go-ecommerce-backend-api/global"
	"time"
)

type IUserAuthRepository interface {
	AddOTP(email string, otp int, exp int64) error
}

type userAuthRepository struct {
}

// AddOTP implements IUserAuthRepository.
func (u *userAuthRepository) AddOTP(email string, otp int, exp int64) error {
	// panic("unimplemented")
	key := fmt.Sprintf("usr:%s:otp", email) // usr:email:otp
	err := global.Rdb.SetEx(context.Background(),key, fmt.Sprintf("%d", otp), time.Duration(exp)).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewUserAuthRepository() IUserAuthRepository {
	return &userAuthRepository{}
}
