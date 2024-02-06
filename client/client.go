package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/zombiemachines/mariadb-grpc-test/grpc"
)

var DELETER bool = true

func main() {
	creds, err := credentials.NewClientTLSFromFile("../tls/cert.pem", "")
	if err != nil {
		log.Fatalf("Failed to load TLS certificate: %v", err)
	}
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create a client for the PersonService
	client := pb.NewPersonServiceClient(conn)

	if !DELETER { // Test the CreatePerson RPC
		createPersonResponse, err := client.CreatePerson(context.Background(), &pb.CreatePersonRequest{Name: "Conner", Age: 22})
		if err != nil {
			log.Fatalf("CreatePerson RPC failed: %v", err)
		}
		log.Printf("CreatePerson response: %v", createPersonResponse)

		// Test the GetPersonByID RPC
		getPersonByIDResponse, err := client.GetPersonByID(context.Background(), &pb.GetPersonByIDRequest{Id: createPersonResponse.Id})
		if err != nil {
			log.Fatalf("GetPersonByID RPC failed: %v", err)
		}
		log.Printf("GetPersonByID response: %v", getPersonByIDResponse)

		// Test the UpdatePerson RPC
		updatePersonResponse, err := client.UpdatePerson(context.Background(), &pb.UpdatePersonRequest{Id: createPersonResponse.Id, Name: "Conner Updated", Age: 33})
		if err != nil {
			log.Fatalf("UpdatePerson RPC failed: %v", err)
		}
		log.Printf("UpdatePerson response: %v", updatePersonResponse) //updatePersonResponse.Success
	} else if DELETER {
		toDeleteID := 5
		deletePersonByIDResponse, err := client.DeletePerson(context.Background(), &pb.DeletePersonByIDRequest{Id: int32(toDeleteID)})
		if err != nil {
			log.Fatalf("DeletePerson RPC failed: %v", err)
		}
		log.Printf("DeletePerson response: %v \n Person with ID: %d is deleted", deletePersonByIDResponse, toDeleteID)
	}
}
