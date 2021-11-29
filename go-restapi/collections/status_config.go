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
func getStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, _ := helper.Connect()

	prot, _ := db.Collection("protocols").Find(context.TODO(), bson.M{})
	stat, _ := db.Collection("statuses").Find(context.TODO(), bson.M{})
	it, _ := db.Collection("items").Find(context.TODO(), bson.M{})

	var statuses []models.Statuses
	var protocols []models.Protocols 
	var items []models.Items 

	var protocol models.Protocols
	for prot.Next(context.TODO()){
		err := prot.Decode(&protocol) 
		if err != nil {
			log.Fatal(err)
		}
		protocols = append(protocols, protocol)	
	}
	var status models.Statuses
	for stat.Next(context.TODO()){
		err := stat.Decode(&status) 
		if err != nil {
			log.Fatal(err)
		}
		statuses = append(statuses, status)	
	}
	var item models.Items
	for it.Next(context.TODO()){
		err := it.Decode(&item) 
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)	
	}

	for i := range statuses {
		filteredItems := getItemById(items, statuses[i].Id)
		protocols[i].Items = filteredItems
	}
	for i := range statuses {
		filteredProtocols := getProtocolById(protocols, statuses[i].Id)
		statuses[i].Protocol = filteredProtocols
	}
	json.NewEncoder(w).Encode(statuses) 

}

func getProtocolById(protocols []models.Protocols, id string) []models.Protocols {
	result := []models.Protocols {}
	for  i := range protocols {
		if protocols[i].Id == id {
			result = append(result, protocols[i])
		}
	}
	return result
}
func getItemById(items []models.Items, id string) []models.Items {
	result := []models.Items {}
	for  i := range items {
		if items[i].Id == id {
			result = append(result, items[i])
		}
	}
	return result
}

func getConfigs(w http.ResponseWriter, r *http.Request) {

	db, _ := helper.Connect()

	prot, _ := db.Collection("protocols").Find(context.TODO(), bson.M{})
	conf, _ := db.Collection("configs").Find(context.TODO(), bson.M{})
	it, _ := db.Collection("items").Find(context.TODO(), bson.M{})

	var configs []models.Config
	var protocols []models.Protocols 
	var items []models.Items 

	var protocol models.Protocols
	for prot.Next(context.TODO()){
		err := prot.Decode(&protocol) 
		if err != nil {
			log.Fatal(err)
		}
		protocols = append(protocols, protocol)	
	}
	var config models.Config
	for conf.Next(context.TODO()){
		err := conf.Decode(&config) 
		if err != nil {
			log.Fatal(err)
		}
		configs = append(configs, config)	
	}
	var item models.Items
	for it.Next(context.TODO()){
		err := it.Decode(&item) 
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)	
	}

	for i := range configs {
		filteredItems := getItemById(items, configs[i].Id)
		protocols[i].Items = filteredItems
	}
	for i := range configs {
		filteredProtocols := getProtocolById(protocols, configs[i].Id)
		configs[i].Protocol = filteredProtocols
	}
	json.NewEncoder(w).Encode(configs) 
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var status models.Statuses
	var params = mux.Vars(r)

	var id string = params["id"]
	filter := bson.M{"id": id}
	db, _ := helper.Connect()
	err := db.Collection("statuses").FindOne(context.TODO(), filter).Decode(&status)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(status)
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var config models.Config
	var params = mux.Vars(r)
	var id string = params["id"]
	filter := bson.M{"id": id}
	db, _ := helper.Connect()

	err := db.Collection("configs").FindOne(context.TODO(), filter).Decode(&config)

	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(config)
}


func updateStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]
	var status models.Statuses
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
	db, _ := helper.Connect()
	err := db.Collection("statuses").FindOneAndUpdate(context.TODO(), filter, update).Decode(&status)

	if err != nil {
		helper.GetError(err, w)
		return
	}
	status.Id = id

	json.NewEncoder(w).Encode(status)
}

func updateConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]

	var config models.Config
	filter := bson.M{"id": id}
	_ = json.NewDecoder(r.Body).Decode(&config)

	db, _ := helper.Connect()

	update := bson.D{
		{"$set", bson.D{
			{"id", config.Id},
			{"type", config.Type},
			{"name", config.Name},
			{"protocol", config.Protocol},
		}},
	}

	err := db.Collection("configs").FindOneAndUpdate(context.TODO(), filter, update).Decode(&config)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	config.Id = id

	json.NewEncoder(w).Encode(config)
}

func MainStatusConfigs() {
	r := helper.Routes
	r.HandleFunc("/descriptor/statuses", getStatuses).Methods("GET")
	r.HandleFunc("/descriptor/status/{id}", getStatus).Methods("GET")
	r.HandleFunc("/descriptor/status/{id}", updateStatuses).Methods("PUT")

}