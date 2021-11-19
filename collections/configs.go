package collections

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	// "fmt"
	

	"github.com/gorilla/mux"
	"github.com/ilhamabdlh/go-restapi/helper"
	"github.com/ilhamabdlh/go-restapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionConfig = helper.ConnectConfigsDB()

func getConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	lookupStage := bson.D{{"$lookup", bson.D{{"from", "protocols"}, {"localField", "id"}, {"foreignField", "id"}, {"as", "protocol"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$id"}, {"preserveNullAndEmptyArrays", false}}}}

	ShowLoadedCursor, err := collectionConfig.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
	if err !=nil{
		log.Fatal(err)
	}

	var ShowLoaded []bson.M
	if err = ShowLoadedCursor.All(ctx, &ShowLoaded); err!= nil{
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(ShowLoaded)
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var config models.Config
	var params = mux.Vars(r)
	var id string = params["id"]
	filter := bson.M{"id": id}

	err := collectionConfig.FindOne(context.TODO(), filter).Decode(&config)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(config)
}

func createConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var config models.Config
	_ = json.NewDecoder(r.Body).Decode(&config)
	result, err := collectionConfig.InsertOne(context.TODO(), config)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]

	var config models.Config
	filter := bson.M{"id": id}
	_ = json.NewDecoder(r.Body).Decode(&config)

	update := bson.D{
		{"$set", bson.D{
			{"id", config.Id},
			{"type", config.Type},
			{"name", config.Name},
			{"protocol", config.Protocol},
		}},
	}

	err := collectionConfig.FindOneAndUpdate(context.TODO(), filter, update).Decode(&config)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	config.Id = id

	json.NewEncoder(w).Encode(config)
}

func deleteConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]
	filter := bson.M{"id": id}

	deleteResult, err := collectionConfig.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

func MainConfigs() {
	r := helper.Routes

	r.HandleFunc("/descriptor/configs", getConfigs).Methods("GET")
	r.HandleFunc("/descriptor/configs/{id}", getConfig).Methods("GET")
	r.HandleFunc("/descriptor/configs/{id}", updateConfigs).Methods("PUT")

}
