package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/hello", func(c *gin.Context) {
		// Parse the JSON request body
		var requestBody map[string]interface{}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Modify the JSON data
		requestBody["message"] = "Modified: " + requestBody["message"].(string)

		// Send the modified JSON as a response
		c.JSON(http.StatusOK, requestBody)
	})

	router.Run(":6000")
}
