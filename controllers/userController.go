package controllers

import (
	"context"
	"net/http"
	"os"
	"porty-go/config"
	"porty-go/models"
	"porty-go/services"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	oauth2api "google.golang.org/api/oauth2/v2"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [post]
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	if (user.Email == "") || (user.Password == "") {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Email and password are required!",
		})
		return
	}

	result, err := services.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    result,
	})
}

// Login godoc
// @Summary Login a user
// @Description Login a user with the input payload
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/login [post]
func LoginUser(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	user, err := services.GetUserByEmail(loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid email or password",
		})
		return
	}

	if user.Password == "" && user.IsGoogle {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  "error",
			Message: "You have registered with Google, please login with Google. Set password to login with email",
		})
		return
	}

	// Compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid email or password",
		})
		return
	}

	// Generate a JWT token for the user
	token, err := services.GenerateToken(user.ID.Hex(), user.Email, user.FullName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to generate token",
		})
		return
	}

	result := models.DataLoginResponse{
		IdUser:      user.ID.Hex(),
		FullName:    user.FullName,
		Email:       user.Email,
		Token:       token,
		SetPassword: false,
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Status:  "success",
		Message: "Login successful",
		Data:    result,
	})
}

// GoogleLogin redirects the user to the Google login page
func GoogleLogin(c *gin.Context) {
	url := config.GoogleOAuthConfig().AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the callback from Google after the user has logged in
func GoogleCallback(c *gin.Context) {
	frontendURL := os.Getenv("FRONT_END_URL") // Replace with your frontend URL
	state := c.Query("state")
	if state != "state" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid state",
		})
		return
	}

	code := c.Query("code")
	redirectURLHome := frontendURL + "/"
	token, err := config.GoogleOAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, redirectURLHome)
	}

	client := config.GoogleOAuthConfig().Client(context.Background(), token)
	oauth2Service, err := oauth2api.New(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to create OAuth2 service",
		})
		return
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to get user info",
		})
		return
	}

	// Create or update the user
	tokenString, err := services.CreateOrUpdateOAuth(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: "Failed to create or update user",
		})
		return
	}
	// Redirect to the frontend with the token as a query parameter
	redirectURLSucces := frontendURL + "/auth/success?token=" + tokenString

	c.Redirect(http.StatusTemporaryRedirect, redirectURLSucces)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := services.GetUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// UpdateUser godoc
// @Summary Update a user by ID
// @Description Update a user by ID with the input payload
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "User"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	result, err := services.UpdateUserById(id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "User updated successfully",
		Data:    result,
	})
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Delete a user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	result, err := services.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "User deleted successfully",
		Data:    result,
	})
}

// VerifyUser godoc
// @Summary Verify a user by Email
// @Description Verify a user by Email
// @Tags users
// @Produce json
// @Param token query string true "Verification Token"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.ErrorResponse
// @Router /users/verify [get]
func VerifyEmail(c *gin.Context) {
	tokenString := c.Param("id")

	// Parse the token
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Status:  "error",
				Message: "Invalid token",
			})
			return
		}
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Bad request",
		})
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  "error",
			Message: "Invalid token",
		})
		return
	}

	// Get the email from the token claims
	email := claims.Subject

	user, err := services.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	if user.IsVerify {
		c.JSON(http.StatusOK, models.Response{
			Status:  "success",
			Message: "Email already verified",
			Data:    user,
		})
		return
	}

	// Check if the token is expired
	if claims.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  "error",
			Message: "Token has expired",
		})
		return
	}

	// Verify the user
	if _, err := services.VerifyUser(user.ID.Hex(), user); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Email verified successfully",
		Data:    user,
	})
}
