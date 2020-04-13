package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Dod_id    string             `json:"dod_id,omitempty" bson:"dod_id,omitempty"`
	User_name string             `json:"user_name,omitempty" bson:"user_name,omitempty"`
}

func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user User
	_ = json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("milhub_user").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}
func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("context-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	collection := client.Database("milhub_user").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `" }`))
	}
	json.NewEncoder(response).Encode(user)
}
func GetUsersEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var users []User
	collection := client.Database("milhub_user").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(users)
}

func main() {
	fmt.Println("Starting application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	fmt.Println(client)
	router := mux.NewRouter()
	router.HandleFunc("/user", CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/users", GetUsersEndpoint).Methods("GET")
	router.HandleFunc("/user/{id}", GetUserEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}
