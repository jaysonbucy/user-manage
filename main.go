package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type User struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	dod_id string             `json:"dod_id,omitempty" bson:"dod_id,omitempty"`
	name   string             `json:"name,omitempty" bson:"name,omitempty"`
}

func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user User
	_ = json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("milhub_user").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}
func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("context-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex((params["id"]))
	var user User
	collection := client.Database("milhub_user").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `" }`))
	}
	json.NewEncoder(response).Encode(user)
}
func GetUsersEndpoint(response http.ResponseWriter, request *http.Request)

func main() {
	fmt.Println("Starting application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/user", CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/users", GetUsersEndpoint).Methods("GET")
	router.HandleFunc("/user{id}", GetUserEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}
