package config

type Config struct {
	App      App
	Database *Database
	Seed     *Seed
	JWT      *JWT
	Redis    *Redis
}

func Load() (*Config, error) {
	database, err := loadDatabase()
	if err != nil {
		return nil, err
	}

	seed, err := loadSeed()
	if err != nil {
		return nil, err
	}

	jwt, err := loadJWT()
	if err != nil {
		return nil, err
	}

	redis, err := loadRedis()
	if err != nil {
		return nil, err
	}

	return &Config{
		App:      loadApp(),
		Database: database,
		Seed:     seed,
		JWT:      jwt,
		Redis:    redis,
	}, nil
}

func (c *Config) validate() error {
	if err := c.App.Validate(); err != nil {
		return err
	}
	if err := c.Database.Validate(); err != nil {
		return err
	}
	if err := c.JWT.Validate(); err != nil {
		return err
	}
	if err := c.Redis.Validate(); err != nil {
		return err
	}
	return nil
}
func (c *Config) ServerAddress() string {
	return ":" + c.App.Port
}
