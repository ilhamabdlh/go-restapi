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
	// "github.com/gorilla/handlers"
	gomain "github.com/ilhamabdlh/go-restapi/main"
)
//Connection mongoDB with helper class
var collectionDescriptor = helper.ConnectDescriptorsDB()

func getDescriptors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created Book array
	var descriptors []models.Descriptor

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collectionDescriptor.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var descriptor models.Descriptor
		// & character returns the memory address of the following variable.
		err := cur.Decode(&descriptor) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		descriptors = append(descriptors, descriptor)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(descriptors) // encode similar to serialize process.
}

func getDescriptor(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var descriptor models.Descriptor
	// we get params with mux.
	var params = mux.Vars(r)
	

	// string to primitive.ObjectID
	var id string = params["id"]

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
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

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&descriptor)

	// insert our book model.
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
	//Get id from parameters
	var id string = params["id"]

	var descriptor models.Descriptor
	
	// Create filter
	filter := bson.M{"id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&descriptor)

	// prepare update model.
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
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	var id string = params["id"]

	// prepare filter.
	filter := bson.M{"id": id}

	deleteResult, err := collectionDescriptor.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}



func MainDescriptors() {
	//Init Router
	r := gomain.Route

  	// arrange our route
	r.HandleFunc("/descriptors/", getDescriptors).Methods("GET")
	r.HandleFunc("/descriptor/{id}", getDescriptor).Methods("GET")
	r.HandleFunc("/descriptor", createDescriptor).Methods("POST")
	r.HandleFunc("/descriptor/New/{id}", updateDescriptor).Methods("PUT")
	r.HandleFunc("/descriptor/{id}", deleteDescriptor).Methods("DELETE")

	// headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", ""})
	// methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	// origin := handlers.AllowedOrigins([]string{"*"})
	// http.ListenAndServe(":4001", handlers.CORS(headers, methods, origin)(r)) 

}