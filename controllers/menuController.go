package controller

import (
	"context"
	"golang-restaurant-management/database"
	appLogger "golang-restaurant-management/logger"
	"golang-restaurant-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := menuCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "get_menus_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while listing the menu items")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing the menu items"})
			return
		}

		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "get_menus_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while retrieving all menus")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while retrieving all menus"})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event": "get_menus_success",
			"time":  time.Now().Format(time.RFC3339),
		}).Info("Successfully retrieved all menus")
		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event":   "get_menu_error",
				"time":    time.Now().Format(time.RFC3339),
				"menu_id": menuId,
				"error":   err,
			}).Error("Error occurred while fetching the menu")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the menu"})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event":   "get_menu_success",
			"time":    time.Now().Format(time.RFC3339),
			"menu_id": menuId,
		}).Info("Successfully retrieved menu")
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&menu); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "create_menu_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while binding JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "create_menu_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": validationErr,
			}).Error("Validation error")
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			msg := "Menu item was not created"
			appLogger.Log.WithFields(logrus.Fields{
				"event": "create_menu_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": insertErr,
			}).Error(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event":   "create_menu_success",
			"time":    time.Now().Format(time.RFC3339),
			"menu_id": menu.Menu_id,
		}).Info("Successfully created menu item")
		c.JSON(http.StatusOK, result)
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.Before(check) && end.After(check)
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "update_menu_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while binding JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}

		var updateObj primitive.D

		if menu.Start_Date != nil && menu.End_Date != nil {
			if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()) {
				msg := "Kindly retype the time"
				appLogger.Log.WithFields(logrus.Fields{
					"event": "update_menu_error",
					"time":  time.Now().Format(time.RFC3339),
					"error": msg,
				}).Error(msg)
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}

			updateObj = append(updateObj, bson.E{"start_date", menu.Start_Date})
			updateObj = append(updateObj, bson.E{"end_date", menu.End_Date})
		}

		if menu.Name != "" {
			updateObj = append(updateObj, bson.E{"name", menu.Name})
		}
		if menu.Category != "" {
			updateObj = append(updateObj, bson.E{"category", menu.Category})
		}

		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", menu.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := menuCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := "Menu update failed"
			appLogger.Log.WithFields(logrus.Fields{
				"event":   "update_menu_error",
				"time":    time.Now().Format(time.RFC3339),
				"menu_id": menuId,
				"error":   err,
			}).Error(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event":   "update_menu_success",
			"time":    time.Now().Format(time.RFC3339),
			"menu_id": menuId,
		}).Info("Successfully updated menu item")
		c.JSON(http.StatusOK, result)
	}
}
