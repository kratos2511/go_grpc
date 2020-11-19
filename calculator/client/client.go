package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/kratos2511/go_grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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

	//doStreamAvgRequest(c)

	//doSteamingMaximum(c)

	doSquareRoot(c)
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

func doSteamingMaximum(c calculatorpb.CalculateSumServiceClient) {
	log.Println("Init client")

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalln("Client encountered error while opening stream. err: ", err)
		return
	}

	waitc := make(chan struct{})
	rand.Seed(time.Now().UnixNano())

	go func() {
		for i := 0; i < 50; i++ {
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: rand.Int63(),
			})
			time.Sleep(100 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Println("Server send EOF")
				break
			} else if err != nil {
				log.Fatalln("Error reading stream. err: ", err)
				break
			}
			log.Println("Maximum: ", res.GetMaximum())
		}
		close(waitc)
	}()
	<-waitc
}

func doSquareRoot(c calculatorpb.CalculateSumServiceClient) {
	numbers := []int32{100, 4, 9, -23, 144}
	for _, i := range numbers {
		req := &calculatorpb.SquareRootRequest{
			Number: int32(i),
		}
		res, err := c.SquareRoot(context.Background(), req)
		if err != nil {
			resErr, ok := status.FromError(err)
			if ok {
				log.Println(resErr.Message(), resErr.Code())
			} else {
				log.Fatalf("Error %v", err)
			}
		} else {
			log.Printf("Response %v\n", res.GetNumberRoot())
		}
	}
}
