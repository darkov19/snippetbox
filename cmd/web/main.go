package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	// Import the models package that we just created. You need to prefix this with
	// whatever module path you set up back in chapter 02.01 (Project Setup and Creating
	// a Module) so that the import statement looks like this:
	// "{your-module-path}/internal/models". If you can't remember what module path you
	// used, you can find it at the top of the go.mod file.
	"github.com/darkov19/snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

var cfg config

func main() {
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "web:1920@tcp(127.0.0.1:33061)/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	var app = &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	app.logger.Info("Starting server on", "addr", fmt.Sprintf("http://localhost%s", cfg.addr))

	err = http.ListenAndServe(cfg.addr, app.routes())

	app.logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
