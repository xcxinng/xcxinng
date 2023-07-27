package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	var uri = "mongodb://localhost:27020,localhost:27021/?replset=rs0"

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	// database := client.Database("test")
	// coll1 := database.Collection("test1")
	// coll2 := database.Collection("test2")

	// start-session
	// wc := writeconcern.New(writeconcern.WMajority())
	// txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.TODO())

	// result, err := session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
	// 	result, err := coll.InsertMany(ctx, []interface{}{
	// 		bson.D{{"title", "The Bluest Eye"}, {"author", "Toni Morrison"}},
	// 		bson.D{{"title", "Sula"}, {"author", "Toni Morrison"}},
	// 		bson.D{{"title", "Song of Solomon"}, {"author", "Toni Morrison"}},
	// 	})
	// 	return result, err
	// }, txnOptions)
}
