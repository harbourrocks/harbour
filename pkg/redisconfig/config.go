package redisconfig

import (
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
)

// ParseViperConfig tries to map a viper configuration to LoggingOptions
func ParseViperConfig() *RedisOptions {
	l := NewDefaultRedisOptions()

	if v := viper.GetString("REDIS_ADDRESS"); v != "" {
		l.Address = v
	}

	if v := viper.GetString("REDIS_PASSWORD"); v != "" {
		l.Password = v
	}

	l.Database = viper.GetInt("REDIS_DB")

	return l
}

// OpenClient creates a new redis client.
func OpenClient(o *RedisOptions) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     o.Address,
		Password: o.Password,
		DB:       o.Database,
	})

	return client
}
