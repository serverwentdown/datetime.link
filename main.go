package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var listen string
var l *zap.Logger

func main() {
	var err error

	flag.StringVar(&listen, "listen", ":8000", "Listen address")
	level := zap.LevelFlag("log", zap.InfoLevel, "zap logging level")
	flag.Parse()

	config := zapConfig()
	config.Level = zap.NewAtomicLevelAt(*level)
	l, err = config.Build()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}
	defer l.Sync()

	app, err := NewDatetime()
	if err != nil {
		l.Panic("failed to create datetime", zap.Error(err))
	}

	server := &http.Server{
		Addr:           listen,
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	l.Info("server is listening", zap.String("addr", listen))
	err = server.ListenAndServe()
	if err != nil {
		l.Panic("server failed", zap.Error(err))
	}
}
