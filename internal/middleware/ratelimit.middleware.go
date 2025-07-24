package middleware

import (
	"fmt"
	"go-ecommerce-backend-api/global"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	redisStore "github.com/ulule/limiter/v3/drivers/store/redis"
)

type RateLimiter struct {
	globalLimiter      *limiter.Limiter
	publicLimiter      *limiter.Limiter
	userPrivateLimiter *limiter.Limiter
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		globalLimiter:      reateLimit("5-M"),
		publicLimiter:      reateLimit("3-M"),
		userPrivateLimiter: reateLimit("1-M"),
	}
}

func reateLimit(interval string) *limiter.Limiter {
	store, err := redisStore.NewStoreWithOptions(global.Rdb, limiter.StoreOptions{
		Prefix:          "rate-limiter",
		MaxRetry:        3,
		CleanUpInterval: time.Hour,
	})

	if err != nil {
		panic(err)
	}
	rate, err := limiter.NewRateFromFormatted(interval)
	if err != nil {
		panic(err)
	}

	limiter := limiter.New(store, rate)
	return limiter
}

// GLOBAL API LIMITER
func (rl *RateLimiter) GlobalAPIRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "global"
		log.Println("key:", key)
		ctx, err := rl.globalLimiter.Get(c, key)
		if err != nil {
			fmt.Println("err:::", err)
			return
		}
		log.Println("ctx:", ctx)
		if ctx.Reached {
			log.Println("ctx.Reached:", ctx.Reached)
			c.AbortWithStatusJSON(429, gin.H{"error": "Too Many Requests"})
			return
		}
		c.Next()
	}
}

// PUBLIC API LIMITER
func (rl *RateLimiter) PublicAPIRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {

		urlPath := c.Request.URL.Path
		reateLimitPath := rl.filterLimitURLPath(urlPath)

		if reateLimitPath != nil {
			log.Println("Client IP:::::", c.ClientIP())

			key := fmt.Sprintf("%s-%s", "111-222-333-44", urlPath)
			log.Println("key:", key)
			ctx, err := reateLimitPath.Get(c, key)
			if err != nil {
				fmt.Println("err:::", err)
				return
			}
			log.Println("ctx:", ctx)
			if ctx.Reached {
				log.Println("ctx.Reached:", ctx.Reached)
				c.AbortWithStatusJSON(429, gin.H{"error": "Too Many Requests"})
				return
			}

			c.Next()
		}
	}

}

// USER PRIVATE API LIMITER
func (rl *RateLimiter) UserPrivateAPIRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {

		urlPath := c.Request.URL.Path
		reateLimitPath := rl.filterLimitURLPath(urlPath)

		if reateLimitPath != nil {
			log.Println("Client IP:::::", c.ClientIP())
			userId := 1001
			key := fmt.Sprintf("%d-%s-%s", userId, "111-222-333-44", urlPath)
			log.Println("key:", key)
			ctx, err := reateLimitPath.Get(c, key)
			if err != nil {
				fmt.Println("err:::", err)
				return
			}
			log.Println("ctx:", ctx)
			if ctx.Reached {
				log.Println("ctx.Reached:", ctx.Reached)
				c.AbortWithStatusJSON(429, gin.H{"error": "Too Many Requests"})
				return
			}

			c.Next()
		}
	}
}

func (rl *RateLimiter) filterLimitURLPath(urlPath string) *limiter.Limiter {
	switch urlPath {
	case "/api/v1/auth/login", "/ping/3":
		return rl.publicLimiter
	case "/api/v1/user/info", "/ping/1":
		return rl.userPrivateLimiter
	}
	return rl.globalLimiter
}
