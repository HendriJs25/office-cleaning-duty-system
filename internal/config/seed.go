package config

type Seed struct {
	AdminEmail    string
	AdminPassword string
}

func loadSeed() (*Seed, error) {
	email, err := getEnvRequired("SEED_ADMIN_EMAIL")
	if err != nil {
		return nil, err
	}

	password, err := getEnvRequired("SEED_ADMIN_PASSWORD")
	if err != nil {
		return nil, err
	}

	return &Seed{
		AdminEmail:    email,
		AdminPassword: password,
	}, nil
}
