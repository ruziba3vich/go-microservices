package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/products", func(c *gin.Context) {
		collection := client.Database("testdb").Collection("products")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cur, err := collection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(ctx)

		var products []bson.M
		if err = cur.All(ctx, &products); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, products)
	})

	router.Run(":8082")
}
