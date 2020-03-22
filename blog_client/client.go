package blogclient

import (
	"context"
	"fmt"
	"log"

	"github.com/yashgandhi-32/GRPC-API-CRUD/blogproto"
	"google.golang.org/grpc"
)

func createBlog(c blogproto.BlogServieClient) {
	blog := &blogproto.CreateBlogRequest{
		Blog: &blogproto.Blog{
			AuthorId: "yash",
			Title:    "My first blog",
			Content:  "No content",
		},
	}
	resp, err := c.CreateBlog(context.Background(), blog)
	if err != nil {
		log.Fatalf("unexpected error %v", err)
	}
	fmt.Println("Blog has been created %v", resp)

	readblogresp, err := c.ReadBlog(context.Background(), &blogproto.ReadBlogRequest{
		BlogId: resp.GetBlog().GetId(),
	})
	if err != nil {
		fmt.Println("An error happened in fetching doc %v", err)
	}
	fmt.Println("Found blog:%v", readblogresp)
}

func StartClient() {
	fmt.Println("Client server started")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("An error occured in grpc client %v", err)
	}
	c := blogproto.NewBlogServieClient(cc)
	defer cc.Close()
	createBlog(c)
}
