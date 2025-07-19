package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func GetHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes)
}
func GenSalt(length int) (string, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func HashPassword(password string, salt string) string {
	saltedPassword := password + salt

	hashPassword := sha256.Sum256([]byte(saltedPassword))
	return hex.EncodeToString(hashPassword[:])
}

func CheckPassword(password, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func MatchPassword(storeHash, password, salt string) bool {
	hashPassword := HashPassword(password, salt)
	return storeHash == hashPassword
}
