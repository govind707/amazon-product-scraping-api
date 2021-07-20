package main

import (
	"time"
	"context"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type dataExtractedMap struct {
	Name			string	`json:"name,omitempty" bson:"name,omitempty"`
	ImageURL		string	`json:"imageURL,omitempty" bson:"imageURL,omitempty"`
	Desc			string	`json:"description,omitempty" bson:"description,omitempty"`
	Price			string	`json:"price,omitempty" bson:"price,omitempty"`
	TotalReviews	string		`json:"totalReviews,omitempty" bson:"totalReviews,omitempty"`
}

type respMapObj struct {
	ID		primitive.ObjectID	`json:"_id,omitempty" bson:"_id,omitempty"`
	URL		string				`json:"url,omitempty" bson:"url,omitempty"`
	Product	dataExtractedMap				`json:"product,omitempty" bson:"product,omitempty"`
	LastUpdate	time.Time		`json:"last_update, omitempty" bson:"last_update, omitempty"`
}

var client *mongo.Client

func CheckDocument(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var new_doc, existing_doc respMapObj
	_ = json.NewDecoder(request.Body).Decode(&new_doc)
	collection := client.Database("amazondb").Collection("amazoncollection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	
	collection.FindOne(ctx, bson.M{"url":new_doc.URL}).Decode(&existing_doc)
	new_doc.LastUpdate = time.Now()
	if existing_doc.URL == "" {
		result, _ := collection.InsertOne(ctx, new_doc)
		json.NewEncoder(response).Encode(result)
	} else {
		result, _ := collection.UpdateOne(	ctx,
								bson.M{"url": new_doc.URL},
								bson.D{
									primitive.E{
										Key: "$set",
										Value: bson.D{
												primitive.E{
													Key: "product",
													Value: new_doc.Product,
												},
										},
									},
								},
							)
		json.NewEncoder(response).Encode(result)
	}
	
}

func SelectAllDocuments(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var docs []respMapObj
	collection := client.Database("amazondb").Collection("amazoncollection")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var doc respMapObj
		cursor.Decode(&doc)
		docs = append(docs, doc)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	
	json.NewEncoder(response).Encode(docs)
}

func main()  {
	localhost := "mongodb://database:27017"
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(localhost)
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/collector", CheckDocument).Methods("POST")
	router.HandleFunc("/collector", SelectAllDocuments).Methods("GET")
	http.ListenAndServe(":3031", router)
}