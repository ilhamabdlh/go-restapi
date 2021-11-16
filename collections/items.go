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

var collectionItem = helper.ConnectItemsDB()

func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	var devices []models.Device

	cur, err := collectionProtocol.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var device models.Device
		
		err := cur.Decode(&device) 
		if err != nil {
			log.Fatal(err)
		}

		devices = append(devices, device)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(devices) 
}

func getItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item models.Device
	var params = mux.Vars(r)
	var id string = params["id"]
	filter := bson.M{"id": id}
	err := collectionItem.FindOne(context.TODO(), filter).Decode(&item)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item models.Device
	_ = json.NewDecoder(r.Body).Decode(&item)
	result, err := collectionItem.InsertOne(context.TODO(), item)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]

	var item models.Device
	filter := bson.M{"id": id}
	_ = json.NewDecoder(r.Body).Decode(&item)

	update := bson.D{
		{"$set", bson.D{
			{"id", item.Id},
			{"type", item.Type},
			{"name", item.Name},
			{"priority", item.Priority},
			{"max", item.Default.Max},
			{"min", item.Default.Max},
			{"description", item.Description},
			{"ui", item.Ui},
			{"persist", item.Persist},
		}},
	}

	err := collectionItem.FindOneAndUpdate(context.TODO(), filter, update).Decode(&item)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	item.Id = id

	json.NewEncoder(w).Encode(item)
}

func deleteItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]
	filter := bson.M{"id": id}

	deleteResult, err := collectionItem.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

func MainItems() {
	r := helper.Routes

	r.HandleFunc("/descriptor/items", getItems).Methods("GET")
	r.HandleFunc("/descriptor/item/{id}", getItem).Methods("GET")
	r.HandleFunc("/descriptor/items", createItem).Methods("POST")
	r.HandleFunc("/descriptor/item/{id}", updateItems).Methods("PUT")
	r.HandleFunc("/descriptor/item/{id}", deleteItems).Methods("DELETE")

}
