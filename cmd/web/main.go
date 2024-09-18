package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	logger *slog.Logger
}

type customFileSystem struct {
	fs http.FileSystem
}

func (cfs customFileSystem) Open(path string) (http.File, error) {
	f, err := cfs.fs.Open(path)

	if err != nil {
		return nil, err
	}

	s, err := f.Stat()

	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")

		if _, err := cfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

func main() {
	mux := http.NewServeMux()

	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP Network Address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	app := &application{
		logger: logger,
	}

	filerServer := http.FileServer(customFileSystem{http.Dir("./ui/static/")})

	mux.Handle("GET /static/", http.StripPrefix("/static", filerServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	logger.Info("starting server", "addr", cfg.addr)

	err := http.ListenAndServe(cfg.addr, mux)

	app.logger.Error(err.Error())
	os.Exit(1)
}
