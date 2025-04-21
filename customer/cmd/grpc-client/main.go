package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/DuongVu089x/interview/customer/proto/customer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	addr    = flag.String("addr", "localhost:50051", "the address to connect to")
	userID  = flag.String("user", "", "user ID to query")
	timeout = flag.Duration("timeout", 5*time.Second, "timeout for the request")
)

func main() {
	flag.Parse()

	if *userID == "" {
		log.Fatal("user ID is required")
	}

	// Set up a connection to the server
	log.Printf("Connecting to gRPC server at %s...", *addr)
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, *addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // Make the connection attempt blocking
	)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected successfully to %s", *addr)

	// Create a new client
	client := pb.NewCustomerServiceClient(conn)

	// Check server health
	log.Printf("Checking server health...")
	healthResp, err := client.Check(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			log.Fatalf("Health check failed with gRPC error: code=%v message=%v", st.Code(), st.Message())
		} else {
			log.Fatalf("Health check failed: %v", err)
		}
	}

	if healthResp.Status != pb.HealthCheckResponse_SERVING {
		log.Fatalf("Server is not healthy: %s", healthResp.Error)
	}
	log.Printf("Server is healthy")

	// Make the request
	log.Printf("Requesting customer with ID: %s", *userID)
	response, err := client.GetCustomer(ctx, &pb.GetCustomerRequest{
		UserId: *userID,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			// This is a gRPC error
			log.Fatalf("gRPC error: code=%v message=%v", st.Code(), st.Message())
		} else {
			// This is a regular error
			log.Fatalf("Error getting customer: %v", err)
		}
	}

	// Print the response
	if !response.Exists {
		fmt.Printf("\nCustomer with ID %s not found\n", *userID)
		return
	}

	customer := response.Customer
	fmt.Printf("\nCustomer found:\n")
	fmt.Printf("  ID: %s\n", customer.Id)
	fmt.Printf("  Name: %s\n", customer.Name)
	fmt.Printf("  Email: %s\n", customer.Email)
	fmt.Printf("  Phone: %s\n", customer.Phone)
	fmt.Printf("  Created At: %s\n", customer.CreatedAt)
	fmt.Printf("  Updated At: %s\n", customer.UpdatedAt)
}
