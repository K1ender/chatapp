package config

import "github.com/ilyakaznacheev/cleanenv"

type EmailConfig struct {
	Host string `env:"EMAIL_HOST" env-required:"true"`
	Port int    `env:"EMAIL_PORT" env-required:"true"`
	User string `env:"EMAIL_USER" env-required:"true"`
	Pass string `env:"EMAIL_PASS" env-required:"true"`
	From string `env:"EMAIL_FROM" env-required:"true"`
}

type Config struct {
	Email EmailConfig
}

func MustInit() Config {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)

	if err != nil {
		panic(err)
	}

	return cfg
}
