package server

import (
	"context"
	"exc8/pb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCService struct {
	pb.UnimplementedOrderServiceServer
}

// Predefined drinks
var drinks = []*pb.Drink{
	{Id: 1, Name: "Spritzer", Price: 2, Description: "Wine with soda"},
	{Id: 2, Name: "Beer", Price: 3, Description: "Hagenberger Gold"},
	{Id: 3, Name: "Coffee", Price: 0, Description: "Mifare isn't that secure"},
}

// In-memory totals per drink
var orderTotals = map[int32]int32{}

func StartGrpcServer() error {
	srv := grpc.NewServer()
	grpcService := &GRPCService{}
	pb.RegisterOrderServiceServer(srv, grpcService)

	lis, err := net.Listen("tcp", ":4000")
	if err != nil {
		return err
	}

	log.Println("gRPC server listening on :4000")
	if err := srv.Serve(lis); err != nil {
		return err
	}
	return nil
}

// GetDrinks implements OrderService.GetDrinks
func (s *GRPCService) GetDrinks(ctx context.Context, _ *emptypb.Empty) (*pb.ListDrinksResponse, error) {
	return &pb.ListDrinksResponse{Drinks: drinks}, nil // return the list of predefined drinks
}

// OrderDrink implements OrderService.OrderDrink
func (s *GRPCService) OrderDrink(ctx context.Context, req *pb.OrderDrinkRequest) (*pb.OrderDrinkResponse, error) {
	if req.Amount <= 0 { // validate that amount is positive
		return &pb.OrderDrinkResponse{
			Success: false,
			Message: "amount must be > 0",
		}, fmt.Errorf("invalid amount: %d", req.Amount) // return error for invalid amount
	}

	var exists bool            // check that drink exists
	for _, d := range drinks { // iterate over drinks
		if d.Id == req.DrinkId { // match drink id
			exists = true // mark as found
			break         // exit loop
		}
	}
	if !exists { // if drink not found
		return &pb.OrderDrinkResponse{
			Success: false,
			Message: fmt.Sprintf("unknown drink id: %d", req.DrinkId),
		}, fmt.Errorf("unknown drink id: %d", req.DrinkId) // return error for unknown drink
	}

	orderTotals[req.DrinkId] += req.Amount // update in-memory order total

	return &pb.OrderDrinkResponse{
		Success: true,
		Message: fmt.Sprintf("ordered %d of drink %d", req.Amount, req.DrinkId),
	}, nil // return success response
}

// GetOrders implements OrderService.GetOrders
func (s *GRPCService) GetOrders(ctx context.Context, _ *emptypb.Empty) (*pb.GetOrdersResponse, error) {
	orders := make([]*pb.Order, 0, len(orderTotals)) // create slice for orders
	for drinkID, amount := range orderTotals {       // iterate over order totals
		orders = append(orders, &pb.Order{ // append each order
			DrinkId: drinkID, // set drink id
			Amount:  amount,  // set amount
		})
	}
	return &pb.GetOrdersResponse{Orders: orders}, nil // return all orders
}
