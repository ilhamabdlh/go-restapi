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
