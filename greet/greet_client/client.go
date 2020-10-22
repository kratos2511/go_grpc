package main

import (
	"context"
	"log"

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
