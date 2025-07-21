package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractKeyToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}
	return "", nil
}
