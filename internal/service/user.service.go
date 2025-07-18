package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/global"
	"go-ecommerce-backend-api/internal/repository"
	"go-ecommerce-backend-api/internal/utils"
	"go-ecommerce-backend-api/internal/utils/crypto"
	"go-ecommerce-backend-api/pkg/response"
	"time"

	"github.com/segmentio/kafka-go"
)

type IUserService interface {
	Register(email, purpose string) int
	//...
}

type userService struct {
	userRepo     repository.IUserRepository
	userAuthRepo repository.IUserAuthRepository
}

// Register implements IUserService.
func (u *userService) Register(email string, purpose string) int {

	// 1. Check Exist
	// exist := u.userRepo.FindByEmail(email) // Example ID, replace with actual logic
	// if exist {
	// 	return response.ErrorNotExistCode
	// }
	// 2. Hash Email
	hashedEmail := crypto.GetHash(email)
	fmt.Println("Hashed Email is:", hashedEmail)

	// 3. Create OTP
	OTP := utils.GenerateSixDigitNumber()
	if purpose == "dev" {
		OTP = 123456 // For development purposes, use a fixed OTP
	}
	// fmt.Println("OTP is:", OTP)
	// 4. Save OTP into Redis with expiration
	err := u.userAuthRepo.AddOTP(hashedEmail, OTP, int64(10*time.Minute)) // 5 minutes expiration
	if err != nil {
		return response.ErrorInValidOTP
	}
	// // 5. Send OTP to Email
	// err = utils.SendTextEmailOTP([]string{email}, utils.User, strconv.Itoa(OTP))
	//
	// if err != nil {
	// 	fmt.Println("Error sending email:", err)
	// 	return response.ErrorSendEmailOTPCode
	// }
	// fmt.Println("Email sent successfully to:", email)

	// 5. Send OTP by Golang Email Service
	// if err := utils.SendEmailToGoByAPI(454545, "nphv852004@gmail.com", "OTP"); err != nil {
	// 	return response.ErrorSendEmailOTPCode
	// }

	body := make(map[string]interface{})
	body["otp"] = OTP
	body["email"] = hashedEmail

	bodyRq, err := json.Marshal(body)

	message := kafka.Message{
		Key:   []byte("otp-auth"),
		Value: []byte(bodyRq),
		Time:  time.Now(),
	}

	err = global.KafkaProducer.WriteMessages(context.Background(), message)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return response.ErrorSendEmailOTPCode
	}

	return response.RegisterSuccessCode
}

func NewUserService(userRepo repository.IUserRepository, userAuthRepo repository.IUserAuthRepository) IUserService {
	return &userService{
		userRepo:     userRepo,
		userAuthRepo: userAuthRepo,
	}
}
