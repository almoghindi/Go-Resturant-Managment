package controller

import (
	"context"
	"golang-restaurant-management/database"
	appLogger "golang-restaurant-management/logger"
	"golang-restaurant-management/models"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
	validate        = validator.New()
)

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"food_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
				}}}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "get_foods_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while listing food items")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing food items"})
			return
		}
		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "get_foods_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while retrieving food items")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while retrieving food items"})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event": "get_foods_success",
			"time":  time.Now().Format(time.RFC3339),
		}).Info("Successfully retrieved food items")
		c.JSON(http.StatusOK, allFoods[0])
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		foodId := c.Param("food_id")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event":   "get_food_error",
				"time":    time.Now().Format(time.RFC3339),
				"food_id": foodId,
				"error":   err,
			}).Error("Error occurred while fetching the food item")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the food item"})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event":   "get_food_success",
			"time":    time.Now().Format(time.RFC3339),
			"food_id": foodId,
		}).Info("Successfully retrieved food item")
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "create_food_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while binding JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(food)
		if validationErr != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "create_food_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": validationErr,
			}).Error("Validation error")
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		if err != nil {
			msg := "menu was not found"
			appLogger.Log.WithFields(logrus.Fields{
				"event":   "create_food_error",
				"time":    time.Now().Format(time.RFC3339),
				"menu_id": food.Menu_id,
				"error":   err,
			}).Error(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		num := toFixed(*food.Price, 2)
		food.Price = &num

		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := "Food item was not created"
			appLogger.Log.WithFields(logrus.Fields{
				"event": "create_food_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": insertErr,
			}).Error(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		appLogger.Log.WithFields(logrus.Fields{
			"event":   "create_food_success",
			"time":    time.Now().Format(time.RFC3339),
			"food_id": food.Food_id,
		}).Info("Successfully created food item")
		c.JSON(http.StatusOK, result)
	}
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		var food models.Food

		foodId := c.Param("food_id")

		if err := c.BindJSON(&food); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "update_food_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while binding JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{"name", food.Name})
		}

		if food.Price != nil {
			updateObj = append(updateObj, bson.E{"price", food.Price})
		}

		if food.Food_image != nil {
			updateObj = append(updateObj, bson.E{"food_image", food.Food_image})
		}

		if food.Menu_id != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
			if err != nil {
				msg := "menu was not found"
				appLogger.Log.WithFields(logrus.Fields{
					"event":   "update_food_error",
					"time":    time.Now().Format(time.RFC3339),
					"menu_id": food.Menu_id,
					"error":   err,
				}).Error(msg)
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{"menu", food.Price})
		}

		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", food.Updated_at})

		upsert := true
		filter := bson.M{"food_id": foodId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := "food item update failed"
			appLogger.Log.WithFields(logrus.Fields{
				"event":   "update_food_error",
				"time":    time.Now().Format(time.RFC3339),
				"food_id": foodId,
				"error":   err,
			}).Error(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event":   "update_food_success",
			"time":    time.Now().Format(time.RFC3339),
			"food_id": foodId,
		}).Info("Successfully updated food item")
		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
