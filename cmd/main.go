package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network address")
	flag.Parse()

	// colors for logs
	colorRed := "\033[31m"
	colorGreen := "\033[32m"

	infoLog := log.New(os.Stdout, colorGreen+"INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, colorRed+"ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	infoLog.Printf("Starting the server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	errorLog.Println(err)
}
