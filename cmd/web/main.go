package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	logger *slog.Logger
}

var cfg config

func main() {

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var app = &application{
		logger: logger,
	}

	app.logger.Info("Starting server on", "addr", fmt.Sprintf("http://localhost%s", cfg.addr))

	err := http.ListenAndServe(cfg.addr, app.routes())

	app.logger.Error(err.Error())
	os.Exit(1)
}
