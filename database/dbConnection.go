package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Dbinstance() *mongo.Client {
	// MongoDb := "mongodb://localhost:27017/"
	MongoDb := "mongodb+srv://pulkitdubey1220:mZCFtKiiK2WODCdj@pulkit.9jagb.mongodb.net/"
	fmt.Println(MongoDb)

	client,err := mongo.NewClient(options.Client().ApplyURI(MongoDb))

	if err != nil{
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Connected to Mongodb")
	return client
}

var Client *mongo.Client = Dbinstance()


func OpenCollection(client *mongo.Client, collectionName string, dbName ...string) *mongo.Collection {
	// Use "event" as the default database name if none is provided
	database := "event"
	if len(dbName) > 0 && dbName[0] != "" {
		database = dbName[0]
	}

	var collection *mongo.Collection = client.Database(database).Collection(collectionName)
	return collection
}