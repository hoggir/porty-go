package services

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"porty-go/models"
	"porty-go/repositories"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	oauth2api "google.golang.org/api/oauth2/v2"
	"gopkg.in/gomail.v2"
)

func RegisterUser(user models.User) (*mongo.InsertOneResult, error) {
	user.ID = primitive.NewObjectID()
	user.IsGoogle = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = nil

	// Encrypt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error encrypting password:", err)
		return nil, err
	}
	user.Password = string(hashedPassword)

	// Check if user already exists
	existingUser, err := GetUserByEmail(user.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println("Error checking if user exists:", err)
		return nil, err
	}

	if existingUser.Email != "" {
		log.Println("User already exists:", user.Email)
		return nil, errors.New("user already exists")
	}

	// Generate the verification token
	token, err := GenerateVerificationToken(user.Email)
	if err != nil {
		log.Println("Error generating verification token:", err)
		return nil, errors.New("failed to generate verification token")
	}
	verificationLink := "http://localhost:3000/users/verify?token=" + token

	// Send welcome email
	if err := SendWelcomeEmail(user.Email, verificationLink); err != nil {
		log.Println("Error sending welcome email:", err)
		return nil, errors.New("failed to send welcome email")
	}

	result, err := repositories.CreateUser(user)
	if err != nil {
		log.Println("Error creating user:", err)
		return nil, err
	}

	return result, nil
}

func GetUser(id string) (models.User, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	return repositories.GetUserById(objID)
}

func GetUserByEmail(email string) (models.User, error) {
	return repositories.GetUserByEmail(email)
}

func UpdateUserById(id string, user models.User) (*mongo.UpdateResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	now := time.Now()
	user.UpdatedAt = &now
	return repositories.UpdateUserById(objID, user)
}

func DeleteUser(id string) (*mongo.DeleteResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	return repositories.DeleteUser(objID)
}

func VerifyUser(id string, user models.User) (*mongo.UpdateResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	user.IsVerify = true
	now := time.Now()
	user.VerifyAt = &now
	return repositories.UpdateUserById(objID, user)
}

func GenerateVerificationToken(email string) (string, error) {
	// Define the token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims, which includes the email and expiry time
	claims := &jwt.StandardClaims{
		Subject:   email,
		ExpiresAt: expirationTime.Unix(),
	}

	// Create the token using your secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println("Failed to generate token:", err)
		return "", err
	}

	// Parse the token
	decode := &jwt.StandardClaims{}
	dctoken, err := jwt.ParseWithClaims(tokenString, decode, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", err
		}
		return "", err
	}

	if !dctoken.Valid {
		return "", err
	}

	// Get the email from the token claims
	// email := decode.Subject
	fmt.Println("Email:", decode.Subject)

	return tokenString, nil
}

func SendWelcomeEmail(to string, verificationLink string) error {
	// Parse the HTML template
	tmpl, err := template.ParseFiles("templates/welcome_email.html")
	if err != nil {
		fmt.Println("Failed to parse template:", err)
		return nil
	}

	// Data to pass to the template
	data := struct {
		Name             string
		VerificationLink string
	}{
		Name:             to,
		VerificationLink: verificationLink,
	}

	// Execute the template with data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		fmt.Println("Failed to execute template:", err)
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_USERNAME"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Welcome to Porty!!!")
	m.SetBody("text/plain", "Thank you for registering with Porty!!!")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"))

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send email:", err)
		return nil
	}
	return nil
}

func CreateOrUpdateOAuth(userInfo *oauth2api.Userinfo) (string, error) {
	user, err := GetUserByEmail(userInfo.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println("Error checking if user exists:", err)
		return "", err
	}

	isCreted := user.ID.IsZero()

	if !isCreted {
		now := time.Now()
		if !user.IsVerify {
			user.IsVerify = true
		}
		if user.VerifyAt == nil {
			user.VerifyAt = &now
		}
		if !user.IsGoogle {
			user.IsGoogle = true
		}
		user.LastLogin = &now
		_, err := UpdateUserById(user.ID.Hex(), user)
		if err != nil {
			log.Println("Error updating user:", err)
			return "", err
		}
	}

	if isCreted {
		now := time.Now()
		user = models.User{
			ID:        primitive.NewObjectID(),
			FullName:  userInfo.Name,
			Email:     userInfo.Email,
			Password:  "",
			IsVerify:  true,
			VerifyAt:  &now,
			IsGoogle:  true,
			LastLogin: &now,
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		}

		_, err := repositories.CreateUser(user)
		if err != nil {
			log.Println("Error creating user:", err)
			return "", err
		}
	}

	token, err := GenerateToken(user.ID.Hex(), user.Email, user.FullName)
	if err != nil {
		log.Println("Error generating token:", err)
		return "", err
	}
	return token, nil
}

// CustomClaims defines the custom claims for the JWT token
type CustomClaims struct {
	FullName string `json:"full_name"`
	jwt.StandardClaims
}

// generateToken generates a JWT token with the user's ID, email, and full name
func GenerateToken(userID, email, fullName string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &CustomClaims{
		FullName: fullName,
		StandardClaims: jwt.StandardClaims{
			Subject:   userID,
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "porty-go",
			Audience:  email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
