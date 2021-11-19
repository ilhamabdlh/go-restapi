package collections

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	// "time"
	"strconv"
	
	"github.com/ilhamabdlh/go-restapi/helper"
	"github.com/ilhamabdlh/go-restapi/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)
var collectionDescriptor = helper.ConnectDescriptorsDB()


func getDescriptors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var descriptors []models.Response

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

	var status int
	if err != nil {
		status = 400
	} else {
		status = 200
	}

	var response models.Response
	response.Data = descriptor
	response.Status = strconv.Itoa(status)
	response.Success = true
	response.Msg = http.StatusText(status)

	descriptors = append(descriptors, response)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	











	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// lookupStage := bson.D{{"$lookup", bson.D{{"from", "configs"}, {"localField", "id"}, {"foreignField", "id"}, {"as", "configs"}}}}
	// lookupStageTwo := bson.D{{"$lookup", bson.D{{"from", "protocols"}, {"localField", "configs.id"}, {"foreignField", "id"}, {"as", "configs.protocol"}}}}
	// lookupStageThree := bson.D{{"$lookup", bson.D{{"from", "items"}, {"localField", "configs.protocol.id"}, {"foreignField", "id"}, {"as", "configs.protocol.items"}}}}
	
	// lookupStageFour := bson.D{{"$lookup", bson.D{{"from", "statuses"}, {"localField", "id"}, {"foreignField", "id"}, {"as", "status"}}}}
	// lookupStageFive := bson.D{{"$lookup", bson.D{{"from", "protocols"}, {"localField", "status.id"}, {"foreignField", "id"}, {"as", "status.protocol"}}}}
	// lookupStageSix := bson.D{{"$lookup", bson.D{{"from", "items"}, {"localField", "status.protocol.id"}, {"foreignField", "id"}, {"as", "status.protocol.items"}}}}
	
	// unwindStage := bson.D{{"$unwind", bson.D{{"path", "$configs"}, {"preserveNullAndEmptyArrays", false}}}}
	// unwindStageTwo := bson.D{{"$unwind", bson.D{{"path", "$configs.protocol"}, {"preserveNullAndEmptyArrays", false}}}}
	// unwindStageThree := bson.D{{"$unwind", bson.D{{"path", "$configs.protocol.items"}, {"preserveNullAndEmptyArrays", false}}}}
	
	// unwindStageFour := bson.D{{"$unwind", bson.D{{"path", "$status"}, {"preserveNullAndEmptyArrays", false}}}}
	// unwindStageFive := bson.D{{"$unwind", bson.D{{"path", "$status.protocol"}, {"preserveNullAndEmptyArrays", false}}}}
	// unwindStageSix := bson.D{{"$unwind", bson.D{{"path", "$status.protocol.items"}, {"preserveNullAndEmptyArrays", false}}}}
	

	// showLoadedCursor, err := collectionDescriptor.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage, lookupStageTwo, unwindStageTwo, lookupStageThree, unwindStageThree, lookupStageFour, unwindStageFour, lookupStageFive, unwindStageFive, lookupStageSix, unwindStageSix})
	// if err !=nil{
	// 	log.Fatal(err)
	// }

	// var showLoaded []bson.M
	// if err = showLoadedCursor.All(ctx, &showLoaded); err!= nil{
	// 	log.Fatal(err)
	// }
	

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

	var status int
	if err != nil {
		status = 400
	} else {
		status = 200
	}

	var response models.Response
	response.Data = descriptor
	response.Status = strconv.Itoa(status)
	response.Success = true
	response.Msg = http.StatusText(status)

	

	json.NewEncoder(w).Encode(response)
}

var collectionConfigDes = helper.ConnectConfigsDB()
var collectionStatusDes = helper.ConnectStatusesDB()
var collectionProtocolDes = helper.ConnectProtocolsDB()
var collectionItemDes = helper.ConnectItemsDB()
func createDescriptor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var descriptor models.Descriptor
	var config models.Config
	var status models.Statuses
	var protocolOne models.Protocols
	var protocolTwo models.Protocols
	var itemOne models.Itemes
	var itemTwo models.Itemes

	_ = json.NewDecoder(r.Body).Decode(&descriptor)
	result, _ := collectionDescriptor.InsertOne(context.TODO(), descriptor)

	
	status = descriptor.Status
	config = descriptor.Configs
	protocolOne = descriptor.Configs.Protocol
	protocolTwo = descriptor.Status.Protocol
	itemOne = descriptor.Configs.Protocol.Items
	itemTwo = descriptor.Status.Protocol.Items

	_ = json.NewDecoder(r.Body).Decode(&status)
	_ = json.NewDecoder(r.Body).Decode(&config)
	_ = json.NewDecoder(r.Body).Decode(&protocolOne)
	_ = json.NewDecoder(r.Body).Decode(&protocolTwo)
	_ = json.NewDecoder(r.Body).Decode(&itemOne)
	_ = json.NewDecoder(r.Body).Decode(&itemTwo)

	conf, errr := collectionConfigDes.InsertOne(context.TODO(), config)
	stat, _ := collectionStatusDes.InsertOne(context.TODO(), status)
	ptOne, _ := collectionProtocolDes.InsertOne(context.TODO(), protocolOne)
	ptTwo, _ := collectionProtocolDes.InsertOne(context.TODO(), protocolTwo)
	itOne, _ := collectionItemDes.InsertOne(context.TODO(), itemOne)
	itTwo, _ := collectionItemDes.InsertOne(context.TODO(), itemTwo)

	
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