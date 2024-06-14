package main

import (
	"flag"
	"fmt"
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

func main() {

	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var app = &application{
		logger: logger,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir(cfg.staticDir)})

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	app.logger.Info("Starting server on", "addr", fmt.Sprintf("http://localhost%s", cfg.addr))

	err := http.ListenAndServe(cfg.addr, mux)

	app.logger.Error(err.Error())
	os.Exit(1)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
