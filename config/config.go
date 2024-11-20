package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func LoadConfig() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	DB = client.Database("tedy")
	if DB == nil {
		log.Fatal("Failed to connect to database")
	}

	log.Println(`                                           __                 __
	    _____   ____     ____     ____     ___     _____   / /_   ___     ____/ /
	   / ___/  / __ \   / __ \   / __ \   / _ \   / ___/  / __/  / _ \   / __  / 
	  / /__   / /_/ /  / / / /  / / / /  /  __/  / /__   / /_   /  __/  / /_/ /  
	  \___/   \____/  /_/ /_/  /_/ /_/   \___/   \___/   \__/   \___/   \__,_/   `)

	return client
}
