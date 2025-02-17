package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection
var integrationServiceCollection *mongo.Collection

func Init(client *mongo.Client) {
	userCollection = client.Database("tedy").Collection("users")
	integrationServiceCollection = client.Database("tedy").Collection("integrationService")
}
