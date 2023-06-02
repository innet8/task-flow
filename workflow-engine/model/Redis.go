package model

import (
	"fmt"
	"time"

	redis "github.com/go-redis/redis"
)

var redisClusterClient *redis.ClusterClient
var redisClient *redis.Client
var clusterIsOpen bool

// RedisOpen 是否连接 redis
var RedisOpen bool

// InitRedis 初始化 redis
func InitRedis() {
	fmt.Println("-------启动 Redis--------")
	if conf.RedisCluster == "true" {
		clusterIsOpen = true
		redisClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{conf.RedisHost + ":" + conf.RedisPort},
			Password: conf.RedisPassword,
		})
		pong, err := redisClusterClient.Ping().Result()
		if err != nil {
			fmt.Printf("------------连接 Redis 集群：%s 失败，原因：%v\n", conf.RedisHost+":"+conf.RedisPort, err)
			return
		}
		RedisOpen = true
		fmt.Printf("---------连接 Redis 集群成功，%v\n", pong)
	} else {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     conf.RedisHost + ":" + conf.RedisPort,
			Password: conf.RedisPassword,
		})
		pong, err := redisClient.Ping().Result()
		if err != nil {
			fmt.Printf("------------连接 Redis：%s 失败，原因：%v\n", conf.RedisHost+":"+conf.RedisPort, err)
			return
		}
		RedisOpen = true
		fmt.Printf("---------连接 Redis 成功，%v\n", pong)
	}
}

// SetRedisVal 将值保存到 Redis
func SetRedisVal(key, value string, expiration time.Duration) error {
	if clusterIsOpen {
		return redisClusterClient.Set(key, value, expiration).Err()
	}
	return redisClient.Set(key, value, expiration).Err()
}

// GetRedisVal 从 Redis 获取值
func GetRedisVal(key string) (string, error) {
	if clusterIsOpen {
		return redisClusterClient.Get(key).Result()
	}
	return redisClient.Get(key).Result()
}
