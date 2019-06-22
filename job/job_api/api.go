package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Company string             `json:"company" bson:"company"`
	Title   string             `json:"title" bson:"title"`
	Content string             `json:"content" bson:"content"`
	Region  string             `json:"region" bson:"region"`
	Country string             `json:"country" bson:"country"`
}

var client *mongo.Client

func CreateJobEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var job Job
	_ = json.NewDecoder(request.Body).Decode(&job)
	fmt.Print(job.Title)
	collection := client.Database("visasponsorshipjobs").Collection("jobs")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	result, _ := collection.InsertOne(ctx, job)
	json.NewEncoder(response).Encode(result)
}

func GetJobEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	vars := mux.Vars(request)
	fmt.Println(vars["id"])
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	fmt.Println(id)
	collection := client.Database("visasponsorshipjobs").Collection("jobs")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	var job Job
	filter := bson.M{"_id": id}
	err := collection.FindOne(ctx, filter).Decode(&job)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(job)

}

func GetJobsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var jobs []Job
	collection := client.Database("visasponsorshipjobs").Collection("jobs")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var job Job
		cursor.Decode(&job)
		jobs = append(jobs, job)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	json.NewEncoder(response).Encode(jobs)
}

func main() {
	fmt.Println("Starting the application ..")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()

	router.HandleFunc("/jobs/new", CreateJobEndpoint).Methods("POST")
	router.HandleFunc("/jobs", GetJobsEndpoint).Methods("GET")
	router.HandleFunc("/jobs/{id}", GetJobEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}
