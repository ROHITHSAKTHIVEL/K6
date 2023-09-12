package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
)

func init() {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(context.TODO(), clientOptions)

	// Check the connection
	err := client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	// Access the database and collection
	database := client.Database("mydb")        // Replace "mydb" with your database name
	collection = database.Collection("tokens") // Replace "tokens" with your collection name
}

func main() {
	router := gin.Default()

	// router.POST("/createuser", func(ctx *gin.Context) {

	// })

	router.POST("/tokens", func(c *gin.Context) {
		// Check if the request has a token in the "Authorization" header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found in the header"})
			return
		}

		// Store the token in MongoDB
		_, err := collection.InsertOne(context.TODO(), map[string]interface{}{"token": token})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Token stored successfully"})
	})

	router.GET("/tokens", func(c *gin.Context) {
		// Query MongoDB to retrieve stored tokens
		cursor, err := collection.Find(context.TODO(), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(context.TODO())

		var tokens []string

		for cursor.Next(context.TODO()) {
			var result map[string]interface{}
			if err := cursor.Decode(&result); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			tokens = append(tokens, result["token"].(string))
		}

		// Check if there are no stored tokens
		if len(tokens) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No tokens found"})
			return
		}

		// Send the retrieved tokens as a response
		c.JSON(http.StatusOK, gin.H{"tokens": tokens})
	})

	router.Run(":6000")
}
