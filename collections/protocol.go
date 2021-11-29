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

func getProtocols(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, _ := helper.Connect()

	prot, _ := db.Collection("protocols").Find(context.TODO(), bson.M{})
	it, _ := db.Collection("items").Find(context.TODO(), bson.M{})

	var protocols []models.Protocols
	var items []models.Items

	var protocol models.Protocols
	for prot.Next(context.TODO()) {
		err := prot.Decode(&protocol)
		if err != nil {
			log.Fatal(err)
		}
		protocols = append(protocols, protocol)
	}
	var item models.Items
	for it.Next(context.TODO()) {
		err := it.Decode(&item)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}
	for i := range protocols {
		filteredItems := getItemByIdProtocol(items, protocols[i].Id)
		protocols[i].Items = filteredItems
	}

	json.NewEncoder(w).Encode(protocols)
}
func getItemByIdProtocol(items []models.Items, id string) []models.Items {
	result := []models.Items{}
	for i := range items {
		if items[i].Id == id {
			result = append(result, items[i])
		}
	}
	return result
}

func getProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var protocol models.Protocols
	var params = mux.Vars(r)
	db, _ := helper.Connect()

	var id string = params["id"]

	filter := bson.M{"id": id}
	err := db.Collection("protocols").FindOne(context.TODO(), filter).Decode(&protocol)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(protocol)
}

func updateProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]
	db, _ := helper.Connect()
	var protocol models.Protocols
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

	err := db.Collection("protocols").FindOneAndUpdate(context.TODO(), filter, update).Decode(&protocol)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	protocol.Id = id

	json.NewEncoder(w).Encode(protocol)
}

func MainProtocols() {
	r := helper.Routes
	r.HandleFunc("/descriptor/protocols", getProtocols).Methods("GET")
	r.HandleFunc("/descriptor/protocol/{id}", getProtocol).Methods("GET")
	r.HandleFunc("/descriptor/protocol/{id}", updateProtocol).Methods("PUT")

}
