package config


import (
	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
	HostDB   	string `env:"DB_HOST" env-default:"localhost"`
    PortDB   	string `env:"DB_PORT" env-default:"5432"`
	UserDB		string `env:"DB_USER" env-default:"admin"`
	PasswordDB 	string `env:"DB_PASSWORD" env-default:"admin"`
	NameDB 		string `env:"DB_NAME" env-default:"banners"`

	ServerPort  string `env:"PORT" env-default:"6050"`

	JWTsecret   string `env:"JWT_SECRET" env-default:"secret"`
}


func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
