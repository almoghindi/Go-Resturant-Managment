package controller

import (
	"context"
	"fmt"
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

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "get_orders_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while listing order items")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing order items"})
			return
		}

		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "get_orders_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while retrieving all orders")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while retrieving all orders"})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event": "get_orders_success",
			"time":  time.Now().Format(time.RFC3339),
		}).Info("Successfully retrieved all orders")
		c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event":   "get_order_error",
				"time":    time.Now().Format(time.RFC3339),
				"order_id": orderId,
				"error":   err,
			}).Error("Error occurred while fetching the order")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the order"})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event":   "get_order_success",
			"time":    time.Now().Format(time.RFC3339),
			"order_id": orderId,
		}).Info("Successfully retrieved order")
		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var table models.Table
		var order models.Order
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&order); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "create_order_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while binding JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(order)
		if validationErr != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "create_order_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": validationErr,
			}).Error("Validation error")
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if order.Table_id != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			if err != nil {
				msg := fmt.Sprintf("message:Table was not found")
				appLogger.Log.WithFields(logrus.Fields{
					"event":    "create_order_error",
					"time":     time.Now().Format(time.RFC3339),
					"table_id": order.Table_id,
					"error":    err,
				}).Error(msg)
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
		}

		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()

		result, insertErr := orderCollection.InsertOne(ctx, order)
		if insertErr != nil {
			msg := fmt.Sprintf("order item was not created")
			appLogger.Log.WithFields(logrus.Fields{
				"event": "create_order_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": insertErr,
			}).Error(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		appLogger.Log.WithFields(logrus.Fields{
			"event":   "create_order_success",
			"time":    time.Now().Format(time.RFC3339),
			"order_id": order.Order_id,
		}).Info("Successfully created order item")
		c.JSON(http.StatusOK, result)
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var table models.Table
		var order models.Order
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var updateObj primitive.D

		orderId := c.Param("order_id")
		if err := c.BindJSON(&order); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "update_order_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while binding JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if order.Table_id != nil {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			if err != nil {
				msg := fmt.Sprintf("message:Table was not found")
				appLogger.Log.WithFields(logrus.Fields{
					"event":    "update_order_error",
					"time":     time.Now().Format(time.RFC3339),
					"table_id": order.Table_id,
					"error":    err,
				}).Error(msg)
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{"table_id", order.Table_id})
		}

		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", order.Updated_at})

		upsert := true
		filter := bson.M{"order_id": orderId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := fmt.Sprintf("order item update failed")
			appLogger.Log.WithFields(logrus.Fields{
				"event":   "update_order_error",
				"time":    time.Now().Format(time.RFC3339),
				"order_id": orderId,
				"error":   err,
			}).Error(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		appLogger.Log.WithFields(logrus.Fields{
			"event":   "update_order_success",
			"time":    time.Now().Format(time.RFC3339),
			"order_id": orderId,
		}).Info("Successfully updated order item")
		c.JSON(http.StatusOK, result)
	}
}

func OrderItemOrderCreator(order models.Order) string {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	_, err := orderCollection.InsertOne(ctx, order)
	if err != nil {
		appLogger.Log.WithFields(logrus.Fields{
			"event": "order_item_order_creator_error",
			"time":  time.Now().Format(time.RFC3339),
			"error": err,
		}).Error("Error occurred while creating order item")
		return ""
	}

	appLogger.Log.WithFields(logrus.Fields{
		"event":   "order_item_order_creator_success",
		"time":    time.Now().Format(time.RFC3339),
		"order_id": order.Order_id,
	}).Info("Successfully created order item")
	return order.Order_id
}
