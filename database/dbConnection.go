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
	MongoDb := "mongodb://localhost:27017/"
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


func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("event").Collection(collectionName)

	return collection
}