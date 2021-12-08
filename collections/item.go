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

func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var items []models.Items
	db, _ := helper.Connect()
	cur, err := db.Collection("items").Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var item models.Items
		err := cur.Decode(&item) 
		if err != nil {
			log.Fatal(err)
		}

		items = append(items, item)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item models.Items
	var params = mux.Vars(r)

	db, _ := helper.Connect()
	var id string = params["id"]

	filter := bson.M{"id": id}
	err := db.Collection("items").FindOne(context.TODO(), filter).Decode(&item)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(item)
} 

func updateItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]
	var item models.Items
	filter := bson.M{"id": id}
	db, _ := helper.Connect()
	_ = json.NewDecoder(r.Body).Decode(&item)

	update := bson.D{
		{"$set", bson.D{
			{"id", item.Id},
			{"type", item.Type},
			{"name", item.Name},
			{"priority", item.Priority},
			{"default", bson.D{
			{"max", item.Default.Max},
			{"min", item.Default.Min},
			}},
			{"description", item.Description},
			{"ui", item.Ui},
			{"persist", item.Persist},
		}},
	}

	err := db.Collection("items").FindOneAndUpdate(context.TODO(), filter, update).Decode(&item)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func MainItems() {
	r := helper.Routes

	r.HandleFunc("/descriptor/items", getItems).Methods("GET")
	r.HandleFunc("/descriptor/item/{id}", getItem).Methods("GET")
	r.HandleFunc("/descriptor/item/{id}", updateItems).Methods("PUT")

	r.HandleFunc("/descriptor/configs", getConfigs).Methods("GET")
	r.HandleFunc("/descriptor/configs/{id}", getConfig).Methods("GET")
	r.HandleFunc("/descriptor/configs/{id}", updateConfigs).Methods("PUT")
}