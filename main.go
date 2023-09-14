package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

var user1 = User{primitive.NewObjectID(), "ali", "ali.ali@yahoo.com", "abcdef"}

var client *mongo.Client

// Helper function to get the user collection from MongoDB
func getUserCollection() *mongo.Collection {
	return client.Database("myproject").Collection("user")
}

// Create a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse the request body into a User struct
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the user into the MongoDB collection
	result, err := getUserCollection().InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created user as JSON with the ID
	user.ID = result.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(user)
}

func main() {

	// Initialize the MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// create router with NewRouter method and assign it to instance r
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.go run main.go

	//after declaring a new router instance, you can use the HandleFunc method of your router
	//instance to assign routes to handler functions along with the request type that the
	//handler function handles.

	r.HandleFunc("/create", createUser).Methods("POST")
	//r.HandleFunc("/getuser/{name}", getUser).Methods("GET")
	//r.HandleFunc("/getusers", getUsers).Methods("GET")
	//r.HandleFunc("/update/{name}", updateUser).Methods("PUT")
	//r.HandleFunc("/patch/{name}", putsomedata).Methods("PATCH")

	//http.ListenAndServe() function to start the server and tell it to listen for
	//new HTTP requests and
	//then serve them using the handler functions you set up
	//You can set up a server using the ListenAndServe method of the http package.
	//The ListenAndServe method takes as arguments the port you want the server to run on
	//and a router instance
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":27017", r))
}
