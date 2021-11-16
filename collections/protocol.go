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
var collectionProtocol = helper.ConnectProtocolsDB()

func getProtocols(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	var protocols []models.Protocol

	cur, err := collectionProtocol.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var protocol models.Protocol
		
		err := cur.Decode(&protocol) 
		if err != nil {
			log.Fatal(err)
		}

		protocols = append(protocols, protocol)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(protocols) 
}

func getProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var protocol models.Protocol
	var params = mux.Vars(r)
	

	var id string = params["id"]

	filter := bson.M{"id": id}
	err := collectionProtocol.FindOne(context.TODO(), filter).Decode(&protocol)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(protocol)
}


func createProtocols(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var protocol models.Protocol

	_ = json.NewDecoder(r.Body).Decode(&protocol)

	result, err := collectionProtocol.InsertOne(context.TODO(), protocol)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]

	var protocol models.Protocol
	
	filter := bson.M{"id": id}

	_ = json.NewDecoder(r.Body).Decode(&protocol)

	update := bson.D{
		{"$set", bson.D{
			{"id", protocol.Id},
			{"type", protocol.Type},
			{"name", protocol.Name},
			{"items", protocol.Items},
		}},
	}

	err := collectionProtocol.FindOneAndUpdate(context.TODO(), filter, update).Decode(&protocol)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	protocol.Id = id

	json.NewEncoder(w).Encode(protocol)
}


func deleteProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	var id string = params["id"]

	filter := bson.M{"id": id}

	deleteResult, err := collectionProtocol.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}



func MainProtocols() {
	r := helper.Routes
	r.HandleFunc("/descriptor/protocols", getProtocols).Methods("GET")
	r.HandleFunc("/descriptor/protocol/{id}", getProtocol).Methods("GET")
	r.HandleFunc("/descriptor/protocols", createProtocols).Methods("POST")
	r.HandleFunc("/descriptor/protocol/{id}", updateProtocol).Methods("PUT")
	r.HandleFunc("/descriptor/protocol/{id}", deleteProtocol).Methods("DELETE")

}