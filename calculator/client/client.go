package main

import (
	"context"
	"log"

	"github.com/kratos2511/go_grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Calculator client init")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error initializing client %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculateSumServiceClient(cc)
	log.Println("Client initialized")

	makeRequest(c)
}

func makeRequest(c calculatorpb.CalculateSumServiceClient) {
	numbers := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := range numbers {
		req := &calculatorpb.CalculateSumRequest{
			Numbers: &calculatorpb.Numbers{
				Numbers: numbers[0 : i+1],
			},
		}
		if res, err := c.CalculateSum(context.Background(), req); err != nil {
			log.Fatalf("Error %v", err)
		} else {
			log.Printf("Response %v\n", res)
		}
	}
}
