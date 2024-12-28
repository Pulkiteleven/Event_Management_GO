package controller

import (
	"context"
	"go_event/database"
	"go_event/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CategoryCollection *mongo.Collection = database.OpenCollection(database.Client, "category")

func CreateCategroy() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var categroy models.Category

		if err := c.BindJSON(&categroy); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(categroy)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation error"})
			defer cancel()
			return
		}

		count, err := CategoryCollection.CountDocuments(ctx, bson.M{
			"name": categroy.Name,
		})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The category already exists"})
			return
		}

		categroy.Created_at, _ = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))
		categroy.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))

		categroy.ID = primitive.NewObjectID()
		categroy.Cat_id = categroy.ID.Hex()

		resultInsertionNumber, insertErr := CategoryCollection.InsertOne(ctx,categroy)

		if insertErr != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Categroy Item was not inserted"})
		}
		defer cancel()

		c.JSON(http.StatusOK,resultInsertionNumber)
	}
}

func GetCategory() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(),100*time.Second)

		defer cancel()

		cursor, err := CategoryCollection.Find(ctx,bson.M{},options.Find())

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error Getting Category"})
			return
		}

		var categories []models.Category

		if err := cursor.All(ctx,&categories);err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error decoding categories"})
			return
		}

		c.JSON(http.StatusOK,gin.H{
			"categories":categories,
		})
	}
}


