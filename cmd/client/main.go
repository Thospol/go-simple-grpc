package main

import (
	"context"
	"encoding/json"
	"fmt"
	"grpc-api/proto"
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

	PrintFormatJSON(response)

	return response, err
}

// GetAvgCalculator get avg from slice of number
func GetAvgCalculator(client proto.CalculatorServiceClient) (float32, error) {
	// Timeout 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.ComputeAverage(ctx)
	if err != nil {
		return 0, err
	}

	request := []*proto.ComputeAverageRequest{
		{Number: 1}, {Number: 2}, {Number: 3}, {Number: 4},
	}

	for i := range request {
		err = stream.Send(request[i])
		if err != nil {
			return 0, err
		}
		fmt.Printf("send request: %+v\n", request[i])
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		return 0, err
	}

	return response.Average, err
}

func main() {
	// Disable transport security for this example only,
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial("127.0.0.1:8000", opts...)
	if err != nil {
		panic(err)
	}

	// clientUser := proto.NewUserServiceClient(conn)
	// _, err = GetProfile(clientUser)
	// if err != nil {
	// 	panic(err)
	// }

	clientCalculator := proto.NewCalculatorServiceClient(conn)
	avg, err := GetAvgCalculator(clientCalculator)
	if err != nil {
		panic(err)
	}

	fmt.Println("avg: ", avg)
}

// PrintFormatJSON print format json
func PrintFormatJSON(n interface{}) {
	b, _ := json.MarshalIndent(n, "", "\t")
	os.Stdout.Write(b)
}
