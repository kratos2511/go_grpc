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
