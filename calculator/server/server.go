package main

import (
	"context"
	"log"
	"net"

	"github.com/kratos2511/go_grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) CalculateSum(ctx context.Context, req *calculatorpb.CalculateSumRequest) (*calculatorpb.CalculateSumResponse, error) {
	log.Println("Serving request:", req)
	numbers := req.GetNumbers().GetNumbers()
	var sum int32 = 0
	for _, number := range numbers {
		sum += number
	}

	return &calculatorpb.CalculateSumResponse{
		Sum: sum,
	}, nil
}

func main() {
	log.Println("Calculator Server Init")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to open a listner %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculateSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}

}
