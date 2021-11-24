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
func getStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "protocols"}, {"localField", "id"}, {"foreignField", "id"}, {"as", "protocol"}}}}
	lookupStageTwo := bson.D{{"$lookup", bson.D{{"from", "items"}, {"localField", "id"}, {"foreignField", "id"}, {"as", "protocol.items"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$protocol"}, {"preserveNullAndEmptyArrays", false}}}}

	db, _ := helper.Connect()
	showLoadedCursor, err := db.Collection("statuses").Aggregate(ctx, mongo.Pipeline{lookupStageTwo, unwindStage, lookupStage})
	if err !=nil{
		log.Fatal(err)
	}

	var ShowLoaded []bson.M
	if err = showLoadedCursor.All(ctx, &ShowLoaded); err!= nil{
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(ShowLoaded) 
}

func getConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	db, _ := helper.Connect()

	lookupStage := bson.D{{"$lookup", bson.D{{"from", "protocols"}, {"localField", "id"}, {"foreignField", "id"}, {"as", "protocol"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$id"}, {"preserveNullAndEmptyArrays", false}}}}

	ShowLoadedCursor, err := db.Collection("configs").Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
	if err !=nil{
		log.Fatal(err)
	}

	var ShowLoaded []bson.M
	if err = ShowLoadedCursor.All(ctx, &ShowLoaded); err!= nil{
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