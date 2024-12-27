package controller

import (
	"context"
	"fmt"
	"go_event/database"
	"go_event/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var eventCollection *mongo.Collection = database.OpenCollection(database.Client, "event")

func CreateEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var event models.Event

		if err := c.BindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(event)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			defer cancel()
			return
		}

		event.Created_at, _ = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))
		event.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))

		event.ID = primitive.NewObjectID()
		event.Event_id = event.ID.Hex()

		resultInsertionNumber, inserErr := eventCollection.InsertOne(ctx, event)

		if inserErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Event was not created"})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func GetEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))

		if err1 != nil || page < 1 {
			page = 1
		}

		skip := (page - 1) * recordPerPage

		filter := bson.M{}

		if userId := c.Query("user_id"); userId != "" {
			filter["owner"] = userId
		}
		if category := c.Query("category"); category != "" {
			filter["category"] = category
		}
		if city := c.Query("city"); city != "" {
			filter["city"] = city
		}
		fmt.Println(filter)

		cursor, err := eventCollection.Find(ctx, filter, options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(recordPerPage)))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Fetching Events"})
			return
		}

		var events []bson.M
		if err := cursor.All(ctx, &events); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding event data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"page":         page,
			"record_count": len(events),
			"events":       events,
		})

	}
}

func GetUserEnrolledEvents() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(),100*time.Second)

		defer cancel()
		
		userId := c.Param("user_id")

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))

		if err1 != nil || page < 1 {
			page = 1
		}

		skip := (page - 1) * recordPerPage

		cursor,err := EnrollCollection.Find(ctx,bson.M{"user_id":userId},options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(recordPerPage)))

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Errr Getting Events"})
			return
		}

		var enrolled []models.Enroll

		if err := cursor.All(ctx,&enrolled);err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error decoding enrolled user"})
			return
		}

		var eventIds []string
		for _,enrolledUser := range enrolled{
			eventIds = append(eventIds, enrolledUser.Event_id)
		}

		var Events[] models.Event

		if len(eventIds) > 0{
			cursor, err = eventCollection.Find(ctx,bson.M{"event_id":bson.M{"$in":eventIds}})
			if err != nil{
				c.JSON(http.StatusInternalServerError,gin.H{"error":"Error fetching User details"})
				return
			}

			if err := cursor.All(ctx,&Events); err != nil{
				c.JSON(http.StatusInternalServerError,gin.H{"error":"Error Decoding Events"})
			}
		}

		c.JSON(http.StatusOK,gin.H{
			"page":page,
			"record_count":len(Events),
			"user":Events,
		})


	}
}
