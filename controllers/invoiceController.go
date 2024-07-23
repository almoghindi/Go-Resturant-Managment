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

type InvoiceViewFormat struct {
	Invoice_id       string
	Payment_method   string
	Order_id         string
	Payment_status   *string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_details    interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := invoiceCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "get_invoices_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while listing invoice items")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing invoice items"})
			return
		}

		var allInvoices []bson.M
		if err = result.All(ctx, &allInvoices); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "get_invoices_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while retrieving all invoices")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while retrieving all invoices"})
			return
		}
		appLogger.Log.WithFields(logrus.Fields{
			"event": "get_invoices_success",
			"time":  time.Now().Format(time.RFC3339),
		}).Info("Successfully retrieved all invoices")
		c.JSON(http.StatusOK, allInvoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		invoiceId := c.Param("invoice_id")
		var invoice models.Invoice

		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event":      "get_invoice_error",
				"time":       time.Now().Format(time.RFC3339),
				"invoice_id": invoiceId,
				"error":      err,
			}).Error("Error occurred while fetching the invoice item")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the invoice item"})
			return
		}

		var invoiceView InvoiceViewFormat
		allOrderItems, err := ItemsByOrder(invoice.Order_id)
		if err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event":      "get_invoice_error",
				"time":       time.Now().Format(time.RFC3339),
				"invoice_id": invoiceId,
				"error":      err,
			}).Error("Error occurred while fetching order items")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching order items"})
			return
		}
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date

		invoiceView.Payment_method = "null"
		if invoice.Payment_method != nil {
			invoiceView.Payment_method = *invoice.Payment_method
		}

		invoiceView.Invoice_id = invoice.Invoice_id
		invoiceView.Payment_status = invoice.Payment_status
		if len(allOrderItems) > 0 {
			invoiceView.Payment_due = allOrderItems[0]["payment_due"]
			invoiceView.Table_number = allOrderItems[0]["table_number"]
			invoiceView.Order_details = allOrderItems[0]["order_items"]
		}

		appLogger.Log.WithFields(logrus.Fields{
			"event":      "get_invoice_success",
			"time":       time.Now().Format(time.RFC3339),
			"invoice_id": invoiceId,
		}).Info("Successfully retrieved invoice item")
		c.JSON(http.StatusOK, invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "create_invoice_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while binding JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var order models.Order
		err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
		if err != nil {
			msg := "Order was not found"
			appLogger.Log.WithFields(logrus.Fields{
				"event":    "create_invoice_error",
				"time":     time.Now().Format(time.RFC3339),
				"order_id": invoice.Order_id,
				"error":    err,
			}).Error(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		status := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}

		invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()

		validationErr := validate.Struct(invoice)
		if validationErr != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "create_invoice_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": validationErr,
			}).Error("Validation error")
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		result, insertErr := invoiceCollection.InsertOne(ctx, invoice)
		if insertErr != nil {
			msg := "Invoice item was not created"
			appLogger.Log.WithFields(logrus.Fields{
				"event": "create_invoice_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": insertErr,
			}).Error(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		appLogger.Log.WithFields(logrus.Fields{
			"event":      "create_invoice_success",
			"time":       time.Now().Format(time.RFC3339),
			"invoice_id": invoice.Invoice_id,
		}).Info("Successfully created invoice item")
		c.JSON(http.StatusOK, result)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoice models.Invoice
		invoiceId := c.Param("invoice_id")

		if err := c.BindJSON(&invoice); err != nil {
			appLogger.Log.WithFields(logrus.Fields{
				"event": "update_invoice_error",
				"time":  time.Now().Format(time.RFC3339),
				"error": err,
			}).Error("Error occurred while binding JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"invoice_id": invoiceId}
		var updateObj primitive.D

		if invoice.Payment_method != nil {
			updateObj = append(updateObj, bson.E{"payment_method", invoice.Payment_method})
		}

		if invoice.Payment_status != nil {
			updateObj = append(updateObj, bson.E{"payment_status", invoice.Payment_status})
		}

		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", invoice.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		status := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}

		result, err := invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := "Invoice item update failed"
			appLogger.Log.WithFields(logrus.Fields{
				"event":      "update_invoice_error",
				"time":       time.Now().Format(time.RFC3339),
				"invoice_id": invoiceId,
				"error":      err,
			}).Error(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		appLogger.Log.WithFields(logrus.Fields{
			"event":      "update_invoice_success",
			"time":       time.Now().Format(time.RFC3339),
			"invoice_id": invoiceId,
		}).Info("Successfully updated invoice item")
		c.JSON(http.StatusOK, result)
	}
}
