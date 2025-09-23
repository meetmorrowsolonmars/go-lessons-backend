package app

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      ServerConfig   `yaml:"server"`
	DebugServer ServerConfig   `yaml:"debug_server"`
	JWT         JWTConfig      `yaml:"jwt"`
	Postgres    PostgresConfig `yaml:"postgres"`
}

func (c *Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config is invalid: %w", err)
	}

	if err := c.DebugServer.Validate(); err != nil {
		return fmt.Errorf("debug server config is invalid: %w", err)
	}

	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("jwt config is invalid: %w", err)
	}

	if err := c.Postgres.Validate(); err != nil {
		return fmt.Errorf("postgres config is invalid: %w", err)
	}

	return nil
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

func (c *ServerConfig) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("address is required")
	}

	return nil
}

type JWTConfig struct {
	SecretKey           string        `yaml:"secret_key"`
	Issuer              string        `yaml:"issuer"`
	AccessTokenDuration time.Duration `yaml:"access_token_duration"`
}

func (c *JWTConfig) Validate() error {
	if c.SecretKey == "" {
		return fmt.Errorf("secret_key is required")
	}

	if c.Issuer == "" {
		return fmt.Errorf("issuer is required")
	}

	if c.AccessTokenDuration == 0 {
		return fmt.Errorf("access_token_duration is required")
	}

	return nil
}

type PostgresConfig struct {
	ConnectionString string `yaml:"connection_string"`
}

func (c *PostgresConfig) Validate() error {
	if c.ConnectionString == "" {
		return fmt.Errorf("connection_string is required")
	}

	return nil
}

func ReadConfig(path string) (Config, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return Config{}, fmt.Errorf("open config file: %w", err)
	}

	defer file.Close()

	var config Config

	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	if err = config.Validate(); err != nil {
		return Config{}, fmt.Errorf("validate config file: %w", err)
	}

	return config, nil
}
