package redisstore

import (
	"context"
	"time"

	"github.com/dashbikash/vidura-sense/internal/system"
	"github.com/redis/go-redis/v9"
)

var (
	log    = system.GetLogger()
	config = system.GetConfig()
	ctx    = context.Background()
)

type RedisStore struct {
	redisClient *redis.Client
}

func DefaultClient() *RedisStore {

	opt, err := redis.ParseURL(config.Data.Redis.RedisUrl)
	if err != nil {
		panic(err)
	}
	return &RedisStore{redisClient: redis.NewClient(opt)}
}

func (rds *RedisStore) GetString(key string, defaultVal string) string {

	result, err := rds.redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Error(err.Error())
		return defaultVal
	}
	rds.redisClient.Close()
	return result
}
func (rds *RedisStore) SetString(key string, val string, ttl time.Duration) bool {
	log.Debug("Adding key to" + key)
	result, err := rds.redisClient.Set(ctx, key, val, time.Duration(ttl)).Result()
	if err != nil {
		log.Error(err.Error())
	}
	rds.redisClient.Close()
	return result == "OK"
}
