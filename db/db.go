package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Client *mongo.Client
	CTX    context.Context
	Cancel context.CancelFunc
)

// Opens a connection to the DB globally
func Connect(uri string) error {
	var err error
	CTX, Cancel = context.WithCancel(context.Background())
	Client, err = mongo.Connect(CTX, options.Client().ApplyURI(uri))
	return err
}

// Closes the global connetion to the DB
func Close() {
	defer Cancel()
	defer Client.Disconnect(CTX)
}

// Checks db connection
func Ping() error {
	return Client.Ping(CTX, readpref.Primary())
}

// Drops the DB if it exists and then creates both the DB and the collection anew
func Reset() error {
	db := Client.Database(os.Getenv("DB_NAME"))

	err := db.Drop(CTX)
	if err != nil {
		return err
	}

	err = db.CreateCollection(CTX, os.Getenv("COLL_NAME"))
	if err != nil {
		return err
	}

	collection := db.Collection(os.Getenv("COLL_NAME"))

	expIndex := mongo.IndexModel{
		Keys:    bson.M{"exp": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	_, err = collection.Indexes().CreateOne(CTX, expIndex)
	if err != nil {
		return err
	}

	guidIndex := mongo.IndexModel{
		Keys:    bson.M{"guid": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(CTX, guidIndex)
	if err != nil {
		return err
	}

	return nil
}

// Returns the collection that has COLL_NAME
func GetTokenCollection() *mongo.Collection {
	return Client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLL_NAME"))
}
