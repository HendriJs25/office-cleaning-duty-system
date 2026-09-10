package config

const (
	defaultAppName = "cleaning-app"
	defaultAppEnv  = "development"
	defaultAppPort = "8080"
)

type App struct {
	Name string
	Env  string
	Port string
}

func loadApp() App {
	return App{
		Name: getEnv("APP_NAME", defaultAppName),
		Env:  getEnv("APP_ENV", defaultAppEnv),
		Port: getEnv("APP_PORT", defaultAppPort),
	}
}

func (a *App) Validate() error {
	if err := validatePort("APP_PORT", a.Port); err != nil {
		return err
	}

	return nil
}
