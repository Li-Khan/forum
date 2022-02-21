package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Li-Khan/forum/pkg/models"
	"github.com/Li-Khan/forum/pkg/models/sqlite"
	"github.com/Li-Khan/forum/web"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network address")
	dsn := flag.String("dsn", "forum.db", "Name of the database")
	flag.Parse()

	// цвета для логгирования
	colorRed := "\033[31m"
	colorGreen := "\033[32m"

	infoLog := log.New(os.Stdout, colorGreen+"INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, colorRed+"ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := models.OpenDB(*dsn)
	if err != nil {
		errorLog.Println(err)
		return
	}

	// Инициализирую кэш шаблонов
	templateCache, err := web.NewTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Println(err)
		return
	}

	// Зависимости всего веб приложения.
	app := web.Application{
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
		Forum:         &sqlite.ForumModel{DB: db},
		TemplateCache: templateCache,
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go web.SessionGC()

	infoLog.Printf("Starting the server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Println(err)
}
