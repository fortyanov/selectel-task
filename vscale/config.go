package vscale

import (
	"selectel-task/env"
)

type config struct {
	XToken         string
	RequestTimeout int
}

func initConfig() (cfg *config, err error) {
	cfg = &config{
		XToken:         env.GetEnv("XTOKEN", "NOT EXISTS"),
		RequestTimeout: env.GetEnvAsInt("REQUEST_TIMEOUT", 120),
	}

	return cfg, err
}
