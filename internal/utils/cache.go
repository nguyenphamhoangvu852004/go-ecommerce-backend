package utils

import "fmt"

func GetUserKey(hashKey string) string {
	userKey := fmt.Sprintf("u:%s:otp", hashKey)
	fmt.Println(userKey)
	return userKey
}
