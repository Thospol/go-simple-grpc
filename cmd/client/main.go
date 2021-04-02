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
)

// GetAvgCalculator get avg from slice of number
func SumValue(client proto.CalculatorServiceClient) error {
	// Timeout 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := proto.SumRequest{
		Number1: 3,
		Number2: 10,
	}
	response, err := client.Sum(ctx, &request)
	if err != nil {
		return err
	}

	fmt.Println("sum: ", response.Sum)

	return err
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

	var result int32 = 0
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			// handle log
			return err
		}

		fmt.Println("primeNumber: ", response.PrimeNumber)
		if result > 0 {
			result = result * response.PrimeNumber
		} else {
			result = response.PrimeNumber
		}
	}
	fmt.Println("number: ", request.Number)
	fmt.Println("result: ", result)

	return nil
}

func FindMaximum(client proto.CalculatorServiceClient) error {
	stream, err := client.FindMaximum(context.Background())
	if err != nil {
		return err
	}

	request := []*proto.FindMaximumRequest{
		{
			Number: 1,
		},
		{
			Number: 5,
		},
		{
			Number: 3,
		},
		{
			Number: 6,
		},
		{
			Number: 2,
		},
		{
			Number: 20,
		},
	}

	go func() {
		for i := 0; i < len(request); i++ {
			err = stream.Send(request[i])
			if err != nil {
				return
			}
			fmt.Println("request number: ", request[i].Number)
			time.Sleep(1 * time.Second)
		}

		err = stream.CloseSend()
		if err != nil {
			return
		}
	}()

	wait := make(chan struct{})
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(wait)
				return
			}

			if err != nil {
				close(wait)
				return
			}
			fmt.Println("response maximum number: ", resp.Maximum)
		}
	}()

	<-wait
	fmt.Println("finish")

	return nil
}

func main() {
	// Disable transport security for this example only,
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial("127.0.0.1:8000", opts...)
	if err != nil {
		panic(err)
	}

	clientCalculator := proto.NewCalculatorServiceClient(conn)
	// client streaming
	err = SumValue(clientCalculator)
	if err != nil {
		panic(err)
	}

	fmt.Println("============================================")

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

	fmt.Println("============================================")

	// bi-directional streaming
	err = FindMaximum(clientCalculator)
	if err != nil {
		panic(err)
	}
}

// PrintFormatJSON print format json
func PrintFormatJSON(n interface{}) {
	b, _ := json.MarshalIndent(n, "", "\t")
	os.Stdout.Write(b)
}
