package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	pb "github.com/zombiemachines/mariadb-grpc-test/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type App struct {
	grpcConn   *grpc.ClientConn
	grpcClient pb.PersonServiceClient
	echo.Echo
}

var app *App

func main() {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize gRPC client
	creds, err := credentials.NewClientTLSFromFile("../tls/cert.pem", "")
	if err != nil {
		log.Fatalf("Failed to load TLS certificate: %v", err)
	}
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	// defer conn.Close()

	// Create a client for the PersonService
	client := pb.NewPersonServiceClient(conn)

	app = &App{grpcConn: conn, grpcClient: client, Echo: *echo.New()}
	defer app.grpcConn.Close()

	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// Routes
	// app.POST("/persons", createPersonHandler)
	// app.GET("/persons/:id", getPersonHandler)
	// app.PUT("/persons/:id", updatePersonHandler)
	app.DELETE("/persons/:id", deletePersonHandler)

	// Start server
	port := os.Getenv("ECHO_PORT")
	app.Logger.Fatal(app.Start(":" + port))
}

func deletePersonHandler(c echo.Context) error {
	// Parse ID from path parameter
	id := c.Param("id")
	tmp, _ := strconv.Atoi(id)
	nuID := int32(tmp)
	// Call gRPC method to delete person
	res, err := app.grpcClient.DeletePerson(context.Background(), &pb.DeletePersonByIDRequest{Id: nuID})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

// func createPersonHandler(c echo.Context) error {
// 	// Parse request body
// 	req := new(pb.CreatePersonRequest)
// 	if err := c.Bind(req); err != nil {
// 		return c.JSON(http.StatusBadRequest, err.Error())
// 	}

// 	// Call gRPC method to create person
// 	res, err := app.grpcClient.CreatePerson(context.Background(), req)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}

// 	return c.JSON(http.StatusOK, res)
// }

// func getPersonHandler(c echo.Context) error {
// 	// Parse ID from path parameter
// 	id := c.Param("id")

// 	tmp, _ := strconv.Atoi(id)
// 	nuID := int32(tmp)
// 	// Call gRPC method to get person by ID
// 	res, err := app.grpcClient.GetPersonByID(context.Background(), &pb.GetPersonByIDRequest{Id: nuID})
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}

// 	return c.JSON(http.StatusOK, res)
// }

// func updatePersonHandler(c echo.Context) error {
// 	// Parse ID from path parameter
// 	id := c.Param("id")

// 	// Parse request body
// 	req := new(pb.UpdatePersonRequest)
// 	if err := c.Bind(req); err != nil {
// 		return c.JSON(http.StatusBadRequest, err.Error())
// 	}

// 	// Call gRPC method to update person
// 	tmp, _ := strconv.Atoi(id)
// 	req.Id = int32(tmp)
// 	res, err := app.grpcClient.UpdatePerson(context.Background(), req)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}

// 	return c.JSON(http.StatusOK, res)
// }
