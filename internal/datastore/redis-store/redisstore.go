package redisstore

import (
	"context"
	"time"

	"github.com/dashbikash/vidura-sense/internal/common"
	"github.com/redis/go-redis/v9"
)

var (
	log    = common.GetLogger()
	config = common.GetConfig()
	ctx    = context.Background()
)

func getClient() *redis.Client {
	opt, err := redis.ParseURL(config.Data.Redis.RedisUrl)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opt)
}
func GetRobotsTxt(domainName string) {
	rdb := getClient()
	result, err := rdb.Get(ctx, config.Data.Redis.Branches.RobotsTxt.Name+":"+domainName).Result()
	if err != nil {
		panic(err.Error())
	}
	rdb.Close()
	log.Println(result)
}
func SetRobotsTxt(domainName string, robotsTxt string) {
	rdb := getClient()
	result, err := rdb.Set(ctx, config.Data.Redis.Branches.RobotsTxt.Name+":"+domainName, robotsTxt, time.Hour*time.Duration(config.Data.Redis.Branches.RobotsTxt.Ttl)).Result()
	if err != nil {
		panic(err.Error())
	}
	log.Println(result)
}
