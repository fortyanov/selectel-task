package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"selectel-task/log"
)

var (
	cfg *config
	srv *http.Server
	ctx context.Context
)

func Init() (err error) {
	if cfg, err = initConfig(); err != nil {
		return err
	}

	r := mux.NewRouter()
	r.HandleFunc("/scalets", scaletsCreate).Methods("POST")
	r.HandleFunc("/scalets", scaletsDelete).Methods("DELETE")

	srv = &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		WriteTimeout: time.Second * time.Duration(cfg.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(cfg.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(cfg.IdleTimeout),
		Handler:      r,
	}

	return
}

func Run(c context.Context, cancel context.CancelFunc) {
	ctx = c
	go func() {
		defer cancel()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server error: ", err)
		}
	}()
}

func Shutdown() {
	c, cancel := context.WithTimeout(ctx, time.Second * time.Duration(30))
	defer cancel()

	if err := srv.Shutdown(c); err != nil {
		log.Error("Gracefully shutdown error: ", err)
	}
}
