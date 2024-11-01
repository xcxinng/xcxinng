package main

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func TestConnectMongo(t *testing.T) {
	type args struct {
		uri            string
		readPreference string
		readConcern    string
	}
	tests := []struct {
		name string
		args args
	}{
		// {
		// 	name: "use the value from uri",
		// 	args: args{
		// 		uri:            "mongodb://localhost:27020,localhost:27021/?replicaSet=rs0",
		// 		readPreference: "",
		// 		readConcern:    ""},
		// },
		// {
		// 	name: "use the value from manual specifying",
		// 	args: args{
		// 		uri:            "mongodb://localhost:27020,localhost:27021/?replicaSet=rs0",
		// 		readPreference: "primary",
		// 		readConcern:    "majority"},
		// },
		{
			name: "use the value from the combination",
			args: args{
				uri: "mongodb://localhost:27017/testdb",
				// readPreference: "secondary",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConnectMongo(tt.args.uri, tt.args.readPreference, tt.args.readConcern)
		})
	}
}

func TestTransaction(t *testing.T) {
	err := WithTransactionExample(context.Background())
	assert.NoError(t, err)
}

// WithTransactionExample is an example of using the Session.WithTransaction function.
func WithTransactionExample(ctx context.Context) error {
	// For a replica set, include the replica set name and a seedlist of the members in the URI string; e.g.
	// uri := "mongodb://mongodb0.example.com:27017,mongodb1.example.com:27017/?replicaSet=myRepl"
	// For a sharded cluster, connect to the mongos instances; e.g.
	// uri := "mongodb://mongos0.example.com:27017,mongos1.example.com:27017/"
	uri := "mongodb://localhost:27020,localhost:27021/?replset=rs0"
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return err
	}
	defer func() { _ = client.Disconnect(ctx) }()
	// Prereq: Create collections.
	wcMajority := writeconcern.New(writeconcern.WMajority(), writeconcern.WTimeout(1*time.Second))
	wcMajorityCollectionOpts := options.Collection().SetWriteConcern(wcMajority)
	fooColl := client.Database("test").Collection("foo", wcMajorityCollectionOpts)
	barColl := client.Database("test").Collection("bar", wcMajorityCollectionOpts)
	// Step 1: Define the callback that specifies the sequence of operations to perform inside the transaction.
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Important: You must pass sessCtx as the Context parameter to the operations for them to be executed in the
		// transaction.
		if _, err := fooColl.InsertOne(sessCtx, bson.D{{"hello", "world"}}); err != nil {
			return nil, err
		}
		return nil, errors.New("error test")
		// return nil,fmt.Errorf("test error")
		if _, err := barColl.InsertOne(sessCtx, bson.D{{"xyz", 999}}); err != nil {
			return nil, err
		}
		return nil, nil
	}
	// Step 2: Start a session and run the callback using WithTransaction.
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	result, err := session.WithTransaction(ctx, callback)
	if err != nil {
		return err
	}
	log.Printf("result: %v\n", result)
	return nil
}
