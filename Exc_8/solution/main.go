package main

import (
	"exc8/client"
	"exc8/server"
	"time"
)

func main() {
	go func() {
		// Start the gRPC server // start the gRPC server in a goroutine
		if err := server.StartGrpcServer(); err != nil { // run the server and check for errors
			panic(err) // panic if server fails
		}
	}()
	time.Sleep(1 * time.Second) // wait for server to start
	// Start the gRPC client // create and run the gRPC client
	c, err := client.NewGrpcClient() // initialize the client
	if err != nil {                  // check for client creation error
		panic(err) // panic if client fails
	}
	if err := c.Run(); err != nil { // run client logic and check for errors
		panic(err) // panic if client run fails
	}
	println("Orders complete!") // print completion message
}
