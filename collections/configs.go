package collections

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ilhamabdlh/go-restapi/helper"
	"github.com/ilhamabdlh/go-restapi/models"
	"go.mongodb.org/mongo-driver/bson"
)

var collectionConfig = helper.ConnectConfigsDB()

func getConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var configs []models.Config
	cur, err := collectionConfig.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var config models.Config
		err := cur.Decode(&config)
		if err != nil {
			log.Fatal(err)
		}

		configs = append(configs, config)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)

	}

	json.NewEncoder(w).Encode(configs)
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
	r.HandleFunc("/descriptor/configs", createConfig).Methods("POST")
	r.HandleFunc("/descriptor/configs/{id}", updateConfigs).Methods("PUT")
	r.HandleFunc("/descriptor/configs/{id}", deleteConfigs).Methods("DELETE")

}
