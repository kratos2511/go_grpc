package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/kratos2511/go_grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) StreamAvg(stream calculatorpb.CalculateSumService_StreamAvgServer) error {
	log.Println("Init StreamAvg request")
	var sum float32
	var count float32

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("Responding to the request")
			return stream.SendAndClose(&calculatorpb.StreamAvgResponse{Average: sum / count})
		} else if err != nil {
			log.Fatalln("Error recieveing params, error: ", err)
		}
		count++
		sum += float32(msg.GetNumber())
		log.Println("total:", sum, "\tcount:", count)
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculateSumService_FindMaximumServer) error {
	log.Println("Init FindMaximum request")
	max := int64(0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("Finished request stream")
			return err
		} else if err != nil {
			log.Fatalln("Request recieve error. err: ", err)
			return err
		}

		log.Println("Number recieved: ", msg.GetNumber())
		if msg.GetNumber() > max {
			max = msg.GetNumber()
		}
		sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
			Maximum: max,
		})
		if sendErr != nil {
			log.Fatalln("Server encountered error while responding. err: ", sendErr)
			return sendErr
		}
	}
}

func (*server) SquareRoot(ctx context.Context, request *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	log.Println("Init SquareRoot")
	number := request.GetNumber()
	log.Println("Requested Number: ", number)
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Recieved a negative number %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
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
