package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// 【本地环境启动】
//
// sudo /opt/homebrew/opt/mongodb-community@7.0/bin/mongod -f /opt/homebrew/etc/mongod.conf
// sudo /opt/homebrew/opt/mongodb-community@7.0/bin/mongod -f /opt/homebrew/etc/mongod-27020.conf
// sudo /opt/homebrew/opt/mongodb-community@7.0/bin/mongod -f /opt/homebrew/etc/mongod-27021.conf
//
// 【连接uri】
// mongosh "mongodb://localhost:27017,localhost:27020,localhost:27021/?replicaSet=rs0"

// ConnectMongo will connect to a mongodb replica set.
//
// The main purpose here is to test how driver behavior in a
// combination of using the clientOptions and connect URI.
//
// Cause in the sit env in Pingan Tech company, things are
// acting weird, try to reproduce them locally to figure out
// what's happing.
func ConnectMongo(uri, readPreference, readConcern string) {
	if uri == "" {
		log.Fatal("URI empty")
	}

	// handle options
	mgOptions := &options.ClientOptions{}
	mgOptions.SetMaxPoolSize(1)
	switch readPreference {
	case "primary":
		mgOptions.ReadPreference = readpref.Primary()
	case "secondary":
		mgOptions.ReadPreference = readpref.Secondary()
	default:
		mgOptions.ReadPreference = readpref.SecondaryPreferred()
	}

	switch readConcern {
	case "majority":
		mgOptions.ReadConcern = readconcern.Majority()
	case "linearizable":
		mgOptions.ReadConcern = readconcern.Linearizable()
	default:
		mgOptions.ReadConcern = readconcern.Local()
	}

	// construct client instance
	client, err := mongo.NewClient(options.Client().ApplyURI(uri), mgOptions)
	if err != nil {
		log.Fatal(err)
	}

	// connect to mongodb
	if err = client.Connect(context.TODO()); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	log.Println("read_concern: ", mgOptions.ReadConcern.GetLevel())
	log.Println("read_preference: ", mgOptions.ReadPreference.Mode().String())
	res := client.Database("test").Collection("test").FindOne(context.TODO(), bson.M{})
	if res.Err() != nil {
		log.Fatal(res.Err().Error())
	}
	data, err := res.DecodeBytes()
	if err != nil {
		log.Fatal(err)
	}

	// fetch data
	user := struct {
		Name string `json:"name" bson:"name"`
	}{}
	bson.Unmarshal(data, &user)
	log.Println(">>>>>>", user)
}
