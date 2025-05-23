package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var DB *mongo.Database

func Connect() {

	connectString := os.Getenv("MONGO_URL")
	database := os.Getenv("MONGO_DATABASE")
	clientOptions := options.Client().ApplyURI(connectString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Erro ao conectar:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Mongo não respondeu:", err)
	}

	DB = client.Database(database)
	log.Println("Mongo conectado")
}
