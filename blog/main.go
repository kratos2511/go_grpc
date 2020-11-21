package main

import (
	"log"
	"os"

	"github.com/kratos2511/go_grpc/blog/client"
	"github.com/kratos2511/go_grpc/blog/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	switch getCommand() {
	case "serve":
		log.Println("Requested server")
		server.Serve()
	case "client":
		log.Println("Requested client")
		client.Request()
	}
}

func getCommand() string {
	if len(os.Args) == 1 {
		return "serve"
	}
	return os.Args[1]
}
