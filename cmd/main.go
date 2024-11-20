package main

import (
	"fmt"
	"os"
	"porty-go/config"
	"porty-go/repositories"
	"porty-go/routes"
	"strings"

	_ "porty-go/docs"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Porty!!! API
// @version 1.0
// @description This is a Doc for a Porty!!! API.
// @BasePath /

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	service := os.Getenv("SERVICE")

	// Set the Swagger host dynamically
	service = strings.ToLower(service)
	var swaggerHost string
	if service == "local" {
		swaggerHost = fmt.Sprintf("localhost:%s", port)
	} else {
		swaggerHost = "porty.up.railway.app"
	}

	// Update Swagger documentation with the dynamic host
	doc := ginSwagger.URL(fmt.Sprintf("https://%s/swagger/doc.json", swaggerHost))

	client := config.LoadConfig()
	repositories.Init(client)

	r := gin.Default()

	// Customize CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8000"}, // Replace with your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register routes
	routes.SetupRouter(r)

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, doc))

	fmt.Println("Server is running at :" + port)
	r.Run(":" + port)
}
