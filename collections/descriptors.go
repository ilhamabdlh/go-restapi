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

)


func getDescriptors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// db, _ := helper.Connect()
	// var descriptors []models.Response
	

	// cur, err := db.Collection("descriptors").Find(context.TODO(), bson.M{})
	// conf, _ := db.Collection("configs").Find(context.TODO(), bson.M{})

	// if err != nil {
	// 	helper.GetError(err, w)
	// 	return
	// }
	// defer cur.Close(context.TODO())

	// var response models.Response
	// for cur.Next(context.TODO()) {

	// 	var descriptor models.Descriptor
	// 	err := cur.Decode(&descriptor) 
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	var status int
	// 	if err != nil {
	// 		status = 400
	// 	} else {
	// 		status = 200
	// 	}

		
	// 	response.Data = descriptor
	// 	response.Status = strconv.Itoa(status)
	// 	response.Success = true
	// 	response.Msg = http.StatusText(status)
		

		
	// }

	// var config models.Config
	// erre := conf.Decode(&config) 
	// if erre != nil {
	// 	log.Fatal(err)
	// }
	// if response.Data.Id == config.Id{
	// 	response.Data.Configs[0] = config
	// }

	// descriptors = append(descriptors, response)

	// if err := cur.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	// json.NewEncoder(w).Encode(descriptors)
	


	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// lookupStageTwo := bson.D{{"$lookup", bson.D{{"from", "configs"}, {"localField", "id"}, {"foreignField", "id"}, {"as", "configs"}}}}
	// lookupStageThree := bson.D{{"$lookup", bson.D{{"from", "protocol"}, {"localField", "id"}, {"foreignField", "configs.id"}, {"as", "configs.protocol"}}}}
	// unwindStage := bson.D{{"$unwind", bson.D{{"path", "$configs"}, {"preserveNullAndEmptyArrays", false}}}}

	// showLoadedCursor, err := db.Collection("descriptors").Aggregate(ctx, mongo.Pipeline{lookupStageTwo, unwindStage, lookupStageThree})
	// if err !=nil{
	// 	log.Fatal(err)
	// }

	// var showLoaded []bson.M
	// if err = showLoadedCursor.All(ctx, &showLoaded); err!= nil{
	// 	log.Fatal(err)
	// }

	

	// descriptors = append(descriptors, descriptor)




	var descriptors []models.Response
	db, _ := helper.Connect()

	cur, err := db.Collection("descriptors").Find(context.TODO(), bson.M{})

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


	json.NewEncoder(w).Encode(descriptors)

	

}

func getDescriptor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

// var collectionConfigDes = helper.ConnectConfigsDB()
// var collectionStatusDes = helper.ConnectStatusesDB()
// var collectionProtocolDes = helper.ConnectProtocolsDB()
// var collectionItemDes = helper.ConnectItemsDB()
func createDescriptor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var descriptor models.Descriptor
	var config models.Config
	var status models.Statuses
	var protocolOne models.Protocols
	var protocolTwo models.Protocols
	var itemOne models.Itemes
	var itemTwo models.Itemes
	db, _ := helper.Connect()

	_ = json.NewDecoder(r.Body).Decode(&descriptor)
	result, _ := db.Collection("descriptors").InsertOne(context.TODO(), descriptor)

	itemOne = descriptor.Configs[0].Protocol[0].Items[0]
	itemTwo = descriptor.Status[0].Protocol[0].Items[0]
	_ = json.NewDecoder(r.Body).Decode(&itemOne)
	_ = json.NewDecoder(r.Body).Decode(&itemTwo)
	

	

	descriptor.Configs[0].Protocol[0].Items = make([]models.Itemes, 0)
	descriptor.Status[0].Protocol[0].Items = make([]models.Itemes, 0)
	
	protocolOne = descriptor.Configs[0].Protocol[0]
	protocolTwo = descriptor.Status[0].Protocol[0]



	// bar := models.Protocols{
	// 	Id: descriptor.Configs[0].Protocol[0].Id,
	// 	Type: descriptor.Configs[0].Protocol[0].Type,
	// 	Name: descriptor.Configs[0].Protocol[0].Name,
	// 	Items: make([]models.Itemes, 0),
	// }

	// bor := models.Statuses{
	// 	Id: descriptor.Status[0].Id,
	// 	Type: descriptor.Status[0].Type,
	// 	Name: descriptor.Status[0].Name,
	// 	Protocol: make([]models.Protocols, 0),
	// }

	// protocoler, errorr := json.Marshal(struct{
	// 	*models.Statuses
	// 	Protocol string `json:"protocol"`
	// }{
	// 	Statuses: &bor,
	// 	Protocol: "",
	// })

	// if errorr != nil {
	// 	panic(errorr)
	// }
	// gerd  := string(protocoler)

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

	// S := models.Statuses{}
	// C := models.Config{}
	// desc := models.Descriptor{
	// 	Id: descriptor.Id,
	// 	Type: descriptor.Type,
	// 	Name: descriptor.Name,
	// 	Version: descriptor.Version,
	// 	Modules: descriptor.Modules,
	// 	Configs: C,
	// 	Status: S,
	// } 

	// val, _ := json.Marshal(desc)
	// value := string(val)

	// descriptor.Status = make([]models.Statuses, 0)
	// descriptor.Configs = make([]models.Config, 0)

	// _ = json.NewDecoder(r.Body).Decode(&descriptor)
	// result, _ := db.Collection("descriptors").InsertOne(context.TODO(), descriptor)

	
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