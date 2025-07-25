package auth

import (
	"errors"
	"fmt"
	"go-ecommerce-backend-api/global"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type PayloadClaims struct {
	jwt.StandardClaims
}

func GenerateToken(payload jwt.Claims) (string, error) {
	fmt.Println("AccessSecret:", global.Config.Jwt.AccessSecret)
	fmt.Println("Kiểu AccessSecret:", reflect.TypeOf(global.Config.Jwt.AccessSecret))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenSigned, err := token.SignedString([]byte(global.Config.Jwt.AccessSecret))

	if global.Config.Jwt.AccessSecret == "" {
		return "", errors.New("access secret key is empty")
	}

	if err != nil {
		return "", err
	}
	return tokenSigned, nil
}

func CreateToken(uuidToken string) (string, error) {
	timeEx := global.Config.Jwt.AccessSecretExpiriedTime
	if timeEx == "" {
		timeEx = "1h"
	}
	fmt.Println("type of timeEx (AccessSecretExpiriedTime): ", reflect.TypeOf(timeEx))
	expiration, err := time.ParseDuration(timeEx)
	if err != nil {
		fmt.Println("Lỗi khi parse:", err)
		return "", err
	}

	timeNow := time.Now()
	exp := timeNow.Add(expiration)

	return GenerateToken(
		&PayloadClaims{
			StandardClaims: jwt.StandardClaims{
				Id:        uuid.New().String(),
				ExpiresAt: exp.Unix(),
				IssuedAt:  timeNow.Unix(),
				Issuer:    "shopdevgo",
				Subject:   uuidToken,
			},
		},
	)
}

func ParseJwtToken(token string) (*PayloadClaims, error) {
	claims := &PayloadClaims{}
	tokenClaims, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.Jwt.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !tokenClaims.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func VerifyToken(token string) (*PayloadClaims, error) {
	claims, err := ParseJwtToken(token)
	if err != nil {
		return nil, err
	}
	if err = claims.Valid(); err != nil {
		return nil, err
	}
	return claims, nil
}
