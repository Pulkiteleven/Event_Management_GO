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

var AttendanceCollection *mongo.Collection = database.OpenCollection(database.Client,"attendance")


func MarkAttendance() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var attendance models.Attendance

		if err := c.BindJSON(&attendance); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(attendance)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			defer cancel()
			return
		}

		startOfDay := time.Now().Truncate(24 * time.Hour)
		endOfDay := startOfDay.Add(24 * time.Hour)

		enrollCount, err := EnrollCollection.CountDocuments(ctx,bson.M{
			"user_id":attendance.User_id,
			"event_id":attendance.Event_id,
		})
		defer cancel()

		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error"})
			return
		}
		if enrollCount <= 0{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Not Enrolled in This Event"})
			return
		}

		count, err := AttendanceCollection.CountDocuments(ctx, bson.M{
			"user_id": attendance.User_id,
			"event_id":attendance.Event_id,
			"attendance_date": bson.M{
				"$gte": startOfDay, // Greater than or equal to the start of the day
				"$lt":  endOfDay,   // Less than the start of the next day
			},
		})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Attendance for this day already marked"})
			return
		}

		attendance.Created_at, _ = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))
		attendance.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))
		attendance.Attendance_date, _ = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))

		attendance.ID = primitive.NewObjectID()
		attendance.Attendance_id = attendance.ID.Hex()

		resultInsertionNumber, insertErr := AttendanceCollection.InsertOne(ctx, attendance)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Attendnace was not inserted"})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func GetAttendance() gin.HandlerFunc {
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


		cursor, err := AttendanceCollection.Find(ctx,bson.M{"event_id":eventId},options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(recordPerPage)))

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error Getting Attendance"})
			return
		}
		var attended []models.Attendance

		if err := cursor.All(ctx,&attended); err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error Decoding Attendance"})
			return
		}

		var userIds []string
		dateMap := make(map[string]time.Time)
		for _,attendedUser := range attended{
			userIds = append(userIds, attendedUser.User_id)
			dateMap[attendedUser.User_id] = attendedUser.Attendance_date
		}

		var users []models.User
		var result[] map[string]interface{}

		if len(userIds) > 0{
			cursor,err = userCollection.Find(ctx,bson.M{"user_id":bson.M{"$in":userIds}})
			if err != nil{
				c.JSON(http.StatusInternalServerError,gin.H{"error":"Error fetching users"})
				return
			}

			if err := cursor.All(ctx,&users); err != nil{
				c.JSON(http.StatusInternalServerError,gin.H{"error":"Error Decoding User Data"})
				return
			}

			for _,user := range users{
				output := map[string]interface{}{
					"user_name": user.User_name,
					"email":     user.Email,
					"avatar":    user.Avatar,
					"attendance_date":dateMap[user.User_id],
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
