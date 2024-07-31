package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"logger-service/data"
	"net/http"
	"time"
)

const (
	webport  = "80"
	rpcport  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	//connect to MongoDB
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	//create a context in order to disconnect from MongoDB

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	//disconnect from MongoDB
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	//start the web server
	//go app.serve()

	log.Println("Starting server on port", webport)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webport),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

//func (app *Config) serve() {
//	srv := &http.Server{
//		Addr:    fmt.Sprintf(":%s", webport),
//		Handler: app.routes(),
//	}
//
//	err := srv.ListenAndServe()
//	if err != nil {
//		log.Panic()
//	}
//}

func connectToMongo() (*mongo.Client, error) {
	//create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to MongoDB:", err)
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	return client, nil
}
