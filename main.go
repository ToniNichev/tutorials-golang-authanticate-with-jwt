package main

import (
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	// In a real-world application, you would perform proper authentication here.
	// For the sake of this example, we'll just check if an API key is present.
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-Auth-Token")
		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Public routes (no authentication required)
	public := router.Group("/public")
	{
		public.GET("/info", func(c *gin.Context) {
			c.String(200, "Public information")
		})
		public.GET("/products", func(c *gin.Context) {
			c.String(200, "Public product list")
		})
	}

	// Private routes (require authentication)
	private := router.Group("/private")
	private.Use(AuthMiddleware())
	{
		private.GET("/data", func(c *gin.Context) {
			c.String(200, "Private data accessible after authentication")
		})
		private.POST("/create", func(c *gin.Context) {
			c.String(200, "Create a new resource")
		})
	}

	router.POST("query", AuthMiddleware(), validateSession, returnData)

	// Run the server on port 8080
	router.Run(":8080")
}
