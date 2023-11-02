package common

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

var cacheInstance *redis.Client

// Init 初始化数据库
func Init() error {
	cacheInstance = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // Redis 服务器地址和端口
		Password: "",               // Redis 认证密码
		DB:       0,                // Redis 数据库索引
	})
	pong, err := cacheInstance.Ping().Result()
	if err != nil {
		return err
	} else {
		fmt.Println(pong)
		return nil
	}
}

// GetCache ...
func GetCacheClient() *redis.Client {
	if cacheInstance == nil {
		Init()
	}
	return cacheInstance
}

func GetCache(key, def string) (string, error) {
	// 从缓存中获取键对应的值
	val, err := GetCacheClient().Get("key1").Result()
	if err == redis.Nil { // 检查不存在的情况
		log.Println("Key does not exist in cache")
		return def, err
	} else if err != nil {
		log.Println("Error retrieving value from cache:", err)
		return def, err
	} else {
		log.Println(val)
		return val, nil
	}
}

func SetCache(key, value string, expiration time.Duration) error {
	return GetCacheClient().Set("key1", value, expiration).Err()
}
