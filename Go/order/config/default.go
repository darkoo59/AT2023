package config

type Config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

func LoadConfig() (config Config, err error) {
	config = Config{DatabaseURL: "postgresql://user:nURyYtI-sLFoL9MS1nREIA@eighth-guppy-8754.8nj.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"}
	return config, err
}
