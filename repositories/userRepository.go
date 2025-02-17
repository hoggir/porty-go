package repositories

import (
	"context"
	"porty-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(user models.User) (*mongo.InsertOneResult, error) {
	return userCollection.InsertOne(context.Background(), user)
}

func GetUserById(id primitive.ObjectID) (models.User, error) {
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	return user, err
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	return user, err
}

func UpdateUserById(id primitive.ObjectID, user models.User) (*mongo.UpdateResult, error) {
	return userCollection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": user})
}

func DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return userCollection.DeleteOne(context.Background(), bson.M{"_id": id})
}
