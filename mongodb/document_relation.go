package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func genIpaddr() string {
	md := rand.New(rand.NewSource(time.Now().Unix()))
	ip := fmt.Sprintf("%d.%d.%d.%d", md.Intn(255), md.Intn(255), md.Intn(255), md.Intn(255))
	return ip
}

func insertPort(d *mongo.Database) error {
	for i := 0; i < 100; i++ {
		var documents []interface{}
		for j := 0; j < 5000; j++ {
			documents = append(documents, bson.M{
				"_id":         primitive.NewObjectID(),
				"port_name":   fmt.Sprintf("port_%d_%d", i, j),
				"port_ip":     genIpaddr(),
				"description": "bla bla ...",
			})
		}
		_, err := d.Collection("instance_port").InsertMany(context.TODO(), documents)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertDocument() error {
	var uri = "mongodb://localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())
	var portIDs []bson.M
	cursor, err := client.Database("test").Collection("instance_port").
		Find(context.TODO(), bson.M{}, options.Find())
	if err != nil {
		return err
	}

	err = cursor.All(context.TODO(), &portIDs)
	if err != nil {
		return err
	}
	var tID = make([]primitive.ObjectID, 0, len(portIDs))
	for _, v := range portIDs {
		id := v["_id"]
		tID = append(tID, id.(primitive.ObjectID))
	}

	document := bson.M{
		"device_name":         "device_100",
		"admin_ip":            "1.1.1.1",
		"kind":                "switch",
		"vendor":              "h3c",
		"device_contain_port": tID,
	}
	_, err = client.Database("test").Collection("instance_device").InsertOne(context.TODO(), document)
	return err
}
