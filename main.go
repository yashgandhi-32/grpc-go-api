package main

import (
	"fmt"
	"os"
	"os/signal"

	blogclient "github.com/yashgandhi-32/GRPC-API-CRUD/blog_client"
	blogserver "github.com/yashgandhi-32/GRPC-API-CRUD/blog_server"
)

func main() {
	go func() { blogserver.StartServer() }()
	go func() { blogclient.StartClient() }()
	// wait for control c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Wait until all signal received
	<-ch
	fmt.Print("Stopping the server")
	fmt.Print("Closing the listener")
}
