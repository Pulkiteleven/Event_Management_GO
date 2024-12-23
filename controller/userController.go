package controller

import (
	"context"
	"fmt"
	"go_event/database"
	"go_event/helpers"
	"go_event/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func GetUsers() gin.HandlerFunc {
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

		cursor, err := userCollection.Find(ctx, bson.M{}, options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(recordPerPage)))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
			return
		}

		var users []bson.M
		if err := cursor.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error error decoding user data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"page":         page,
			"record_count": len(users),
			"users":        users,
		})
	}
}

// func GetUsers() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
// 		if err != nil || recordPerPage < 1 {
// 			recordPerPage = 10
// 		}

// 		page, err1 := strconv.Atoi(c.Query("page"))

// 		if err1 != nil || page < 1 {
// 			page = 1
// 		}

// 		startIndex := (page - 1) * recordPerPage
// 		startIndex, err = strconv.Atoi(c.Query("startIndex"))

// 		matchStage := bson.D{{"$match", bson.D{{}}}}
// 		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", "1"}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
// 		projectStage := bson.D{
// 			{"$project", bson.D{
// 				{"_id", 0},
// 				{"total_count", 1},
// 				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
// 			}}}

// 		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
// 			matchStage, groupStage, projectStage})
// 		defer cancel()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
// 		}

// 		var allUsers []bson.M
// 		if err = result.All(ctx, &allUsers); err != nil {
// 			log.Fatal(err)
// 		}
// 		c.JSON(http.StatusOK, allUsers[0])
// 	}
// }

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		userId := c.Param("user_id")
		fmt.Println("hello",userId)

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not Found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			defer cancel()
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			defer cancel()
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email Already exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.User_name, user.User_id)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			// msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User item was not created"})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		var user models.User
		var foundUser models.User

		defer cancel()
		if err := c.BindJSON(&user);err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		err := userCollection.FindOne(ctx,bson.M{"email":user.Email}).Decode(&foundUser)

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"user not found"})
		}


		passwordIsValid, msg := VerifyPassword(*user.Password,*foundUser.Password)

		if passwordIsValid != true{
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		token, referehToken, _ := helpers.GenerateAllTokens(foundUser.Email,foundUser.User_name,foundUser.User_id)

		helpers.UpdateAllTokens(token,referehToken,foundUser.User_id)

		defer cancel()

		c.JSON(http.StatusOK,foundUser)
	}
}

func UpdateUser() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		var user models.User
		userId := c.Param("user_id")
		fmt.Printf(userId);

		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest,err.Error())
			return
		}

		filter := bson.M{"user_id":userId}
		
		var updateObj primitive.D
		
			updateObj = append(updateObj, bson.E{"avatar",user.Avatar})
			updateObj = append(updateObj, bson.E{"city",user.City})
			updateObj = append(updateObj, bson.E{"categories",user.Categories})
		
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", user.Updated_at})

		upsert := true
		
		opt := options.UpdateOptions{
			Upsert:&upsert,
		}

		result,err := userCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set",updateObj},
			},
			&opt,
		)

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"User not Updated"})
			return
		}

		defer cancel()


		c.JSON(http.StatusOK,result)



	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword, providedPassword string) (bool,string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil{
		msg = fmt.Sprintf("login or password is incorrect")
		check = false
	}

	return check,msg
}
