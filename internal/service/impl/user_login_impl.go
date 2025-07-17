package impl

import (
	"context"
	"database/sql"
	"fmt"
	"go-ecommerce-backend-api/global"
	consts "go-ecommerce-backend-api/internal/const"
	"go-ecommerce-backend-api/internal/database"
	"go-ecommerce-backend-api/internal/dto"
	"go-ecommerce-backend-api/internal/utils"
	"go-ecommerce-backend-api/pkg/response"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type sUserLogin struct {
	r *database.Queries
}

// Login implements service.IUserLogin.
func (s *sUserLogin) Login(ctx context.Context) error {
	panic("unimplemented")
}

// Register implements service.IUserLogin.
func (s *sUserLogin) Register(ctx context.Context, in *dto.RegisterInput) (codeResult int, err error) {

	// 1. Hash email
	fmt.Printf("VerifyKey: %s | VerifyType: %d | VerifyPurpose: %s\n", in.VerifyKey, in.VerifyType, in.VerifyPurpose)
	hashKey := utils.GetHash(in.VerifyKey)

	rs, err := s.r.CheckUserBaseExists(ctx, in.VerifyKey)
	// 2. Check user exists in user base
	if err != nil {
		return response.ErrorExistData, err
	}

	if rs > 0 {
		return response.ErrorExistData, fmt.Errorf("user has already registered")
	}

	// 3. Create user base
	usrKey := utils.GetUserKey(hashKey)

	otpFound, err := global.Rdb.Get(ctx, usrKey).Result()
	switch {
	case err == redis.Nil:
		fmt.Println("OTP not found in redis")
	case err != nil:
		fmt.Println("Error getting OTP from redis:", err)
		return response.ErrorInValidOTP, err
	case otpFound != "":
		return response.ErrorOTPNotExists, fmt.Errorf("OTP already exists for this email")
	}

	otpNew := utils.GenerateSixDigitNumber()
	if in.VerifyPurpose == "dev" {
		otpNew = 123456 // For development purposes, use a fixed OTP
	}
	fmt.Println("OTP new is:", otpNew)

	// 5. Save OTP into Redis with expiration
	err = global.Rdb.SetEx(ctx, usrKey, fmt.Sprintf("%d", otpNew), time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err() // 5 minutes expiration
	if err != nil {
		return response.ErrorInValidOTP, err
	}

	// 6. send OTP
	switch in.VerifyType {
	case consts.EMAIL:
		err := utils.SendTextEmailOTP([]string{in.VerifyKey}, utils.User, strconv.Itoa(otpNew))
		if err != nil {
			return response.ErrorSendEmailOTPCode, err
		}
		// 7. save OTP to mysql
		result, err := s.r.CreateVerify(
			ctx,
			database.CreateVerifyParams{
				VerifyOtp:     strconv.Itoa(otpNew),
				VerifyKeyHash: hashKey,
				VerifyType:    sql.NullInt32{1, true},
				VerifyKey:     in.VerifyKey,
			},
		)
		if err != nil {
			return response.ErrorSendEmailOTPCode, err
		}
		lastIdVerifyUser, err := result.LastInsertId()
		if err != nil {
			return response.ErrorSendEmailOTPCode, err
		}
		fmt.Println("lastIdVerifyUser: ", lastIdVerifyUser)
		return response.SuccessSendEmailOTPCode, nil
	case consts.MOBILE:
		return response.ErrorSuccessCode, nil
	}
	return response.ErrorSuccessCode, nil
}

// VerifyOTP implements service.IUserLogin.
func (s *sUserLogin) VerifyOTP(ctx context.Context, in *dto.VerifyInput) (out dto.VerifyOutput, err error) {
	hashKey := utils.GetHash(strings.ToLower(in.VerifyKey))

	otpFound, err := global.Rdb.Get(ctx, utils.GetUserKey(hashKey)).Result()
	if err != nil {
		return out, err
	}

	if in.VerifyCode != otpFound {
		// Nếu nó sai 5 lần trong 1 phút thì sao???

		return out, fmt.Errorf("OTP not match with the storage otp in system")
	}

	infoOtp, err := s.r.GetVerifyOTP(ctx, hashKey)
	if err != nil {
		return out, err
	}

	if err := s.r.UpdateVerifyToVerified(ctx, hashKey); err != nil {
		return out, err
	}

	//output
	out.Token = infoOtp.VerifyKeyHash
	out.Message = "success"

	return out, err
}

// UpdatePasswordRegister implements service.IUserLogin.
func (s *sUserLogin) UpdatePasswordRegister(ctx context.Context, token string, password string) (userId int, err error) {

}

func NewUserLogin(r *database.Queries) *sUserLogin {
	return &sUserLogin{
		r: r,
	}
}
