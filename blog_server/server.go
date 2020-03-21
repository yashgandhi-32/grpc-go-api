package blogserver

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/yashgandhi-32/GRPC-API-CRUD/blogproto"
	"google.golang.org/grpc"
)

type server struct {
}

func StartServer() {
	log.Print("Blog Service started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("An error occured in starting server %v", err)
	}

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
