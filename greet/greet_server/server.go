package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/kratos2511/go_grpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	return &greetpb.GreetResponse{
		Result: fmt.Sprintf("Hello %v %v", firstName, lastName),
	}, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimeRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Println("GreetManyTimes invoked with", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	for i := 0; i < 50; i++ {
		result := &greetpb.GreetManyTimeResponse{
			Result: fmt.Sprintf("Hello %v %v. Packet: %v", firstName, lastName, i),
		}
		stream.Send(result)
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Println("LongGreet invoked with", stream)
	result := ""
	for {
		msg, err := stream.Recv()
		log.Println("LongGreet message: ", msg)
		if err == io.EOF {
			//Respond to client
			log.Println("LongGreet response to close stream")
			return stream.SendAndClose(&greetpb.LongGreetResponse{Result: result})
		}
		if err != nil {
			log.Fatalln("Server encountered error on recieve", err)
		}
		log.Println("Hello", msg.GetGreeting().GetFirstName(), msg.GetGreeting().GetLastName())
		result += "Hello " + msg.GetGreeting().GetFirstName() + " " + msg.GetGreeting().GetLastName() + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Println("GreetEveryone invoked with", stream)
	for {
		msg, err := stream.Recv()
		log.Println("GreetEveryone message: ", msg)
		if err == io.EOF {
			//Respond to client
			log.Println("GreetEveryone response to close stream")
			return nil
		} else if err != nil {
			log.Fatalln("GreetEveryone caught error, err: ", err)
			return err
		}
		result := "Hello " + msg.GetGreeting().GetFirstName() + " " + msg.GetGreeting().GetLastName()
		log.Println(result)
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalln("Server send error. err: ", err)
			return err
		}
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	log.Println("Init Greet with Deadline")

	time.Sleep(3 * time.Second)
	if ctx.Err() == context.Canceled {
		log.Println("Deadline exceeded, no need to furnish response.")
		return nil, status.Error(codes.Canceled, "Deadline has exceeded")
	}

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	return &greetpb.GreetWithDeadlineResponse{
		Result: fmt.Sprintf("Hello %v %v", firstName, lastName),
	}, nil
}

func main() {
	log.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
