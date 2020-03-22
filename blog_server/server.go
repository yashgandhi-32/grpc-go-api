package blogserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/yashgandhi-32/GRPC-API-CRUD/blogproto"
	"github.com/yashgandhi-32/GRPC-API-CRUD/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
}

var mongoconn *mongodb.ConnectionManager

func (server) CreateBlog(ctx context.Context, req *blogproto.CreateBlogRequest) (*blogproto.CreateBlogReseponse, error) {
	blog := req.GetBlog()
	data := mongodb.BlogItem{
		AutorID: blog.GetAuthorId(),
		Content: blog.GetContent(),
		Title:   blog.GetTitle(),
		ID:      blog.GetId(),
	}
	res, err := mongoconn.Db.Collection("xyz").InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error %v", err),
		)
	}
	oid := res.InsertedID.(primitive.ObjectID)
	return &blogproto.CreateBlogReseponse{
		Blog: &blogproto.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Content:  blog.GetContent(),
			Title:    blog.GetTitle(),
		},
	}, nil
}

func initMongo() {
	err, conn := mongodb.ConnectDB()
	if err != nil {
		fmt.Print(err.Msg)
	}
	mongoconn = conn
}

func StartServer() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Print("Blog Service started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("An error occured in starting server %v", err)
	}

	//Start Mongodb
	initMongo()

	s := grpc.NewServer()
	blogproto.RegisterBlogServieServer(s, server{})
	go func() {
		fmt.Print("Server Starting")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve %v", err)
		}
	}()
	// wait for control c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Wait until all signal received
	<-ch
	fmt.Print("Stopping the server")
	s.Stop()
	fmt.Print("Closing the listener")
	lis.Close()
}
