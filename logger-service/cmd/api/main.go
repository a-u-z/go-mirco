package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = ":80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

type Config struct {
	Models Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// @@@@@@@@@@@@@@
	s := NewServer()
	s.Models = New(client)

	s.router.Run(webPort)
	// @@@@@@@@@@@@@@

}

// func (app *Config) serve() {
// 	srv := &http.Server{
// 		Addr: fmt.Sprintf(":%s", webPort),
// 		Handler: app.routes(),
// 	}

// 	err := srv.ListenAndServe()
// 	if err != nil {
// 		log.Panic()
// 	}
// }

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}
