package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"selectel-task/log"
	"selectel-task/pidfile"
	"selectel-task/server"
	"selectel-task/vscale"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "No .env file found")
	}
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	var (
		err error
		cfg *config
	)

	if cfg, err = initConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Init logger error: %s", err)
		return 1
	}

	if err = log.Init(cfg.LogDest, cfg.LogTag, cfg.LogLevel); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Init logger error: %s", err)
		return 1
	}
	defer log.Close()

	if err = pidfile.Write(cfg.PidFile); err != nil {
		log.Error("Create PID file error:", err)
		return 1
	}
	defer pidfile.Unlink(cfg.PidFile)

	if err = server.Init(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Init server error: %s", err)
		return 1
	}

	if err = vscale.Init(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Init vscale error: %s", err)
		return 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server.Run(ctx, cancel)
	defer server.Shutdown()

	signalEvents := make(chan os.Signal, 1)
	signal.Notify(signalEvents, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	select {
	case s := <-signalEvents:
		log.Info(fmt.Sprintf("Caught signal %v: terminating", s))
	case <-ctx.Done():
		return 1
	}

	return 0
}
