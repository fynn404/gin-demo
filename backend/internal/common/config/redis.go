package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client
var Ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis() {
	redisCfg := GlobalConfig.Redis

	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.DB,

		// 连接池设置
		PoolSize:     redisCfg.PoolSize,     // 连接池最大连接数
		MinIdleConns: redisCfg.PoolSize / 4, // 最小空闲连接数，建议是PoolSize的1/4

		// 连接超时设置
		DialTimeout:  5 * time.Second, // 建立连接超时时间
		ReadTimeout:  3 * time.Second, // 读取超时时间
		WriteTimeout: 3 * time.Second, // 写入超时时间
		PoolTimeout:  4 * time.Second, // 当连接池满了，等待连接超时时间

		// 空闲连接检查设置
		IdleCheckFrequency: 60 * time.Second, // 空闲连接检查的频率
		IdleTimeout:        5 * time.Minute,  // 空闲超时时间，超过该时间的空闲连接会被关闭
		MaxConnAge:         30 * time.Minute, // 连接存活时间，超过该时间的连接会被关闭重连
	})

	// 测试连接
	if err := Redis.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Printf("Redis connected successfully. Pool size: %d", redisCfg.PoolSize)
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if Redis != nil {
		if err := Redis.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}
}
