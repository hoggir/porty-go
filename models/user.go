package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" swaggerignore:"true"`
	FullName  string             `bson:"fullName"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	LastLogin *time.Time         `bson:"lastLogin" json:"lastLogin" swaggerignore:"true"`
	IsGoogle  bool               `bson:"isGoogle" json:"isGoogle" swaggerignore:"true"`
	IsVerify  bool               `bson:"isVerify"`
	VerifyAt  *time.Time         `bson:"VerifyAt" json:"VerifyAt" swaggerignore:"true"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt" swaggerignore:"true"`
	UpdatedAt *time.Time         `bson:"updatedAt" json:"updatedAt" swaggerignore:"true"`
}
