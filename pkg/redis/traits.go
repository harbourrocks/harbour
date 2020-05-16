package redisconfig

// RedisTrait returns redis config
type RedisTrait interface {
	GetRedisConfig() RedisOptions
	SetRedisConfig(RedisOptions)
}

// RequestModel holds the request
type RedisModel struct {
	redisOptions RedisOptions
}

func (m RedisModel) GetRedisConfig() RedisOptions {
	return m.redisOptions
}

func (m *RedisModel) SetRedisConfig(s RedisOptions) {
	m.redisOptions = s
}

func AddRedis(trait RedisTrait, s RedisOptions) {
	trait.SetRedisConfig(s)
}
