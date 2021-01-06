package chat

import (
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

var latestID interface{}
var collection *mongo.Collection

func Connect() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("testgo1").Collection("names")
}

type Server struct {
	UnimplementedChatServiceServer
}

func (s *Server) SetHello(ctx context.Context, message *HelloRequest) (*HelloReply,
	error) {
	res, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		fmt.Println(err)
	}
	latestID = res.InsertedID.(primitive.ObjectID)
	return &HelloReply{Message: "Hello" + " " + message.Name}, nil
}

func (s *Server) GetHello(ctx context.Context, message *HelloRequest) (*HelloReply,
	error) {
	name := collection.FindOne(context.TODO(), bson.M{"_id": latestID})
	return &HelloReply{Message: "hello" + " " + fmt.Sprint(name.Decode(&bson.M{}))}, nil
}
