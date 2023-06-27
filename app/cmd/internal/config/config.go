package config

import (
	"auth-service/app/cmd/package/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	App      App      `yaml:"app"`
	Token    Token    `yaml:"token"`
	Postgres Postgres `yaml:"postgresql"`
}

type App struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Post     string `yaml:"post"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Token struct {
	AccessTokenPrivateKeyPath  string `yaml:"access_token_private_key_path"`
	AccessTokenPublicKeyPath   string `yaml:"access_token_public_key_path"`
	RefreshTokenPrivateKeyPath string `yaml:"refresh_token_private_key_path"`
	RefreshTokenPublicKeyPath  string `yaml:"refresh_token_public_key_path"`
	JwtExpiration              int64  `yaml:"jwt_expiration"`
}

var instance *Config
var once sync.Once

func NewConfig(l logger.Logger) *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", &instance); err != nil {
			helpText := "can not read config"
			helpMessage, _ := cleanenv.GetDescription(instance, &helpText)
			l.Error("failed to read config, err: %s", helpMessage)
			log.Println(helpMessage)
			log.Fatalln(err)
		}
	})

	return instance
}
