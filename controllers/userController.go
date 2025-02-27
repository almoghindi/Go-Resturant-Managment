package controller

import (
	"context"
	"golang-restaurant-management/database"
	helper "golang-restaurant-management/helpers"
	appLogger "golang-restaurant-management/logger"
	"golang-restaurant-management/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

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

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, projectStage})
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "get_users_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while listing user items")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing user items"})
			return
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "get_users_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while retrieving user items")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while retrieving user items"})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event": "get_users_success",
			"time":  time.Now().Format(time.RFC3339),
		}).Info("Successfully retrieved user items")
		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")
		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event":   "get_user_error",
				"time":    time.Now().Format(time.RFC3339),
				"user_id": userId,
				"error":   err,
			}).Error("Error occurred while fetching the user")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the user"})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event":   "get_user_success",
			"time":    time.Now().Format(time.RFC3339),
			"user_id": userId,
		}).Info("Successfully retrieved user")
		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "sign_up_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while binding JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "sign_up_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": validationErr,
			}).Error("Validation error")
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "sign_up_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while checking for the email")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "sign_up_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while checking for the phone number")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the phone number"})
			return
		}

		if count > 0 {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "sign_up_error",
				"time":  time.Now().Format(time.RFC3339),
			}).Error("This email or phone number already exists")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := "User item was not created"
			appLogger.Log.WithFields(logrus.Fields{
				"event": "sign_up_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": insertErr,
			}).Error(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event":   "sign_up_success",
			"time":    time.Now().Format(time.RFC3339),
			"user_id": user.User_id,
		}).Info("Successfully created user item")
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "login_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while binding JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event":   "login_error",
				"time":    time.Now().Format(time.RFC3339),
				"user_id": user.User_id,
				"error":   err,
			}).Error("User not found, login seems to be incorrect")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found, login seems to be incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "login_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": msg,
			}).Error("Login or password is incorrect")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		appLogger.Log.WithFields(logrus.Fields{
			"event":   "login_success",
			"time":    time.Now().Format(time.RFC3339),
			"user_id": foundUser.User_id,
		}).Info("Successfully logged in")
		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		appLogger.Log.WithFields(logrus.Fields{
			"event": "hash_password_error",
			"time":  time.Now().Format(time.RFC3339),
			"error": err,
		}).Panic("Error occurred while hashing password")
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "login or password is incorrect"
		check = false
	}
	return check, msg
}
