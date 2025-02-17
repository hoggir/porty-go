package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type IntegrationService struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" swaggerignore:"true"`
	ServiceName string             `json:"serviceName,omitempty"`
	ServiceUrl  string             `json:"serviceUrl"`
	Token       string             `json:"token"`
	UserName    string             `json:"userName"`
	Password    string             `json:"password"`
	Model       string             `json:"model"`
}
