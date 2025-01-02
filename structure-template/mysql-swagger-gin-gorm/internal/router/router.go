package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

// InitRouter initializes the Gin engine with necessary middleware and routes.
func InitRouter() *gin.Engine {

	// setupGin configures the Gin engine with necessary middleware.
	gin.SetMode(os.Getenv("RUN_MODE"))
	gin.ForceConsoleColor()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create a new route group for API version 1
	rg := r.Group("/api/v1")
	{
		// Register project-related routes
		RegisterProjectRoutes(rg)
	}
	return r
}
