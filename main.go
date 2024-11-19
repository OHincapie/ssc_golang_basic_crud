package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://admin:abcd1234@cluster0.x02gr.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Database("scc").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	collection := client.Database("scc").Collection("user")
	initialUser := User{
		Name: "test",
	}
	result, err := collection.InsertOne(context.Background(), initialUser)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted a single document: ", result.InsertedID)
}

type User struct {
	ID          uint64 `json:"identification"`
	Name        string `json:"name"`
	LastName    string `json:"lastName"`
	PublicForce string `json:"publicForce"`
	Range       string `json:"range"`
	ForceID     int    `json:"forceId"`
	Email       string `json:"email"`
}
