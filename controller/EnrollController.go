package controller

import (
	"context"
	"go_event/database"
	"go_event/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var EnrollCollection *mongo.Collection = database.OpenCollection(database.Client, "enroll")


func EnrollUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var enroll models.Enroll

		if err := c.BindJSON(&enroll); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(enroll)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			defer cancel()
			return
		}


		count, err := EnrollCollection.CountDocuments(ctx, bson.M{
			"user_id": enroll.User_id,
			"event_id":enroll.Event_id})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "You are Already Enrolled"})
			return
		}

		enroll.Created_at, _ = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))
		enroll.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))

		enroll.ID = primitive.NewObjectID()
		enroll.Enroll_id = enroll.ID.Hex()

		resultInsertionNumber, insertErr := EnrollCollection.InsertOne(ctx, enroll)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Enroll item was not inserted"})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func ApprovedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var enroll models.Enroll
		userId := c.Param("user_id")
		eventId := c.Param("event_id")

		filter := bson.M{"user_id": userId,"event_id":eventId}

		var updateObj primitive.D
		updateObj = append(updateObj, bson.E{"approved", true})
		enroll.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", enroll.Updated_at})

		usert := true

		opt := options.UpdateOptions{
			Upsert: &usert,
		}

		result, err := EnrollCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not approved"})
		}

		defer cancel()

		c.JSON(http.StatusOK, result)

	}
}

func GetEnrollUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()
		eventId := c.Param("event_id")
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))

		if err1 != nil || page < 1 {
			page = 1
		}

		skip := (page - 1) * recordPerPage

		cursor, err := EnrollCollection.Find(ctx,bson.M{"event_id":eventId},options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(recordPerPage)))

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Error Getting enrolled User"})
			return
		}

		var enrolled []models.Enroll

		if err := cursor.All(ctx,&enrolled); err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error decoding enrolled user"})
			return
		}

		var userIds []string
		approvedMap := make(map[string]bool)
		for _, enrolledUser := range enrolled{
			userIds = append(userIds, enrolledUser.User_id)
			approvedMap[enrolledUser.User_id] = enrolledUser.Approved
		}

		var users []models.User
		var result[]map[string]interface{}

		if len(userIds) > 0{
			cursor, err = userCollection.Find(ctx,bson.M{"user_id":bson.M{"$in":userIds}})
			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error":"Error fetching User details"})
				return
			}

			if err := cursor.All(ctx,&users); err != nil{
				c.JSON(http.StatusInternalServerError,gin.H{"error":"Error Decoding User Data"})
				return
			}

			for _, user := range users {
				output := map[string]interface{}{
					"user_name": user.User_name,
					"email":     user.Email,
					"avatar":    user.Avatar,
					"approved":  approvedMap[user.User_id],
				}
				result = append(result, output)
			}
		}

		c.JSON(http.StatusOK,gin.H{
			"page":page,
			"record_count":len(users),
			"users":result,
		})

	
	}
}

