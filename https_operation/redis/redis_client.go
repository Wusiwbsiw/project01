package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// AppRedisClient 结构体，用于持有 Redis 客户端
type AppRedisClient struct {
	Client *redis.Client
}

// NewRedisClient 初始化一个新的 Redis 客户端连接
func NewRedisClient(addr, password string, db int) (*AppRedisClient, error) {
	// 创建 Redis 客户端选项
	opts := &redis.Options{
		Addr:     addr,     // Redis 地址，例如 "localhost:6379"
		Password: password, // Redis 密码，如果没有则留空 ""
		DB:       db,       // 使用哪个数据库，0 是默认的
	}

	// 创建客户端
	client := redis.NewClient(opts)

	// 使用 context 来 Ping Redis 服务器，检查连接是否正常
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("无法连接到 Redis: %w", err)
	}

	fmt.Println("Redis 连接成功！")
	return &AppRedisClient{Client: client}, nil
}

// Close 关闭 Redis 连接
func (c *AppRedisClient) Close() error {
	return c.Client.Close()
}
