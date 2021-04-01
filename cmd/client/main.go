package main

import (
	"context"
	"encoding/json"
	"fmt"
	"grpc-api/proto"
	"io"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetProfile get profile by user
func GetProfile(client proto.UserServiceClient) (*proto.User, error) {
	// Timeout 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := proto.UserRequest{
		Id: "64548",
	}
	response, err := client.GetProfile(ctx, &request)
	statusCode := status.Code(err)
	if statusCode != codes.OK {
		return nil, err
	}

	fmt.Println("user: ", response)

	return response, err
}

// GetAvgCalculator get avg from slice of number
func GetAvgCalculator(client proto.CalculatorServiceClient) error {
	// Timeout 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.ComputeAverage(ctx)
	if err != nil {
		return err
	}

	request := []*proto.ComputeAverageRequest{
		{Number: 1}, {Number: 2}, {Number: 3}, {Number: 4},
	}

	for i := range request {
		err = stream.Send(request[i])
		if err != nil {
			return err
		}
		fmt.Printf("send request: %+v\n", request[i])
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	fmt.Println("avg: ", response.Average)

	return err
}

func GetPrimeNumberDecomposition(client proto.CalculatorServiceClient) error {
	// Timeout 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := proto.PrimeNumberDecompositionRequest{
		Number: 210,
	}
	stream, err := client.PrimeNumberDecomposition(ctx, &request)
	if err != nil {
		return err
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		fmt.Println("primeNumber: ", response.PrimeNumber)
	}

	return nil
}

func main() {
	// Disable transport security for this example only,
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial("127.0.0.1:8000", opts...)
	if err != nil {
		panic(err)
	}

	clientUser := proto.NewUserServiceClient(conn)
	// unary request/ response
	_, err = GetProfile(clientUser)
	if err != nil {
		panic(err)
	}
	fmt.Println("============================================")

	clientCalculator := proto.NewCalculatorServiceClient(conn)
	// client streaming
	err = GetAvgCalculator(clientCalculator)
	if err != nil {
		panic(err)
	}

	fmt.Println("============================================")

	// server streaming
	err = GetPrimeNumberDecomposition(clientCalculator)
	if err != nil {
		panic(err)
	}
}

// PrintFormatJSON print format json
func PrintFormatJSON(n interface{}) {
	b, _ := json.MarshalIndent(n, "", "\t")
	os.Stdout.Write(b)
}
