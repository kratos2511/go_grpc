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

	//makeRequest(c)

	doStreamAvgRequest(c)
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

func doStreamAvgRequest(c calculatorpb.CalculateSumServiceClient) {
	numbers := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stream, err := c.StreamAvg(context.Background())
	if err != nil {
		log.Fatalln("There was an error connecting to the server. error:", err)
	}
	for _, number := range numbers {
		log.Println("Client sending number:", number)
		if err := stream.Send(&calculatorpb.StreamAvgRequest{Number: number}); err != nil {
			log.Fatalln("There was error sending message. error:", err)
		}
	}
	if res, err := stream.CloseAndRecv(); err != nil {
		log.Fatalln("There was an issue closing the stream. error:", err)
	} else {
		log.Println("Client recieved average:", res.GetAverage())
	}
}
