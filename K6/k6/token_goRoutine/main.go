package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
	mu         sync.Mutex // Mutex for protecting shared resources
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

	// POST endpoint for storing tokens concurrently
	router.POST("/tokens", func(c *gin.Context) {
		// Check if the request has a token in the "Authorization" header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token not found in the header"})
			return
		}

		// Start a goroutine to store the token in MongoDB
		go func() {
			if err := storeToken(token); err != nil {
				// Handle the error, log it, etc.
			}
		}()

		c.JSON(http.StatusOK, gin.H{"message": "Token storage request received"})
	})

	// GET endpoint for retrieving stored tokens concurrently
	router.GET("/tokens", func(c *gin.Context) {
		// Start a goroutine to retrieve tokens from MongoDB
		go func() {
			tokens, err := retrieveTokens()
			if err != nil {
				// Handle the error, log it, etc.
				return
			}

			c.JSON(http.StatusOK, gin.H{"tokens": tokens})
		}()
	})

	router.Run(":6000")
}

// storeToken stores a token in MongoDB
func storeToken(token string) error {
	mu.Lock()
	defer mu.Unlock()

	_, err := collection.InsertOne(context.TODO(), map[string]interface{}{"token": token})
	return err
}

// retrieveTokens retrieves stored tokens from MongoDB
func retrieveTokens() ([]string, error) {
	mu.Lock()
	defer mu.Unlock()

	cursor, err := collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tokens []string

	for cursor.Next(context.TODO()) {
		var result map[string]interface{}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		tokens = append(tokens, result["token"].(string))
	}

	return tokens, nil
}
