package main

import (
	"authentification-service/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentification service on port", webPort)

	// TODO connect to DB
	conn, err := connectToDB()
	if err != nil {
		log.Panicf("Could not connect to DB: %v", err)
	}

	//set up config

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panicf("server failed to start: %v", err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() (*sql.DB, error) {
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Could not connect to DB. Retrying in 2 seconds")
			counts++
		} else {
			log.Println("Connected to DB")
			return connection, nil
		}

		if counts >= 5 {
			return nil, fmt.Errorf("Connection lost after %d retries", counts)
		}

		time.Sleep(2 * time.Second)
	}
}
