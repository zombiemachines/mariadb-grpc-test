package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	pb "github.com/zombiemachines/mariadb-grpc-test/grpc"
)

type server struct {
	pb.UnimplementedPersonServiceServer
	db *gorm.DB
}

type Person struct {
	ID   int `gorm:"primaryKey"`
	Name string
	Age  int
}

func (s *server) CreatePerson(ctx context.Context, req *pb.CreatePersonRequest) (*pb.CreatePersonResponse, error) {
	person := &Person{Name: req.Name, Age: int(req.Age)}
	if err := s.db.Create(person).Error; err != nil {
		return nil, err
	}
	return &pb.CreatePersonResponse{Id: int32(person.ID)}, nil
}

func (s *server) GetPersonByID(ctx context.Context, req *pb.GetPersonByIDRequest) (*pb.GetPersonByIDResponse, error) {
	var person Person
	if err := s.db.First(&person, int(req.Id)).Error; err != nil {
		return nil, err
	}
	return &pb.GetPersonByIDResponse{Name: person.Name, Age: int32(person.Age)}, nil
}

func (s *server) UpdatePerson(ctx context.Context, req *pb.UpdatePersonRequest) (*pb.UpdatePersonResponse, error) {
	var person Person
	if err := s.db.First(&person, int(req.Id)).Error; err != nil {
		return nil, err
	}
	person.Name = req.Name
	person.Age = int(req.Age)
	if err := s.db.Save(&person).Error; err != nil {
		return nil, err
	}
	//now that updating the record is successful set he Success boolean to true
	return &pb.UpdatePersonResponse{Success: true}, nil
}

func (s *server) DeletePerson(ctx context.Context, req *pb.DeletePersonByIDRequest) (*pb.DeletePersonByIDResponse, error) {
	var person Person
	if err := s.db.Delete(&person, int(req.Id)).Error; err != nil {
		log.Printf("Failed to delete person with ID %d: %v", req.Id, err)
		return &pb.DeletePersonByIDResponse{Success: false}, err
	}
	// Deletion successful, return success response
	return &pb.DeletePersonByIDResponse{Success: true}, nil
}

func main() {
	// Connect to the database
	dsn := "root:root@tcp(192.168.1.200:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&Person{}); err != nil {
		log.Fatal("Failed to auto-migrate schema:", err)
	}

	// Create gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile("../tls/cert.pem", "../tls/key.pem")
	if err != nil {
		log.Fatalf("Failed to load TLS keys: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterPersonServiceServer(grpcServer, &server{db: db})
	log.Println("gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
