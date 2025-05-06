package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var DB *mongo.Database

func Connect() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://felipeprogramadordm:JFQx6WJeaeX6yd0A@clusterdev.hmufydl.mongodb.net/")
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

	DB = client.Database("4play-dev")
	log.Println("Mongo conectado")
}

