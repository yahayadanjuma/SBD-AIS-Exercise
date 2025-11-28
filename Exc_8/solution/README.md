# Exercise 8 Solution: gRPC Drink Order System

## Overview
This solution demonstrates how to build a simple gRPC-based client-server application in Go for ordering drinks. The project is split into three main parts:
- **Protobuf Definition**: Defines the service and messages for drink ordering.
- **Server**: Implements the gRPC service, manages drink data, and handles orders.
- **Client**: Connects to the server, lists drinks, places orders, and retrieves totals.




## How the Solution Works

### 1. Protobuf Definition (pb/orders.proto)
- Describes the OrderService with three RPCs:
  - `GetDrinks`: Returns a list of available drinks.
  - `OrderDrink`: Places an order for a drink.
  - `GetOrders`: Returns the total orders for each drink.
- The `.proto` file is compiled using `protoc` to generate Go code for client and server communication.

### 2. Server (`server/grpc_server.go`)
- Implements the `OrderService` defined in the proto file.
- Stores a list of drinks and keeps track of orders in memory.
- Handles requests from the client:
  - Returns available drinks.
  - Validates and processes drink orders.
  - Returns the total number of drinks ordered.
- Starts a gRPC server on port 4000.

### 3. Client (`client/grpc_client.go`)
- Connects to the gRPC server at `localhost:4000`.
- Requests the list of drinks from the server and displays them.
- Places two rounds of orders for each drink (first 2, then 6 of each).
- Requests the total orders and displays the results.

### 4. Main (`main.go`)
- Starts the server in a background goroutine.
- Waits briefly to ensure the server is running.
- Creates and runs the client to interact with the server.
- Prints a completion message when all orders are done.

---

## How to Run the Solution
1. **Generate Protobuf Code**
   - Run the provided script (e.g., `generate_pb.sh`) to generate Go files from the proto definition.
2. **Build and Run**
   - Use `go run .` or `go build -o exc8main main.go` then `./exc8main` to start the application.
   - The server will start, and the client will automatically connect, place orders, and print results.



