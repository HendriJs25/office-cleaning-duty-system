package config

type Redis struct {
	RedisURL string
}

func loadRedis() (*Redis, error) {
	redisURL, err := getEnvRequired("REDIS_URL")
	if err != nil {
		return nil, err
	}

	return &Redis{
		RedisURL: redisURL,
	}, nil
}
