//Package client implements client for BlogService
package client

import (
	"context"
	"fmt"
	"log"

	"github.com/kratos2511/go_grpc/blog/blogpb"
	"google.golang.org/grpc"
)

// Request sets up the client and make requests to the
// BlogService Server
func Request() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("BlogService Client Init")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Failed to Dial to server. err: ", err)
		return
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)
	log.Println("Successfully initialised client: ", c)
	requestServer(c)
}

func requestServer(c blogpb.BlogServiceClient) {
	for {
		showMenu()
		switch getOption() {
		case 1:
			createBlog(c)
		case 2:
			readBlog(c)
		default:
			fmt.Println("Thanks for using BlogService Client. Bye..")
			return
		}
	}

}

func showMenu() {
	fmt.Println("Select option by entering number from 0 to 1")
	fmt.Println("1. Create Blog")
	fmt.Println("2. Read Blog")
	fmt.Println("0. Quit")
}

func getOption() int {
	input := 0
	fmt.Scanf("%d", &input)
	return input
}

func createBlog(c blogpb.BlogServiceClient) {
	log.Println("CreateBlog RPC start")
	blogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorID: "Rahul Sachan",
			Title:    "New Blog",
			Content:  "Content for new blog",
		},
	})
	if err != nil {
		log.Println("CreateBlog Failed, err: ", err)
	} else {
		log.Println("CreateBlog Successful. id: ", blogRes.GetBlog().GetID())
	}
}

func readBlog(c blogpb.BlogServiceClient) {
	log.Println("ReadBlog RPC start")
	blogRes, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogID: "5fb8221561f6c2bb0ef890d3",
	})
	if err != nil {
		log.Println("ReadBlog Failed, err: ", err)
	} else {
		log.Println("ReadBlog Successful. id: ", blogRes.GetBlog().GetTitle())
	}
}
