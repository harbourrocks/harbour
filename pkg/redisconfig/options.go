package redisconfig

// RedisOptions configures all important redis parameter.
type RedisOptions struct {
	Address  string
	Password string
	Database int
}

// NewDefaultRedisOptions returns the default redis options.
func NewDefaultRedisOptions() *RedisOptions {
	s := RedisOptions{
		Address:  "127.0.0.1:6379",
		Password: "",
		Database: 0,
	}

	return &s
}
