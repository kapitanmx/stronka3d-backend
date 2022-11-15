package controllers

import (
	"stronka3d-backend/models"
	"stronka3d-backend/db"

	"fmt"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	productCollection = *mongo.GetCollection(db.DB, "products")
) 

func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var product models.Product
		defer cancel()

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		defer cancel()
		product.ID = primitive.NewObjectID()
		result, err := productCollection.InsertOne(ctx, product)
		if err != nil {
			msg := fmt.Sprintf("Product not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error" : msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, product)
	}
}

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		productId := c.Param("_id")
		var product models.Product
		defer cancel()

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := productCollection.FindOne(ctx, bson.M{"product_id": productId})
		if err != nil {
			msg := Sprintf("Product not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err := strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"products", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},}}}
		}
		result, err := productCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupSrtage, projectStage})
		defer cancel() 
		if err != nil {
			msg := Sprintf("Error occured while listing products")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		var allProducts []bson.M
		if err = result.App(ctx, &allProducts); err != nil {
			log.Fatal(err)
		}
		
		c.JSON(http.StatusOK, allProducts[0])
	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		productId := c.Param("_id")
		var product models.Product
		defer cancel()

		result, err := productCollection.DeleteOne(ctx, bson.M{"product_id": productId})
		if err != nil {
			msg := Sprintf("Error: Unable to remove product")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}
