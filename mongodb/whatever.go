package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client())
	if err != nil {
		log.Fatal(err)
	}

	client.Database("test").Collection("student").InsertOne(context.TODO(), nil)

}
