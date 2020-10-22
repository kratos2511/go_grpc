package main

import (
	"context"
	"io"
	"log"

	"github.com/kratos2511/go_grpc/primef/primefpb"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Init Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Issue connecting to server", err)
	}
	c := primefpb.NewPrimeFactorServiceClient(cc)

	doFactorStreaming(c)
}

func doFactorStreaming(c primefpb.PrimeFactorServiceClient) {
	req := &primefpb.GetPrimeFactorsRequest{
		Number: 120,
	}
	if resStream, err := c.GetPrimeFactors(context.Background(), req); err != nil {
		log.Fatalln("Request Failed", err)
	} else {
		for {
			msg, err := resStream.Recv()
			if err == io.EOF {
				return
			}

			log.Println("Response: ", msg)
		}
	}
}
