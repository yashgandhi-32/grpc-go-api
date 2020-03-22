package blogserver

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/yashgandhi-32/GRPC-API-CRUD/blogproto"
	"github.com/yashgandhi-32/GRPC-API-CRUD/mongodb"
	"go.mongodb.org/mongo-driver/bson"
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
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
		ID:       blog.GetId(),
	}
	res, err := mongoconn.Db.Collection("posts").InsertOne(context.Background(), data)
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

func (server) ReadBlog(ctx context.Context, req *blogproto.ReadBlogRequest) (*blogproto.ReadBlogResponse, error) {
	blogId := req.GetBlogId()
	data := &mongodb.BlogItem{}
	id, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse Id"),
		)
	}
	filter := bson.M{"_id": id}
	err = mongoconn.Db.Collection("posts").FindOne(context.Background(), filter).Decode(data)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("No document found"),
		)
	}
	return &blogproto.ReadBlogResponse{
		Blog: &blogproto.Blog{
			Id:       data.ID,
			Title:    data.Title,
			Content:  data.Content,
			AuthorId: data.AuthorID,
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
	log.Println("Blog Service started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("An error occured in starting server %v", err)
	}

	//Start Mongodb
	initMongo()

	s := grpc.NewServer()
	blogproto.RegisterBlogServieServer(s, server{})
	go func() {
		fmt.Println("Server Starting")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve %v", err)
		}
	}()
}
