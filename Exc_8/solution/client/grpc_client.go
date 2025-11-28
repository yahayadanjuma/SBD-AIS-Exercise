package client

import (
	"context" // for context.Background()
	"exc8/pb" // generated gRPC code
	"fmt"     // for fmt.Println / Printf

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcClient struct {
	client pb.OrderServiceClient
}

func NewGrpcClient() (*GrpcClient, error) {
	// grpc.NewClient does not exist; use grpc.Dial instead
	conn, err := grpc.Dial("localhost:4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewOrderServiceClient(conn)
	return &GrpcClient{client: client}, nil
}

func (c *GrpcClient) Run() error {
	fmt.Println("Requesting drinks 🍹🍺☕")                                          // print request for drinks
	drinksResp, err := c.client.GetDrinks(context.Background(), &emptypb.Empty{}) // call GetDrinks RPC
	if err != nil {                                                               // check for error from GetDrinks
		return err // return error if GetDrinks fails
	}
	fmt.Println("Available drinks:")      // print available drinks
	for _, d := range drinksResp.Drinks { // iterate over drinks
		fmt.Printf("\t> id:%d  name:\"%s\"  price:%d  description:\"%s\"\n", d.Id, d.Name, d.Price, d.Description) // print drink details
	}

	fmt.Println() // print a blank line before the next section

	fmt.Println("Ordering drinks 👨‍🍳⏱️🍻🍻") // print ordering drinks
	for _, d := range drinksResp.Drinks {  // iterate over drinks for first order
		req := &pb.OrderDrinkRequest{DrinkId: d.Id, Amount: 2}      // create order request for 2 drinks
		resp, err := c.client.OrderDrink(context.Background(), req) // call OrderDrink RPC
		if err != nil {                                             // check for error from OrderDrink
			return err // return error if OrderDrink fails
		}
		fmt.Printf("\t> Ordering: %d x %s\n", req.Amount, d.Name) // print order details
		fmt.Println("\t", resp.Message)                           // print order response message
	}

	fmt.Println() // print a blank line before the next section

	fmt.Println("Ordering another round of drinks 👨‍🍳⏱️🍻🍻") // print ordering another round
	for _, d := range drinksResp.Drinks {                   // iterate over drinks for second order
		req := &pb.OrderDrinkRequest{DrinkId: d.Id, Amount: 6}      // create order request for 6 drinks
		resp, err := c.client.OrderDrink(context.Background(), req) // call OrderDrink RPC
		if err != nil {                                             // check for error from OrderDrink
			return err // return error if OrderDrink fails
		}
		fmt.Printf("\t> Ordering: %d x %s\n", req.Amount, d.Name) // print order details
		fmt.Println("\t", resp.Message)                           // print order response message
	}

	fmt.Println() // print a blank line before the next section

	fmt.Println("Getting the bill 💹💹💹")                                           // print getting the bill
	ordersResp, err := c.client.GetOrders(context.Background(), &emptypb.Empty{}) // call GetOrders RPC
	if err != nil {                                                               // check for error from GetOrders
		return err // return error if GetOrders fails
	}
	totals := make(map[int32]int32)       // map to store totals
	for _, o := range ordersResp.Orders { // iterate over orders
		totals[o.DrinkId] += o.Amount // accumulate total amount per drink
	}
	for _, d := range drinksResp.Drinks { // iterate over drinks to print totals
		fmt.Printf("\t> Total: %d x %s\n", totals[d.Id], d.Name) // print total per drink
	}
	return nil // return nil on success
}
