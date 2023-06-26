package db

import (
	"goTestProject/config"
	"goTestProject/tools"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client
var RedisSessClient *redis.Client

func init() {
	InitPublishRedisClient()
}

func InitPublishRedisClient() (err error) {
	redisOpt := tools.RedisOption{
		Address:  config.Conf.Common.CommonRedis.RedisAddress,
		Password: config.Conf.Common.CommonRedis.RedisPassword,
		Db:       config.Conf.Common.CommonRedis.Db,
	}
	RedisClient = tools.GetRedisInstance(redisOpt)
	if pong, err := RedisClient.Ping().Result(); err != nil {
		logrus.Infof("RedisCli Ping Result pong: %s,  err: %s", pong, err)
	}
	//this can change use another redis save session data
	RedisSessClient = RedisClient
	return err
}
