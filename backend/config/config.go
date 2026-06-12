package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
}

type ServerConfig struct {
	Port        int    `mapstructure:"port"`
	AllowOrigin string `mapstructure:"allow_origin"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type AuthConfig struct {
	SecretKey string `mapstructure:"secret_key"`
	Expire    string `mapstructure:"expire"`
}

type EncryptionConfig struct {
	SecretKey string `mapstructure:"secret_key"`
}

func (a AuthConfig) ExpireDuration() (time.Duration, error) {
	d, err := time.ParseDuration(a.Expire)
	if err != nil {
		return 0, fmt.Errorf("auth.expire is invalid: %w", err)
	}
	return d, nil
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Server.Port == 0 {
		return fmt.Errorf("server.port is required")
	}
	if c.Server.AllowOrigin == "" {
		return fmt.Errorf("server.allow_origin is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if c.Database.Port == 0 {
		return fmt.Errorf("database.port is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database.dbname is required")
	}
	if c.Database.Username == "" {
		return fmt.Errorf("database.username is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("database.password is required")
	}
	if c.Auth.SecretKey == "" {
		return fmt.Errorf("auth.secret_key is required")
	}
	if c.Auth.Expire == "" {
		return fmt.Errorf("auth.expire is required")
	}
	if _, err := c.Auth.ExpireDuration(); err != nil {
		return err
	}
	if c.Encryption.SecretKey == "" {
		return fmt.Errorf("encryption.secret_key is required")
	}
	return nil
}
