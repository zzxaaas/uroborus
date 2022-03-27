package redis

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

func NewRedisClient() *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host"),
		DB:       viper.GetInt("redis.db"),
		Password: viper.GetString("redis.password"),
	})
	_, err := cli.Ping().Result()
	if err != nil {
		panic(err)
	}
	return cli
}
