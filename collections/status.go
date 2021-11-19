package collections

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	
	"github.com/ilhamabdlh/go-restapi/helper"
	"github.com/ilhamabdlh/go-restapi/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	
)
var collectionStatus = helper.ConnectStatusesDB()
func getStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "protocols"}, {"localField", "id"}, {"foreignField", "id"}, {"as", "protocol"}}}}
	lookupStageTwo := bson.D{{"$lookup", bson.D{{"from", "items"}, {"localField", "id"}, {"foreignField", "id"}, {"as", "protocol.items"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$protocol"}, {"preserveNullAndEmptyArrays", false}}}}

	showLoadedCursor, err := collectionStatus.Aggregate(ctx, mongo.Pipeline{lookupStageTwo, unwindStage, lookupStage})
	if err !=nil{
		log.Fatal(err)
	}

	var ShowLoaded []bson.M
	if err = showLoadedCursor.All(ctx, &ShowLoaded); err!= nil{
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(ShowLoaded) 
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var status models.Statuses
	var params = mux.Vars(r)
	

	var id string = params["id"]

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

	var status models.Statuses

	_ = json.NewDecoder(r.Body).Decode(&status)

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

	err := collectionStatus.FindOneAndUpdate(context.TODO(), filter, update).Decode(&status)


	if err != nil {
		helper.GetError(err, w)
		return
	}
	status.Id = id

	json.NewEncoder(w).Encode(status)
}

func deleteStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	var id string = params["id"]

	filter := bson.M{"id": id}

	deleteResult, err := collectionStatus.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

func MainStatus() {
	r := helper.Routes
	r.HandleFunc("/descriptor/statuses", getStatuses).Methods("GET")
	r.HandleFunc("/descriptor/status/{id}", getStatus).Methods("GET")
	r.HandleFunc("/descriptor/status/{id}", updateStatuses).Methods("PUT")

}