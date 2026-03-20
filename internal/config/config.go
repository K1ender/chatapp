package config

import "github.com/ilyakaznacheev/cleanenv"

type EmailConfig struct {
	Host string `env:"EMAIL_HOST" env-required:"true"`
	Port int    `env:"EMAIL_PORT" env-required:"true"`
	User string `env:"EMAIL_USER" env-required:"true"`
	Pass string `env:"EMAIL_PASS" env-required:"true"`
	From string `env:"EMAIL_FROM" env-required:"true"`
}

type DatabaseConfig struct {
	Host string `env:"DB_HOST" env-required:"true"`
	Port int    `env:"DB_PORT" env-required:"true"`
	User string `env:"DB_USER" env-required:"true"`
	Pass string `env:"DB_PASS" env-required:"true"`
	Name string `env:"DB_NAME" env-required:"true"`
}

type Config struct {
	Email    EmailConfig
	Database DatabaseConfig
	Salt     string `env:"SALT" env-required:"true"`
}

func MustInit() Config {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)

	if err != nil {
		panic(err)
	}

	return cfg
}
