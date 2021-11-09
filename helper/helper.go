package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectConfigsDB() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection := client.Database("mongosDB").Collection("configs")

	return collection
}

func ConnectProtocolsDB() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection := client.Database("mongosDB").Collection("protocols")

	return collection
}

func ConnectDescriptorsDB() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection := client.Database("mongosDB").Collection("descriptors")

	return collection
}

func ConnectStatusesDB() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection := client.Database("mongosDB").Collection("statuses")

	return collection
}


type errNotFound struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}


func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = errNotFound{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
type Configuration struct {
	Port             string
	ConnectionString string
}