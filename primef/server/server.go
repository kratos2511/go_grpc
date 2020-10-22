package main

import (
	"log"
	"net"

	"github.com/kratos2511/go_grpc/primef/primefpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) GetPrimeFactors(req *primefpb.GetPrimeFactorsRequest, stream primefpb.PrimeFactorService_GetPrimeFactorsServer) error {
	log.Println("Request Recieved", req)
	number := req.GetNumber()
	var k int32 = 2
	for number > 1 {
		if number%k == 0 {
			stream.Send(&primefpb.GetPrimeFactorsResponse{Factor: k})
			number = number / k
		} else {
			k = k + 1
		}
	}
	return nil
}

func main() {
	log.Println("Init Server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalln("Server binding failed", err)
	}
	s := grpc.NewServer()

	primefpb.RegisterPrimeFactorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed launching server", err)
	}

}
