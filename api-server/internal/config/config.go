package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `env:"envtype" env-default:"local"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	S3       S3Config       `yaml:"s3"`
	Secrets  SecretConfig
}

type ServerConfig struct {
	Port    int           `yaml:"port" env:"SERVER_PORT"`
	Timeout time.Duration `yaml:"timeout"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASS" env-required:"true"`
	Dbname   string `env:"DB_NAME" env-required:"true"`
}

type S3Config struct {
	Key    string `env:"YAS3_KEY" env-required:"true"`
	Secret string `env:"YAS3_SECRET" env-required:"true"`
	Bucket string `yaml:"bucket"`
}

type SecretConfig struct {
	InviteCode string `env:"INVITE_CODE"`
	JwtCode    string `env:"JWT_CODE"`
	PassSalt   string `env:"PASS_SALT"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	return cfg, err
}

// type Config struct {
// 	Server  ServerConfig
// 	Storage StorageConfig
// }

// type ServerConfig struct {
// 	Port int `yaml:"port"`
// }

// type StorageConfig struct {
// 	MemeStore  PsqlConfig
// 	BoardStore PsqlConfig
// 	UserStore  PsqlConfig
// 	MediaStore PsqlConfig
// }

// type PsqlConfig struct {
// 	Host     string `yaml:"host"`
// 	Port     int    `yaml:"port"`
// 	User     string `env:"DB_USER" env-required:"true"`
// 	Password string `env:"DB_PASS" env-required:"true"`
// 	Dbname   string `env:"DB_NAME" env-required:"true"`
// }
