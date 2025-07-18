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
	"go-ecommerce-backend-api/internal/utils/crypto"
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
	hashKey := crypto.GetHash(in.VerifyKey)

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
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))

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
	infoOTP, err := s.r.GetVerifyOTP(ctx, token)
	if err != nil {
		return response.ErrorOTPNotExists, err
	}
	if infoOTP.IsVerified.Int32 == 0 {
		return response.ErrorOTPNotExists, fmt.Errorf("User OTP not verified")
	}
	//update user_base

	userBase := database.AddUserBaseParams{}

	userBase.UserAccount = infoOTP.VerifyKey
	userSalt, err := crypto.GenSalt(16)
	if err != nil {
		return response.ErrorOTPNotExists, err
	}
	userBase.UserSalt = userSalt
	userBase.UserPassword = crypto.HashPassword(password, userSalt)

	//add to table
	newUserBase, err := s.r.AddUserBase(ctx, userBase)
	if err != nil {
		return response.ErrorOTPNotExists, err
	}
	user_id, err := newUserBase.LastInsertId()
	if err != nil {
		return response.ErrorOTPNotExists, err
	}

	//add user_id to user_indo table

	newUserInfo, err := s.r.AddUserHaveUserId(ctx, database.AddUserHaveUserIdParams{
		UserID:               uint64(user_id),
		UserAccount:          infoOTP.VerifyKey,
		UserNickname:         sql.NullString{String: infoOTP.VerifyKey, Valid: true},
		UserAvatar:           sql.NullString{String: "", Valid: true},
		UserMobile:           sql.NullString{String: "", Valid: true},
		UserGender:           sql.NullInt16{Int16: 0, Valid: true},
		UserState:            uint8(1),
		UserBirthday:         sql.NullTime{Time: time.Time{}, Valid: false},
		UserEmail:            sql.NullString{String: infoOTP.VerifyKey, Valid: true},
		UserIsAuthentication: uint8(1),
	})

	if err != nil {
		return response.ErrorOTPNotExists, err
	}
	user_id, err = newUserInfo.LastInsertId()
	if err != nil {
		return response.ErrorOTPNotExists, err
	}
	return int(user_id), nil
}

func NewUserLogin(r *database.Queries) *sUserLogin {
	return &sUserLogin{
		r: r,
	}
}
