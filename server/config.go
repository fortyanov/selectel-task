package server

import (
	"selectel-task/env"
)

type config struct {
	Host         string
	Port         string
	WriteTimeout int
	ReadTimeout  int
	IdleTimeout  int
}

func initConfig() (cfg *config, err error) {
	cfg = &config{
		Host:         env.GetEnv("SERVER_HOST", "127.0.0.1"),
		Port:         env.GetEnv("SERVER_PORT", "8080"),
		WriteTimeout: env.GetEnvAsInt("SERVER_WRITE_TIMEOUT", 120),
		ReadTimeout:  env.GetEnvAsInt("SERVER_READ_TIMEOUT", 120),
		IdleTimeout:  env.GetEnvAsInt("SERVER_IDLE_TIMEOUT", 120),
	}

	return cfg, err
}
