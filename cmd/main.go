package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Li-Khan/forum/web"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network address")
	flag.Parse()

	// colors for logs
	colorRed := "\033[31m"
	colorGreen := "\033[32m"

	infoLog := log.New(os.Stdout, colorGreen+"INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, colorRed+"ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := web.Application{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting the server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Println(err)
}
