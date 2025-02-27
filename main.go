package main

import (
	"os"
	"time"

	"golang-restaurant-management/database"
	"golang-restaurant-management/logger"
	"golang-restaurant-management/middleware"
	"golang-restaurant-management/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	logger.Init()

	router := gin.New()
	router.Use(gin.LoggerWithWriter(logger.Log.Out))
	router.Use(gin.RecoveryWithWriter(logger.Log.Out))

 	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"}, 
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	logger.Log.WithFields(logrus.Fields{
		"event": "application_start",
		"time":  time.Now().Format(time.RFC3339),
	}).Info("Application started")

	router.Run(":" + port)
}
