package middleware

import (
	"context"
	"go-ecommerce-backend-api/internal/utils/auth"
	"log"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the request path

		uri := c.Request.URL.Path
		log.Println("uri:" + uri)

		jwtToken, err := auth.ExtractKeyToken(c)
		if err != nil {
			log.Println("err:", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		// validation token
		claims, err := auth.VerifyToken(jwtToken)
		if err != nil {
			log.Println("err:", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		// update claims to context
		log.Println("claims::::Subject:::", claims.Subject)
		ctx := context.WithValue(c.Request.Context(), "subjectUUID", claims.Subject)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
