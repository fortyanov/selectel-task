package main

import (
	"os"
	"path"

	"selectel-task/env"
)

type config struct {
	PidFile  string
	LogDest  string
	LogTag   string
	LogLevel string
}

func initConfig() (cfg *config, err error) {
	cfg = &config{
		PidFile:  env.GetEnv("PID_FILE", "/var/run/"+path.Base(os.Args[0])+".pid"),
		LogDest:  env.GetEnv("LOG_DESTINATION", "stderr"),
		LogTag:   env.GetEnv("LOG_TAG", path.Base(os.Args[0])),
		LogLevel: env.GetEnv("LOG_LEVEL", "debug"),
	}

	return cfg, err
}
