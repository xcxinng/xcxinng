package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// 验证writeConcern=Local(w=1)时，主从之间数据的延后性，复制有延迟
func VerifyLocalWriteConcern() {
	var uri = "mongodb://localhost:27017,localhost:27020,localhost:27021/?replicaSet=rs0"
	opt := options.Client().SetTimeout(time.Second * 5)
	client, err := mongo.Connect(context.Background(), opt.ApplyURI(uri))
	if err != nil {
		fmt.Println("111")
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	w := options.Collection().SetWriteConcern(writeconcern.New(writeconcern.W(1)))
	res, err := client.Database("test").Collection("users", w).
		InsertOne(context.TODO(), bson.M{"name": "xcx", "age": 10, "gender": "male"})
	if err != nil {
		panic(err)
	}
	fmt.Println("insertedID: ", res.InsertedID)
	// output: (example)
	// insertedID:  ObjectID("673d8d9e917e0b1abde21dd7")

	r := options.Collection().
		SetReadConcern(readconcern.Local()).
		SetReadPreference(readpref.SecondaryPreferred())
	cursor := client.Database("test").Collection("users", r).
		FindOne(context.TODO(), bson.M{"name": "xcx"})
	if err != nil {
		panic(err)
	}
	var data map[string]interface{}
	err = cursor.Decode(&data)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		log.Println("no document found")
		// output:(example)
		// 2024/11/20 15:19:58 no document found
		return
	}
	fmt.Println("[SecondaryPreferred] data: ", data)
}

// func main() {
// 	VerifyLocalWriteConcern()
// }
