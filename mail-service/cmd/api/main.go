package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
}

const webport = "80"

func main() {
	app := Config{}

	log.Println("Starting mail service on port", webport)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webport),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
