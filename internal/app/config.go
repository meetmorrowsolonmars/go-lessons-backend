package app

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Server      ServerConfig
	DebugServer ServerConfig
	JWT         JWTConfig
}

type ServerConfig struct {
	Address string
}

type JWTConfig struct {
	SecretKey           []byte
	Issuer              string
	AccessTokenDuration time.Duration
}

func ReadConfig() (Config, error) {
	const (
		envVarServerAddress          = "SERVER_ADDRESS"
		envVarDebugServerAddress     = "DEBUG_SERVER_ADDRESS"
		envVarJWTSecretKey           = "JWT_SECRET_KEY"
		envVarJWTIssuer              = "JWT_ISSUER"
		envVarJWTAccessTokenDuration = "JWT_ACCESS_TOKEN_DURATION"
	)

	serverAddress := os.Getenv(envVarServerAddress)
	if serverAddress == "" {
		return Config{}, fmt.Errorf("config variable %s is empty", envVarServerAddress)
	}

	debugServerAddress := os.Getenv(envVarDebugServerAddress)
	if debugServerAddress == "" {
		return Config{}, fmt.Errorf("config variable %s is empty", envVarDebugServerAddress)
	}

	jwtSecretKey := os.Getenv(envVarJWTSecretKey)
	if jwtSecretKey == "" {
		return Config{}, fmt.Errorf("config variable %s is empty", envVarJWTSecretKey)
	}

	jwtIssuer := os.Getenv(envVarJWTIssuer)
	if jwtIssuer == "" {
		return Config{}, fmt.Errorf("config variable %s is empty", envVarJWTIssuer)
	}

	jwtAccessTokenDurationString := os.Getenv(envVarJWTAccessTokenDuration)
	if jwtAccessTokenDurationString == "" {
		return Config{}, fmt.Errorf("config variable %s is empty", envVarJWTAccessTokenDuration)
	}

	jwtAccessTokenDuration, err := time.ParseDuration(jwtAccessTokenDurationString)
	if err != nil {
		return Config{}, fmt.Errorf("config variable %s is invalid", envVarJWTAccessTokenDuration)
	}

	return Config{
		Server: ServerConfig{
			Address: serverAddress,
		},
		DebugServer: ServerConfig{
			Address: debugServerAddress,
		},
		JWT: JWTConfig{
			SecretKey:           []byte(jwtSecretKey),
			Issuer:              jwtIssuer,
			AccessTokenDuration: jwtAccessTokenDuration,
		},
	}, nil
}
