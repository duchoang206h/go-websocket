package internal

import "github.com/redis/go-redis/v9"

type RedisOpt struct {
	Addr     string
	Password string
	DB       int
}

func NewRedis(opt RedisOpt) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})
	return rdb
}
