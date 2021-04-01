package main

import (
	"context"
	"fmt"
	"grpc-api/proto"
	"io"
	"net"

	"google.golang.org/grpc"
)

type CalculatorService struct{}

func (s *CalculatorService) Sum(ctx context.Context, request *proto.SumRequest) (*proto.SumResponse, error) {
	return &proto.SumResponse{
		Sum: request.Number1 + request.Number2,
	}, nil
}

func (s CalculatorService) ComputeAverage(stream proto.CalculatorService_ComputeAverageServer) error {
	fmt.Println("service start calculator...")
	var sum, count int32
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			avg := float32(sum) / float32(count)
			response := &proto.ComputeAverageResponse{
				Average: avg,
			}

			return stream.SendAndClose(response)
		}

		if err != nil {
			return err
		}

		sum += request.GetNumber()
		count++
	}
}

func (s CalculatorService) PrimeNumberDecomposition(request *proto.PrimeNumberDecompositionRequest, stream proto.CalculatorService_PrimeNumberDecompositionServer) error {
	number := request.GetNumber()
	var divisor int32 = 2
	for number > 1 {
		if number%divisor == 0 {
			response := proto.PrimeNumberDecompositionResponse{PrimeNumber: divisor}
			err := stream.Send(&response)
			if err != nil {
				return err
			}

			number = number / divisor
		} else {
			divisor++
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":8000")

	grpcServer := grpc.NewServer()
	calculatorSrv := CalculatorService{}
	proto.RegisterCalculatorServiceServer(grpcServer, &calculatorSrv)

	// Start grpcServer
	if err = grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
