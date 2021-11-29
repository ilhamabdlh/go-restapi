package collections

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	
	"github.com/ilhamabdlh/go-restapi/helper"
	"github.com/ilhamabdlh/go-restapi/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func getDescriptors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, _ := helper.Connect()

	cur, err := db.Collection("descriptors").Find(context.TODO(), bson.M{})
	conf, _ := db.Collection("configs").Find(context.TODO(), bson.M{})
	prot, _ := db.Collection("protocols").Find(context.TODO(), bson.M{})
	stat, _ := db.Collection("statuses").Find(context.TODO(), bson.M{})
	it, _ := db.Collection("items").Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())
	
	var descriptors []models.Descriptor
	var responses []models.Response
	var configs []models.Config
	var statuses []models.Statuses
	var protocols []models.Protocols 
	var items []models.Items 

	for cur.Next(context.TODO()) {
		var descriptor models.Descriptor
		err := cur.Decode(&descriptor) 
		if err != nil {
			log.Fatal(err)
		}
		descriptors = append(descriptors, descriptor)
	}
	var config models.Config
	for conf.Next(context.TODO()){
		err := conf.Decode(&config) 
		if err != nil {
			log.Fatal(err)
		}
		configs = append(configs, config)	
	}
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

	for i := range descriptors {
		filteredItems := getItemByIdFromDescriptor(items, descriptors[i].Id)
		protocols[i].Items = filteredItems
		filteredItemsTwo := getItemByIdFromDescriptor(items, descriptors[i].Id)
		protocols[i].Items = filteredItemsTwo
	}

	for i := range descriptors {
		filteredProtocols := getProtocolByIdFromDescriptor(protocols, descriptors[i].Id)
		configs[i].Protocol = filteredProtocols
		filteredProtocolsTwo := getProtocolByIdFromDescriptor(protocols, descriptors[i].Id)
		statuses[i].Protocol = filteredProtocolsTwo
	}

	for i := range descriptors {
		filteredConfigs := getConfigByIdFromDescriptor(configs, descriptors[i].Id)
		descriptors[i].Configs = filteredConfigs
		filteredStatuses := getStatusByIdFromDescriptor(statuses, descriptors[i].Id)
		descriptors[i].Status = filteredStatuses
	}
	var msg int
	var success bool
	if err != nil {
		msg = 400
		success = false
	} else {
		msg = 200
		success = true
	}

	var response models.Response
	response.Data = descriptors
	response.Status = strconv.Itoa(msg)
	response.Success = success
	response.Msg = http.StatusText(msg)
	responses = append(responses, response)
	
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(responses)
}

func getConfigByIdFromDescriptor(configs []models.Config, id string) []models.Config {
	result := []models.Config {}
	for  i := range configs {
		if configs[i].Id == id {
			result = append(result, configs[i])
		}
	}
	return result
}

func getStatusByIdFromDescriptor(statuses []models.Statuses, id string) []models.Statuses {
	result := []models.Statuses {}
	for  i := range statuses {
		if statuses[i].Id == id {
			result = append(result, statuses[i])
		}
	}
	return result
}

func getProtocolByIdFromDescriptor(protocols []models.Protocols, id string) []models.Protocols {
	result := []models.Protocols {}
	for  i := range protocols {
		if protocols[i].Id == id {
			result = append(result, protocols[i])
		}
	}
	return result
}
func getItemByIdFromDescriptor(items []models.Items, id string) []models.Items {
	result := []models.Items {}
	for  i := range items {
		if items[i].Id == id {
			result = append(result, items[i])
		}
	}
	return result
}

func getDescriptor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var descriptors []models.Descriptor
	var descriptor models.Descriptor
	var params = mux.Vars(r)
	
	var id string = params["id"]

	filter := bson.M{"id": id}
	db, _ := helper.Connect()
	err := db.Collection("descriptors").FindOne(context.TODO(), filter).Decode(&descriptor)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	descriptors = append(descriptors, descriptor)
	var status int
	if err != nil {
		status = 400
	} else {
		status = 200
	}

	var response models.Response
	response.Data = descriptors
	response.Status = strconv.Itoa(status)
	response.Success = true
	response.Msg = http.StatusText(status)

	json.NewEncoder(w).Encode(response)
}

func createDescriptor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var descriptor models.Descriptor
	var config models.Config
	var status models.Statuses
	var protocolOne models.Protocols
	var protocolTwo models.Protocols
	var itemOne models.Items
	var itemTwo models.Items
	db, _ := helper.Connect()

	_ = json.NewDecoder(r.Body).Decode(&descriptor)

	itemOne = descriptor.Configs[0].Protocol[0].Items[0]
	itemTwo = descriptor.Status[0].Protocol[0].Items[0]
	_ = json.NewDecoder(r.Body).Decode(&itemOne)
	_ = json.NewDecoder(r.Body).Decode(&itemTwo)
	
	descriptor.Configs[0].Protocol[0].Items = make([]models.Items, 0)
	descriptor.Status[0].Protocol[0].Items = make([]models.Items, 0)
	
	protocolOne = descriptor.Configs[0].Protocol[0]
	protocolTwo = descriptor.Status[0].Protocol[0]

	_ = json.NewDecoder(r.Body).Decode(&protocolOne)
	_ = json.NewDecoder(r.Body).Decode(&protocolTwo)

	descriptor.Status[0].Protocol = make([]models.Protocols, 0)
	descriptor.Configs[0].Protocol = make([]models.Protocols, 0)

	status = descriptor.Status[0]
	config = descriptor.Configs[0]
	_ = json.NewDecoder(r.Body).Decode(&status)
	_ = json.NewDecoder(r.Body).Decode(&config)

	conf, errr := db.Collection("configs").InsertOne(context.TODO(), config)
	stat, _ := db.Collection("statuses").InsertOne(context.TODO(), status)
	ptOne, _ := db.Collection("protocols").InsertOne(context.TODO(), protocolOne)
	ptTwo, _ := db.Collection("protocols").InsertOne(context.TODO(), protocolTwo)
	itOne, _ := db.Collection("items").InsertOne(context.TODO(), itemOne)
	itTwo, _ := db.Collection("items").InsertOne(context.TODO(), itemTwo)

	descriptor.Status = make([]models.Statuses, 0)
	descriptor.Configs = make([]models.Config, 0)

	_ = json.NewDecoder(r.Body).Decode(&descriptor)
	result, _ := db.Collection("descriptors").InsertOne(context.TODO(), descriptor)
	
	if errr != nil {
		helper.GetError(errr, w)
		return
	}
	json.NewEncoder(w).Encode(result)
	json.NewEncoder(w).Encode(stat)
	json.NewEncoder(w).Encode(conf)
	json.NewEncoder(w).Encode(ptOne)
	json.NewEncoder(w).Encode(ptTwo)
	json.NewEncoder(w).Encode(itOne)
	json.NewEncoder(w).Encode(itTwo)
}

func updateDescriptor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]

	var descriptor models.Descriptor
	
	filter := bson.M{"id": id}

	_ = json.NewDecoder(r.Body).Decode(&descriptor)
	update := bson.D{
		{"$set", bson.D{
			{"id", descriptor.Id},
			{"type", descriptor.Type},
			{"name", descriptor.Name},
			{"version", descriptor.Version},
			{"modules", descriptor.Modules},
			{"configs", descriptor.Configs},
			{"status", descriptor.Status},
		}},
	}
	db, _ := helper.Connect()

	err := db.Collection("descriptors").FindOneAndUpdate(context.TODO(), filter, update).Decode(&descriptor)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	descriptor.Id = id
	json.NewEncoder(w).Encode(descriptor)
}

func MainDescriptors() {
	r := helper.Routes

	r.HandleFunc("/descriptors/", getDescriptors).Methods("GET")
	r.HandleFunc("/descriptor/{id}", getDescriptor).Methods("GET")
	r.HandleFunc("/descriptor/New", createDescriptor).Methods("POST")
	r.HandleFunc("/descriptor/{id}", updateDescriptor).Methods("PUT")

}