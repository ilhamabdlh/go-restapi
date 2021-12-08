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

func getProtocols(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, _ := helper.Connect()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	lookupStageTwo := bson.D{{"$lookup", bson.D{{"from", "items"}, {"localField", "id"}, {"foreignField", "id"}, {"as", "items"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$id"}, {"preserveNullAndEmptyArrays", false}}}}

	showLoadedCursor, err := db.Collection("protocols").Aggregate(ctx, mongo.Pipeline{lookupStageTwo, unwindStage})
	if err !=nil{
		log.Fatal(err)
	}

	var showLoaded []bson.M
	if err = showLoadedCursor.All(ctx, &showLoaded); err!= nil{
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(showLoaded) 
}

func getProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var protocol models.Protocols
	var params = mux.Vars(r)
	db, _ := helper.Connect()
	

	var id string = params["id"]

	filter := bson.M{"id": id}
	err := db.Collection("protocol").FindOne(context.TODO(), filter).Decode(&protocol)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(protocol)
}

func updateProtocol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	var id string = params["id"]
	db, _ := helper.Connect()
	var protocol models.Protocols
	filter := bson.M{"id": id}
	_ = json.NewDecoder(r.Body).Decode(&protocol)

	update := bson.D{
		{"$set", bson.D{
			{"id", protocol.Id},
			{"type", protocol.Type},
			{"name", protocol.Name},
			{"items", protocol.Items},
		}},
	}

	err := db.Collection("protocol").FindOneAndUpdate(context.TODO(), filter, update).Decode(&protocol)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	protocol.Id = id

	json.NewEncoder(w).Encode(protocol)
}



func MainProtocols() {
	r := helper.Routes
	r.HandleFunc("/descriptor/protocols", getProtocols).Methods("GET")
	r.HandleFunc("/descriptor/protocol/{id}", getProtocol).Methods("GET")
	r.HandleFunc("/descriptor/protocol/{id}", updateProtocol).Methods("PUT")

}