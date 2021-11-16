package collections

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	
	"github.com/ilhamabdlh/go-restapi/helper"
	"github.com/ilhamabdlh/go-restapi/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	
)
var collectionStatus = helper.ConnectStatusesDB()
// var collectionStatusProtocol = helper.ConnectStatusesDB()
func getStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var statuses []models.Status
	cur, err := collectionStatus.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var status models.Status
		err := cur.Decode(&status) 
		if err != nil {
			log.Fatal(err)
		}

		statuses = append(statuses, status)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(statuses) 
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var status models.Status
	var params = mux.Vars(r)
	

	var id string = params["id"]

	filter := bson.M{"id": id}
	err := collectionStatus.FindOne(context.TODO(), filter).Decode(&status)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(status)
}

func createStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var status models.Status

	_ = json.NewDecoder(r.Body).Decode(&status)

	result, err := collectionStatus.InsertOne(context.TODO(), status)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}


func updateStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]

	var status models.Status
	
	filter := bson.M{"id": id}

	_ = json.NewDecoder(r.Body).Decode(&status)

	update := bson.D{
		{"$set", bson.D{
			{"id", status.Id},
			{"type", status.Type},
			{"name", status.Name},
			{"protocol", status.Protocol},
		}},
	}

	err := collectionStatus.FindOneAndUpdate(context.TODO(), filter, update).Decode(&status)


	if err != nil {
		helper.GetError(err, w)
		return
	}
	status.Id = id

	json.NewEncoder(w).Encode(status)
}

func deleteStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	var id string = params["id"]

	filter := bson.M{"id": id}

	deleteResult, err := collectionStatus.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

func MainStatus() {
	r := helper.Routes
	r.HandleFunc("/descriptor/statuses", getStatuses).Methods("GET")
	r.HandleFunc("/descriptor/status/{id}", getStatus).Methods("GET")
	r.HandleFunc("/descriptor/status", createStatus).Methods("POST")
	r.HandleFunc("/descriptor/status/{id}", updateStatuses).Methods("PUT")
	r.HandleFunc("/descriptor/status/{id}", deleteStatuses).Methods("DELETE")

}