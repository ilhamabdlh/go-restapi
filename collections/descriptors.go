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
var collectionDescriptor = helper.ConnectDescriptorsDB()


func getDescriptors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var descriptors []models.Descriptor

	cur, err := collectionDescriptor.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

	var descriptor models.Descriptor
	
	err := cur.Decode(&descriptor) 
	if err != nil {
		log.Fatal(err)
	}

	descriptors = append(descriptors, descriptor)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(descriptors)
}

func getDescriptor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var descriptor models.Descriptor
	var params = mux.Vars(r)
	
	var id string = params["id"]

	filter := bson.M{"id": id}
	err := collectionDescriptor.FindOne(context.TODO(), filter).Decode(&descriptor)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(descriptor)
}


func createDescriptor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var descriptor models.Descriptor

	_ = json.NewDecoder(r.Body).Decode(&descriptor)

	result, err := collectionDescriptor.InsertOne(context.TODO(), descriptor)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
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

	err := collectionDescriptor.FindOneAndUpdate(context.TODO(), filter, update).Decode(&descriptor)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	descriptor.Id = id

	json.NewEncoder(w).Encode(descriptor)
}


func deleteDescriptor(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
var params = mux.Vars(r)

var id string = params["id"]
filter := bson.M{"id": id}

deleteResult, err := collectionDescriptor.DeleteOne(context.TODO(), filter)

if err != nil {
	helper.GetError(err, w)
	return
}

json.NewEncoder(w).Encode(deleteResult)
}



func MainDescriptors() {
	r := helper.Routes

	r.HandleFunc("/descriptors/", getDescriptors).Methods("GET")
	r.HandleFunc("/descriptor/{id}", getDescriptor).Methods("GET")
	r.HandleFunc("/descriptor/New", createDescriptor).Methods("POST")
	r.HandleFunc("/descriptor/{id}", updateDescriptor).Methods("PUT")
	r.HandleFunc("/descriptor/{id}", deleteDescriptor).Methods("DELETE")

	// headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", ""})
	// methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	// origin := handlers.AllowedOrigins([]string{"*"})
	// http.ListenAndServe(":4001", handlers.CORS(headers, methods, origin)(r)) 

}