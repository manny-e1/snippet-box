package main

import (
	"database/sql"
	"flag"
	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql"
	"github.com/manny-e1/snippetbox/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type application struct {
	errorLogger        *log.Logger
	infoLogger         *log.Logger
	snippets           *models.SnippetModel
	transactionExample *models.TransactionExample
	templateCache      map[string]*template.Template
	formDecoder        *form.Decoder
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	addr := flag.String("addr", ":5000", "HTTP Network Address")
	dsn := flag.String("dsn", "letsgo:LetsGo123!@/snippetbox?parseTime=true", "Mysql datasource name")

	flag.Parse()

	var infoLogger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	var errorLogger = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)

	templateCache, err := newTemplateCache()

	if err != nil {
		errorLogger.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	app := &application{
		errorLogger: errorLogger,
		infoLogger:  infoLogger,
		snippets: &models.SnippetModel{
			DB: db,
		},
		transactionExample: &models.TransactionExample{DB: db},
		templateCache:      templateCache,
		formDecoder:        formDecoder,
	}

	if err != nil {
		errorLogger.Fatal(err)
	}

	defer db.Close()

	mux := app.routes()
	infoLogger.Printf("Starting server on %s ", *addr)

	srv := &http.Server{
		ErrorLog: errorLogger,
		Addr:     *addr,
		Handler:  mux}
	err = srv.ListenAndServe()
	errorLogger.Fatal(err)
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
