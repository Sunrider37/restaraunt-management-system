package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaraunt-management/database"
	"restaraunt-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(), 100*time.Second)
		
		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "error occured while listing orders"})
			return
		}
		var allOrders []bson.M
		if err := result.All(ctx,&allOrders); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "error while fetching order item"})
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var table models.Table
	var order models.Order

	if err := c.BindJSON(&order);err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
	}
	validationError := validate.Struct(order)
	if validationError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
	}
	if order.Table_id != nil{
		err := tableCollection.FindOne(ctx,bson.M{"table_id":order.Table_id}).Decode(&table)
		defer cancel()
		if err != nil{
			msg := fmt.Sprintf("message:Table was not found")
			c.JSON(http.StatusNotFound, gin.H{"error": msg})
			return
		}
	}
	order.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()
	result, insertErr := orderCollection.InsertOne(ctx,order)
	if insertErr != nil{
		msg := fmt.Sprintf("order item was not created")
		c.JSON(http.StatusNotFound, gin.H{"error": msg})
	}
	defer cancel()
	c.JSON(http.StatusAccepted,result)
}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var table models.Table
		var order models.Order

		var updateObj primitive.D

		orderId := c.Param("order_id")
		if err := c.BindJSON(&order);err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		}
		if order.Table_id != nil{
			err := menuCollection.FindOne(ctx,bson.M{"table_id":order.Table_id}).Decode(&table)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("message: Menu not found")
				c.JSON(http.StatusNotFound, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{"menu", order.Table_id})
		}

		order.Updated_At,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", order.Updated_At})
		upsert := true
		filter := bson.M{"order_id": orderId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderCollection.UpdateOne(ctx, filter, bson.M{"$set": updateObj}, &opt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
	}

	func OrderItemOrderCreator(order models.Order) string{
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		order.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()
		orderCollection.InsertOne(ctx,order)
		defer cancel()
		return order.Order_id
	}