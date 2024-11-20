package routes

import (
	"porty-go/controllers"

	"github.com/gin-gonic/gin"
)

// UserRoutes defines the user-related routes
func UserRoutes(r *gin.Engine) {
	r.GET("/users/:id", controllers.GetUser)
	r.GET("/users/verify/:id", controllers.VerifyEmail)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)

	// Google OAuth routes
	r.POST("/auth/register", controllers.RegisterUser)
	r.POST("/auth/login", controllers.LoginUser)
	r.GET("/auth/google/login", controllers.GoogleLogin)
	r.GET("/auth/google/callback", controllers.GoogleCallback)

	// http://localhost:8000/auth/google/callback
}
