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
	"github.com/gorilla/handlers"
)
//Connection mongoDB with helper class
var collectionConfig = helper.ConnectConfigsDB()

func getConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created Book array
	var configs []models.Config

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collectionConfig.Find(context.TODO(), bson.M{})

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
		var config models.Config
		// & character returns the memory address of the following variable.
		err := cur.Decode(&config) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		configs = append(configs, config)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(configs) // encode similar to serialize process.
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var config models.Config
	// we get params with mux.
	var params = mux.Vars(r)
	

	// string to primitive.ObjectID
	var id string = params["id"]

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"id": id}
	err := collectionConfig.FindOne(context.TODO(), filter).Decode(&config)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(config)
}


func createConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var config models.Config

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&config)

	// insert our book model.
	result, err := collectionConfig.InsertOne(context.TODO(), config)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "application/json")

	var params = mux.Vars(r)
	//Get id from parameters
	var id string = params["id"]

	var config models.Config
	
	// Create filter
	filter := bson.M{"id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&config)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"id", config.Id},
			{"type", config.Type},
			{"name", config.Name},
			{"protocol", config.Protocol},
		}},
	}

	err := collectionConfig.FindOneAndUpdate(context.TODO(), filter, update).Decode(&config)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	config.Id = id

	json.NewEncoder(w).Encode(config)
}


func deleteConfigs(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	var id string = params["id"]

	// prepare filter.
	filter := bson.M{"id": id}

	deleteResult, err := collectionConfig.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}



func MainConfigs() {
	//Init Router
	r := mux.NewRouter()

  	// arrange our route
	r.HandleFunc("/descriptor/configs", getConfigs).Methods("GET")
	r.HandleFunc("/descriptor/configs/{id}", getConfig).Methods("GET")
	r.HandleFunc("/descriptor/configs", createConfig).Methods("POST")
	r.HandleFunc("/descriptor/configs/{id}", updateConfigs).Methods("PUT")
	r.HandleFunc("/descriptor/configs/{id}", deleteConfigs).Methods("DELETE")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", ""})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origin := handlers.AllowedOrigins([]string{"*"})
	http.ListenAndServe(":4001", handlers.CORS(headers, methods, origin)(r)) 

}