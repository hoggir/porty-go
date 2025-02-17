package main

import (
	"fmt"
	"net/http"
	"os"
	"porty-go/config"
	"porty-go/models"
	"porty-go/repositories"
	"porty-go/routes"
	"strings"
	"time"

	_ "porty-go/docs"

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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use format: "Bearer {your_token}"
func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	service := os.Getenv("WEB_SERVICE")

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
	if service == "local" {
		doc = ginSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", swaggerHost))
	}

	client := config.LoadConfig()
	repositories.Init(client)

	r := gin.Default()

	// Customize CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.AllowedOrigins, // Replace with your frontend URL
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

	// Add a "Not Found" route
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, models.Response{
			Status:  "error",
			Message: "Route not found",
		})
	})

	// Print the Swagger URL to the console
	swaggerURL := fmt.Sprintf("http://%s/swagger/index.html", swaggerHost)
	fmt.Printf("Swagger documentation available at: %s\n", swaggerURL)

	fmt.Println("Server is running at :" + port)
	r.Run(":" + port)
}
