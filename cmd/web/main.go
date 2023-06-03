package main

import (
	"flag"
	"log"
	"net/http"

	"reading.maniizzle.io/internal/models"
)

type application struct {
	readingList *models.ReadingListModel
}

func main() {
	addr := flag.String("addr", ":80", "http network address")
	endpoint := flag.String("endpoint", "http://localhost:4000/v1/books", "Endpoint for the readinglist web service")

	app := &application{
		readingList: &models.ReadingListModel{
			Endpoint: *endpoint,
		},
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	log.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
