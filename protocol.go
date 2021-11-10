package collections

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	
	"github.com/ilhamabdlh/go-restapi/helper"
	"github.com/ilhamabdlh/go-restapi/models"
	"github.com/ilhamabdlh/go-restapi/routers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	// "github.com/gorilla/handlers"
)
//Connection mongoDB with helper class
var collectionProtocol = helper.ConnectProtocolsDB()

func getProtocols(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created Book array
	var protocols []models.Protocol

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collectionProtocol.Find(context.TODO(), bson.M{})

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
		var protocol models.Protocol
		// & character returns the memory address of the following variable.
		err := cur.Decode(&protocol) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		protocols = append(protocols, protocol)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(protocols) // encode similar to serialize process.
}

func getProtocol(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var protocol models.Protocol
	// we get params with mux.
	var params = mux.Vars(r)
	

	// string to primitive.ObjectID
	var id string = params["id"]

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
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

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&protocol)

	// insert our book model.
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
	//Get id from parameters
	var id string = params["id"]

	var protocol models.Protocol
	
	// Create filter
	filter := bson.M{"id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&protocol)

	// prepare update model.
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
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	var id string = params["id"]

	// prepare filter.
	filter := bson.M{"id": id}

	deleteResult, err := collectionProtocol.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}



func MainProtocols() {
	//Init Router
	r := routers.Route

  	// arrange our route
	r.HandleFunc("/descriptor/protocols", getProtocols).Methods("GET")
	r.HandleFunc("/descriptor/protocol/{id}", getProtocol).Methods("GET")
	r.HandleFunc("/descriptor/protocols", createProtocols).Methods("POST")
	r.HandleFunc("/descriptor/protocol/{id}", updateProtocol).Methods("PUT")
	r.HandleFunc("/descriptor/protocol/{id}", deleteProtocol).Methods("DELETE")

	// headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", ""})
	// methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	// origin := handlers.AllowedOrigins([]string{"*"})
	// http.ListenAndServe(":4001", handlers.CORS(headers, methods, origin)(r)) 

}