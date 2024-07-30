package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define the http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the http server
	if err := srv.ListenAndServe(); err != nil {
		log.Panic("server failed to start: %v", err)
	}
}
