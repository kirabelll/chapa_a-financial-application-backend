package main

import (
	"os"
	"net/http"

	"github.com/gin-gonic/gin"

	"chapa_a/internal/database"

	"go.uber.org/zap"
	"github.com/joho/godotenv"
)

var logger *zap.Logger

func initLogger() {
	var err error
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
}

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize logger
	initLogger()
	defer logger.Sync()

	logger.Info("Starting chapa_a server")




	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Info("Database connected successfully")
	_ = db // Use db in your handlers

	// Get host from environment
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	// Get HTTP port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := host + ":" + port

	logger.Info("Starting HTTP server", zap.String("address", addr))

	// Create Gin router
	r := gin.Default()


	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to chapa_a!",
		})
	})

	// Start server
	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
