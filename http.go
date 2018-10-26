package main

import (
	stdlog "log"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var logger = stdlog.New(log.StandardLogger().Writer(), "", 0)

func newServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: handler,

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  10 * time.Second,
		ErrorLog:     logger,
	}
}
