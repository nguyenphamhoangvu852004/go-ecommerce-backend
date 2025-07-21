package context

import (
	"context"
	"errors"
	"go-ecommerce-backend-api/internal/utils/cache"
)

type InfoUserUUID struct {
	UserID      uint64
	UserAccount string
}

func GetSubUUID(ctx context.Context) (string, error) {
	sUUID, ok := ctx.Value("subjectUUID").(string)
	if !ok {
		return "", errors.New("subject not found")
	}
	return sUUID, nil
}

func GetUserIDFromUUID(ctx context.Context) (uint64, error) {
	sUUid, err := GetSubUUID(ctx)
	if err != nil {
		return 0, err
	}
	// get cache
	var infoUser InfoUserUUID
	err = cache.GetCache(ctx, sUUid, &infoUser)
	if err != nil {
		return 0, err
	}
	return infoUser.UserID, nil
}
