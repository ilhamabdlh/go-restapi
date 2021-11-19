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
var collectionItemProtocol = helper.ConnectProtocolsDB()
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
	// var di string = params["id"]
	

	var item models.Itemes
	// var protocol models.Protocols
	
	
	// filtered := bson.M{"items.id": di}
	filter := bson.M{"id": id}

	
	_ = json.NewDecoder(r.Body).Decode(&item)
	// _ = json.NewDecoder(r.Body).Decode(&protocol)

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

	// updating := bson.D{
	// 	{"$set", bson.D{
	// 		{"id", protocol.Id},
	// 		{"type", protocol.Type},
	// 		{"name", protocol.Name},
	// 		{"items", bson.D{
	// 			{"id", protocol.Items.Id},
	// 			{"type", protocol.Items.Type},
	// 			{"name", protocol.Items.Name},
	// 			{"priority", protocol.Items.Priority},
	// 			{"max", protocol.Items.Default.Max},
	// 			{"min", protocol.Items.Default.Max},
	// 			{"description", protocol.Items.Description},
	// 			{"ui", protocol.Items.Ui},
	// 			{"persist", protocol.Items.Persist},
	// 		}},
	// 	}},
	// }

	// update := bson.D{
	// 	{"$set", bson.D{
	// 		{"id", protocol.Id},
	// 		{"type", protocol.Type},
	// 		{"name", protocol.Name},
	// 		{"items", bson.D{
	// 			{"id", protocol.Items.Id},
	// 			{"type", protocol.Items.Type},
	// 			{"name", protocol.Items.Name},
	// 			{"priority", protocol.Items.Priority},
	// 			{"max", protocol.Items.Default.Max},
	// 			{"min", protocol.Items.Default.Max},
	// 			{"description", protocol.Items.Description},
	// 			{"ui", protocol.Items.Ui},
	// 			{"persist", protocol.Items.Persist},
	// 		}},
	// 	}},
	// }

	err := collectionItem.FindOneAndUpdate(context.TODO(), filter, update).Decode(&item)
	// err := collectionItemProtocol.FindOneAndUpdate(context.TODO(), filtered, updating).Decode(&protocol)


	if err != nil {
		helper.GetError(err, w)
		return
	}
	//  else if err != nil{
	// 	helper.GetError(err, w)
	// 	return
	// }
	// item.Id = id
	// protocol.Id = id

	json.NewEncoder(w).Encode(item)
	// json.NewEncoder(w).Encode(protocol)
	// json.NewEncoder(w).Encode(protocol)
	// func(p helper.Protocols) SetField(value string) bson.D{
	// 	return bson.D{"items.name", value}
	// }

	// x:= Protocols{}
	// var updateFields bson.D
	// updateFields = append(updateFields, x.SetField("updated"))






}



func MainItems() {
	r := helper.Routes
	r.HandleFunc("/descriptor/itemes", GetItems).Methods("GET")
	r.HandleFunc("/descriptor/item/{id}", getItem).Methods("GET")
	r.HandleFunc("/descriptor/items", createItems).Methods("POST")
	r.HandleFunc("/descriptor/item/{id}", updateItems).Methods("PUT")

}