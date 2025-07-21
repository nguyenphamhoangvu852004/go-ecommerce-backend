package service

import (
	"context"
	"go-ecommerce-backend-api/internal/dto"
)

type (
	IUserLogin interface {
		Login(ctx context.Context, in *dto.LoginUserInput) (code int, out dto.LoginUserOutput, err error)
		Register(ctx context.Context, in *dto.RegisterInput) (int, error)
		VerifyOTP(ctx context.Context, in *dto.VerifyInput) (dto.VerifyOutput, error)
		UpdatePasswordRegister(ctx context.Context, token string, password string) (userId int, err error)
		//check
		IsTwoFactorEnabled(ctx context.Context, userId int) (code int, rs bool, err error)
		//setup authentication
		SetupTwoFactorAuth(ctx context.Context, in *dto.SetupTwoFactorAuthInput) (code int, err error)
		//VerifyTwoFactorAuth
		VerifyTwoFactorAuth(ctx context.Context, in *dto.TwoFactorVerifyInput) (code int, err error)
	}
	IUserInfo interface {
		GetUserInfoByUserID(ctx context.Context) error
		GetAllUser(ctx context.Context) error
	}
	IUserAdmin interface {
		RemoveUser(ctx context.Context) error
		FindOneUser(ctx context.Context) error
	}
)

var (
	localUserLogin IUserLogin
	localUserInfo  IUserInfo
	localUserAdmin IUserAdmin
)

func UserLogin() IUserLogin {
	if localUserLogin == nil {
		panic("localUserLogin is nil for interface")
	}
	return localUserLogin
}

func InitUserLogin(i IUserLogin) {
	localUserLogin = i
}

func UserAdmin() IUserAdmin {
	if localUserAdmin == nil {
		panic("localUserAdmin is nil for interface")
	}
	return localUserAdmin
}

func InitUserAdmin(i IUserAdmin) {
	localUserAdmin = i
}

func UserInfo() IUserInfo {
	if localUserInfo == nil {
		panic("localUserInfo is nil for interface")
	}
	return localUserInfo
}

func InitUserInfo(i IUserInfo) {
	localUserInfo = i
}
