package repositories

import (
	"context"
	"porty-go/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetIntegrationServiceByName(name string) (models.IntegrationService, error) {
	var service models.IntegrationService
	err := integrationServiceCollection.FindOne(context.Background(), bson.M{"serviceName": name}).Decode(&service)
	return service, err
}
