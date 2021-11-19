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
var collectionItem = helper.ConnectItemsDB()
// var collectionItemProtocol = helper.ConnectProtocolsDB()
func GetItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var items []models.Itemes
	cur, err := collectionItem.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var item models.Itemes
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

	var item models.Itemes
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

func createItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item models.Itemes

	_ = json.NewDecoder(r.Body).Decode(&item)

	resalt, err := collectionItem.InsertOne(context.TODO(), item)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(resalt)
}


func updateItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]
	var item models.Itemes
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

	json.NewEncoder(w).Encode(item)

}



func MainItems() {
	r := helper.Routes
	r.HandleFunc("/descriptor/items", GetItems).Methods("GET")
	r.HandleFunc("/descriptor/item/{id}", getItem).Methods("GET")
	r.HandleFunc("/descriptor/item/{id}", updateItems).Methods("PUT")

}