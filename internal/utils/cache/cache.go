package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/global"

	"github.com/redis/go-redis/v9"
)

func GetCache(ctx context.Context, key string, obj interface{}) error {
	rs, err := global.Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("redis: not found")
	} else if err != nil {
		return err
	}

	// convert rs to object

	if err := json.Unmarshal([]byte(rs), obj); err != nil {
		return fmt.Errorf("redis: unmarshal error %w", err)
	}

	return nil
}
