package config

import (
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "configs/config.yml"

type Config struct {
	ListenConfig     `yaml:"listen"`
	MongodbConfig    `yaml:"mongodb"`
	PostgresdbConfig `yaml:"postgresdb"`
	AuthConfig       `yaml:"auth"`
	Oauth2Config     `yaml:"oauth2"`
}

type ListenConfig struct {
	Port string `yaml:"port"`
}

type MongodbConfig struct {
	URI      string `yaml:"uri"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	AuthDb   string `yaml:"auth_db"`
}

type PostgresdbConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type AuthConfig struct {
	JWT          JWTConfig `yaml:"jwt"`
	PasswordSalt string    `yaml:"password_salt"`
}

type JWTConfig struct {
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	SecretKey       string        `yaml:"secret_key"`
}

type Oauth2Config struct {
	RedirectURL  string   `yaml:"redirect_url"`
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	Scopes       []string `yaml:"scopes"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {

		instance = &Config{}
		if err := cleanenv.ReadConfig(configPath, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			panic(help)
		}
	})
	return instance
}
