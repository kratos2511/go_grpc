// Package server implements server for BlogService
package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/kratos2511/go_grpc/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type server struct{}

// BlogItem struct holds the blog bson object
type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

func (*server) CreateBlog(ctx context.Context, blogReq *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	log.Println("CreateBlog Start")
	blog := blogReq.GetBlog()

	data := BlogItem{
		AuthorID: blog.GetAuthorID(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v", err),
		)
	}

	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(
			codes.Internal,
			fmt.Sprintf("Internal Error, Conversion failure:  %v", err),
		)
	}
	log.Println("Inserted into MongoDB, title: ", blog.GetTitle())
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			ID:       objID.Hex(),
			AuthorID: blog.GetAuthorID(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, readBlogRequest *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	log.Println("ReadBlog Start")
	objID, err := primitive.ObjectIDFromHex(readBlogRequest.GetBlogID())
	if err != nil {
		log.Println("Internal error, could not parse id. err: ", err)
		return nil, status.Error(codes.InvalidArgument, "Non parsable ID")
	}

	data := &BlogItem{}
	filter := bson.M{"_id": objID}
	res := collection.FindOne(ctx, filter)
	err = res.Decode(data)
	if err != nil {
		log.Println("Data not found, err: ", err)
		return nil, status.Errorf(codes.InvalidArgument, "Cannot find blog with specified ID. err: %v", err)
	}

	log.Println("Blog foung, title: ", data.Title)

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			ID:       data.ID.Hex(),
			AuthorID: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil
}

// Serve starts the server for BlogService
func Serve() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("BlogSever Init")

	//Setting up mongo colloction for Blog
	waitch := make(chan int)
	defer close(waitch)
	go getBlogCollection(waitch)

	lst, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalln("Failed to announce Listner")
	}

	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})
	reflection.Register(s)

	go func() {
		log.Println("Starting BlogServer")
		if err := s.Serve(lst); err != nil {
			log.Fatalln("Failed to serve. err: ", err)
		}
	}()

	//Signal based channel
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block on channel
	<-ch
	<-waitch
	log.Println("Stopping server")
	s.Stop()
	log.Println("Closing Listener")
	lst.Close()
}

func getBlogCollection(waitch chan<- int) {
	log.Println("Connecting to mongoDB")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalln("Failed connection to Mongo")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln("Client Connection failed with err: ", err)
	}
	collection = client.Database("blogDB").Collection("blog")

	defer func() {
		log.Println("Disconnecting client")
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalln("Client Disconnect failed with err: ", err)
		}
	}()
	waitch <- 1
}
