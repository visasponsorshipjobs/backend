package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"visasponsorshipjobs/backend/job/jobpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct {
}

type jobItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	RecruiterID string             `bson:"recruiter_id"`
	Title       string             `bson:"title"`
	Content     string             `bson:"content"`
	Region      string             `bson:"region"`
	Country     string             `bson:"country"`
}

func (*server) ReadJob(ctx context.Context, req *jobpb.ReadJobRequest) (*jobpb.ReadJobResponse, error) {
	jobId := req.GetJobId()
	oid, err := primitive.ObjectIDFromHex(jobId)

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	//create an empty struct
	data := &jobItem{}

	filter := bson. .NewDocument(bson.EC.ObjectID("_id", oid))
	result := collection.FindOne(context.Background(), filter)

}

func (*server) CreateJob(ctx context.Context, req *jobpb.CreateJobRequest) (*jobpb.CreateJobResponse, error) {
	job := req.GetJob()

	data := jobItem{
		RecruiterID: job.GetRecruiterId(),
		Title:       job.GetContent(),
		Content:     job.GetContent(),
		Region:      job.GetRegion(),
		Country:     job.GetCountry(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cant convert to OID: %v", err),
		)
	}

	return &jobpb.CreateJobResponse{
		Job: &jobpb.Job{
			Id:          oid.Hex(),
			RecruiterId: job.GetRecruiterId(),
			Title:       job.GetTitle(),
			Content:     job.GetContent(),
			Region:      job.GetRegion(),
			Country:     job.GetCountry(),
		},
	}, nil

}

func main() {

	//if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Job service started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Faield to listen: %v", err)
	}

	s := grpc.NewServer()
	jobpb.RegisterJobServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting server ...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//BLock until a signal is received
	<-ch

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("closing the listener")
	lis.Close()
	fmt.Println("End of Program")
}
