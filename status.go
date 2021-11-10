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

var collectionStatus = helper.ConnectStatusesDB()
func getStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created Book array
	var statuses []models.Status

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collectionStatus.Find(context.TODO(), bson.M{})

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
		var status models.Status
		// & character returns the memory address of the following variable.
		err := cur.Decode(&status) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		statuses = append(statuses, status)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(statuses) // encode similar to serialize process.
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var status models.Status
	// we get params with mux.
	var params = mux.Vars(r)
	

	// string to primitive.ObjectID
	var id string = params["id"]

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"id": id}
	err := collectionStatus.FindOne(context.TODO(), filter).Decode(&status)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(status)
}

func createStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var status models.Status

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&status)

	// insert our book model.
	result, err := collectionStatus.InsertOne(context.TODO(), status)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}


func updateStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	//Get id from parameters
	var id string = params["id"]

	var status models.Status
	
	// Create filter
	filter := bson.M{"id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&status)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"id", status.Id},
			{"type", status.Type},
			{"name", status.Name},
			{"protocol", status.Protocol},
		}},
	}

	err := collectionStatus.FindOneAndUpdate(context.TODO(), filter, update).Decode(&status)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	status.Id = id

	json.NewEncoder(w).Encode(status)
}

func deleteStatuses(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	var id string = params["id"]

	// prepare filter.
	filter := bson.M{"id": id}

	deleteResult, err := collectionStatus.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

func MainStatus() {
	//Init Router
	r := routers.Route

  	// arrange our route
	r.HandleFunc("/descriptor/statuses", getStatuses).Methods("GET")
	r.HandleFunc("/descriptor/status/{id}", getStatus).Methods("GET")
	r.HandleFunc("/descriptor/status", createStatus).Methods("POST")
	r.HandleFunc("/descriptor/status/{id}", updateStatuses).Methods("PUT")
	r.HandleFunc("/descriptor/status/{id}", deleteStatuses).Methods("DELETE")

	// headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", ""})
	// methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	// origin := handlers.AllowedOrigins([]string{"*"})

	// log.Fatal(http.ListenAndServe(":4001", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r))) 

}