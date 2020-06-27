package vscale

import (
	"selectel-task/env"
)

type config struct {
	XToken string
}

func initConfig() (cfg *config, err error) {
	cfg = &config{
		XToken: env.GetEnv("XTOKEN", "NOT EXISTS"),
	}

	return cfg, err
}
