package apibaseappgateway

import (
	"backend_base_app/shared/log"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (r GatewayApiBaseApp) GenerateID(ctx context.Context) string {
	log.Info(ctx, "called")

	id, err := gonanoid.Generate("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", 4)
	if err != nil {
		return "abcd"
	}

	return id
}

func (r GatewayApiBaseApp) EncryptPassword(ctx context.Context, text string) string {
	log.Info(ctx, "called")

	pass := md5.Sum([]byte(text))

	return hex.EncodeToString(pass[:])
}

func testCache(cacheConnection *cache.Cache) {
	ctx := context.TODO()
	key := "testCacheApimanager"
	obj := "ok"
	if err := cacheConnection.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: obj,
		TTL:   time.Hour,
	}); err != nil {
		fmt.Println("error >>>", err.Error)
	} else {
		fmt.Println("redis connected !")
	}
}
