package config

type Config struct {
	App App
}

func Load() (*Config, error) {
	return &Config{
		App: loadApp(),
	}, nil
}

func (c *Config) validate() error {
	if err := c.App.Validate(); err != nil {
		return err
	}

	return nil
}
func (c *Config) ServerAddress() string {
	return ":" + c.App.Port
}
