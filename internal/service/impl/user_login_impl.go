package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/global"
	consts "go-ecommerce-backend-api/internal/const"
	"go-ecommerce-backend-api/internal/database"
	"go-ecommerce-backend-api/internal/dto"
	"go-ecommerce-backend-api/internal/utils"
	"go-ecommerce-backend-api/internal/utils/auth"
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

// VerifyTwoFactorAuth implements service.IUserLogin.
func (s *sUserLogin) VerifyTwoFactorAuth(ctx context.Context, in *dto.TwoFactorVerifyInput) (code int, err error) {

	// check enabled
	isTwoFactorEnabled, err := s.r.IsTwoFactorEnabled(ctx, int32(in.UserId))
	if err != nil {
		return response.ErrorCodeTwoFactorAuthenVerify, err
	}
	if isTwoFactorEnabled > 0 {
		return response.ErrorCodeTwoFactorAuthenVerify, fmt.Errorf("Two factor already enabled")
	}

	keyHash := crypto.GetHash("2fa:" + strconv.Itoa(int(in.UserId)))

	// check otp in redis available
	otpVerifyAuth, err := global.Rdb.Get(ctx, keyHash).Result()
	if err == redis.Nil {
		return response.ErrorCodeTwoFactorAuthenVerify, fmt.Errorf("OTP %s not found in redis", keyHash)
	} else if err != nil {
		return response.ErrorCodeTwoFactorAuthenVerify, err
	}

	// match otp
	if otpVerifyAuth != in.TwoFactorCode {
		fmt.Printf("otpVerifyAuth: %s | in.TwoFactorCode: %s\n", otpVerifyAuth, in.TwoFactorCode)
		return response.ErrorCodeTwoFactorAuthenVerify, fmt.Errorf("OTP not match with the storage otp in system")
	}

	// update otp to verified
	err = s.r.UpdateTwoFactorStatus(ctx, database.UpdateTwoFactorStatusParams{
		UserID:            int32(in.UserId),
		TwoFactorAuthType: database.PreGoAccUserTwoFactor9999TwoFactorAuthTypeEMAIL,
	})

	if err != nil {
		return response.ErrorCodeTwoFactorAuthenVerify, err
	}

	// remove otp from redis
	err = global.Rdb.Del(ctx, keyHash).Err()
	if err != nil {
		return response.ErrorCodeTwoFactorAuthenVerify, err
	}

	return 200, nil
}

// IsTwoFactorEnabled implements service.IUserLogin.
func (s *sUserLogin) IsTwoFactorEnabled(ctx context.Context, userId int) (code int, rs bool, err error) {
	panic("unimplemented")
}

// SetupTwoFactorAuth implements service.IUserLogin.
func (s *sUserLogin) SetupTwoFactorAuth(ctx context.Context, in *dto.SetupTwoFactorAuthInput) (code int, err error) {

	// 1. kiểm tra xem nó có bật tính năng lên chưa -> rồi thì return luôn
	isTrue, err := s.r.IsTwoFactorEnabled(ctx, int32(in.UserId))
	if err != nil {
		return response.ErrorCodeTwoFactorAuthenSetup, err
	}
	if isTrue > 0 {
		return response.ErrorCodeTwoFactorAuthenSetup, fmt.Errorf("Two factor already enabled")
	}

	// 2. Enable nó lên
	error := s.r.EnableTwoFactorTypeEmail(ctx, database.EnableTwoFactorTypeEmailParams{
		UserID:            int32(in.UserId),
		TwoFactorAuthType: database.PreGoAccUserTwoFactor9999TwoFactorAuthType(in.TwoFactorAuthType),
		TwoFactorEmail:    sql.NullString{String: in.TwoFactorEmail, Valid: true},
	})
	if error != nil {
		return response.ErrorCodeTwoFactorAuthenSetup, err
	}

	// 3. Gữi OTP qua in.TwoFactorEmail
	keyHash := crypto.GetHash("2fa:" + strconv.Itoa(int(in.UserId)))
	go global.Rdb.SetEx(ctx, keyHash, "123456", time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()

	return response.SetupTwoFactorAuthCodeSuccess, nil
}

// Login implements service.IUserLogin.
func (s *sUserLogin) Login(ctx context.Context, in *dto.LoginUserInput) (codeResult int, out dto.LoginUserOutput, err error) {
	userBase, err := s.r.GetOneUserInfo(ctx, in.UserAccount)
	if err != nil {
		return response.ErrorAuthFailed, out, err
	}

	// check password
	if !crypto.MatchPassword(userBase.UserPassword, in.UserPassword, userBase.UserSalt) {
		return response.ErrorAuthFailed, out, fmt.Errorf("Password does not match")
	}

	// check two factor authentication
	isTwoFactoerEnabled, err := s.r.IsTwoFactorEnabled(ctx, int32(userBase.UserID))
	if err != nil {
		return response.ErrorAuthFailed, out, err
	}
	if isTwoFactoerEnabled > 0 {
		// send otp to email
		keyUserLoginTwoFactor := crypto.GetHash("2fa:otp" + strconv.Itoa(int(userBase.UserID)))
		err := global.Rdb.SetEx(ctx, keyUserLoginTwoFactor, "111111", time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
		if err != nil {
			return response.ErrorAuthFailed, out, err
		}

		// send otp to email
		// get email 2fa

		infoUserTwoFactor, err := s.r.GetTwoFactorMethodByIDAndType(ctx, database.GetTwoFactorMethodByIDAndTypeParams{
			UserID:            int32(userBase.UserID),
			TwoFactorAuthType: database.PreGoAccUserTwoFactor9999TwoFactorAuthTypeEMAIL,
		})
		if err != nil {
			return response.ErrorAuthFailed, out, err
		}
		fmt.Println("infoUserTwoFactor", infoUserTwoFactor)

		go utils.SendTextEmailOTP([]string{infoUserTwoFactor.TwoFactorEmail.String}, "nguyenphamhoangvu852004@gmail.com", "111111")

		out.Message = "Send OTP to email"

		return response.SuccessSendEmailOTPCode, out, nil

	}

	// update passworrd time
	go s.r.LoginUserBase(ctx, database.LoginUserBaseParams{
		UserLoginIp:  sql.NullString{String: "127.0.0.1", Valid: true},
		UserAccount:  in.UserAccount,
		UserPassword: in.UserPassword,
	})

	// Create UUID user
	subToken := utils.GenerateCliTokenUUID(int(userBase.UserID))
	fmt.Println("SubToken is ", subToken)
	// Get user info table
	userInfo, err := s.r.GetUserByAccount(ctx, uint64(userBase.UserID))
	if err != nil {
		return response.ErrorAuthFailed, out, err
	}

	// convert to json
	userInfoJSON, err := json.Marshal(userInfo)
	if err != nil {
		return response.ErrorAuthFailed, out, err
	}

	// give userrInfoJSON into Redis with SubToken
	if err := global.Rdb.Set(ctx, subToken, userInfoJSON, time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err(); err != nil {
		return response.ErrorAuthFailed, out, err
	}

	// create token
	out.Token, err = auth.CreateToken(subToken)
	out.Message = "Login Sucess"
	if err != nil {
		return response.ErrorAuthFailed, out, err
	}
	return 200, out, nil
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
