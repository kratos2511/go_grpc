package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/kratos2511/go_grpc/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Hello I'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)

	log.Printf("created client: %f", c)

	//doUnaryRequest(c)
	//doServerStreamingRequest(c)
	doClientStream(c)
}

func doUnaryRequest(c greetpb.GreetServiceClient) {
	for i := 0; i < 50; i++ {
		req := &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rahul",
				LastName:  "Sachan",
			},
		}
		if res, err := c.Greet(context.Background(), req); err != nil {
			log.Fatalf("error while calling Greet RPC %v", err)
		} else {
			log.Printf("Respose %v", res)
		}
	}
}

func doServerStreamingRequest(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimeRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rahul",
			LastName:  "Sachan",
		},
	}
	if resStream, err := c.GreetManyTimes(context.Background(), req); err != nil {
		log.Fatalf("Error calling GreetManyTimes RPC %v", err)
	} else {
		for {
			msg, err := resStream.Recv()
			if err == io.EOF {
				return
			}
			log.Printf("Respose: %v", msg)
		}
	}
}

func doClientStream(c greetpb.GreetServiceClient) {
	log.Println("Client Streaming Init")
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalln("Error calling server")
	}
	data := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rahul",
				LastName:  "Sachan",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Abhay",
				LastName:  "Sachan",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Isha",
				LastName:  "Singh",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Pratesh",
				LastName:  "Jhari",
			},
		},
	}
	for _, val := range data {
		log.Println("Client sending data")
		stream.Send(val)
		time.Sleep(100 * time.Millisecond)
	}

	if res, err := stream.CloseAndRecv(); err != nil {
		log.Fatalln("Client encountered error", err)
	} else {
		log.Println("Response:", res)
	}

}
