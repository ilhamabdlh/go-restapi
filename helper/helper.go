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


func Connect() (*mongo.Database, error) {
    clientOptions := options.Client()
    clientOptions.ApplyURI("mongodb://localhost:27017")
    client, err := mongo.NewClient(clientOptions)
	ctx := context.Background()
    if err != nil {
        return nil, err
    }
	
    err = client.Connect(ctx)
    if err != nil {
        return nil, err
    }

	collection := client.Database("mongosDB")

    return collection, nil
}


// func ConnectConfigsDB() *mongo.Collection {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
// 	client, err := mongo.Connect(context.TODO(), clientOptions)

// 	if err != nil {
// 		log.Fatal(err)
// 	}


// 	// fmt.Println("Connected to MongoDB")

// 	collection := client.Database("mongosDB").Collection("configs")

// 	return collection
// }


// func ConnectProtocolsDB() *mongo.Collection {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
// 	client, err := mongo.Connect(context.TODO(), clientOptions)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// fmt.Println("Connected to MongoDB")

// 	collection := client.Database("mongosDB").Collection("protocols")

// 	return collection
// }

// func ConnectDescriptorsDB() *mongo.Collection {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
// 	client, err := mongo.Connect(context.TODO(), clientOptions)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Connected to MongoDB")

// 	collection := client.Database("mongosDB").Collection("descriptors")

// 	return collection
// }

// func ConnectStatusesDB() *mongo.Collection {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
// 	client, err := mongo.Connect(context.TODO(), clientOptions)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// fmt.Println("Connected to MongoDB")

// 	collection := client.Database("mongosDB").Collection("statuses")

// 	return collection
// }

// func ConnectItemsDB() *mongo.Collection {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
// 	client, err := mongo.Connect(context.TODO(), clientOptions)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// fmt.Println("Connected to MongoDB")

// 	collection := client.Database("mongosDB").Collection("items")

// 	return collection
// }

type ErrNotFound struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func GetError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var Response = ErrNotFound{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(Response)

	w.WriteHeader(Response.StatusCode)
	w.Write(message)
	fmt.Println(message)
}
