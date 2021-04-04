package main

import (
	"context"
	"fmt"
	"grpc-api/proto"
	"io"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CalculatorService struct {
	mutex sync.Mutex
}

func (s *CalculatorService) FindMaximum(stream proto.CalculatorService_FindMaximumServer) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var maximum int32
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		if request.Number == 0 {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("error: number=%d is invalid", request.GetNumber()))
		}

		if maximum < request.Number {
			maximum = request.Number
			err = stream.Send(&proto.FindMaximumResponse{
				Maximum: maximum,
			})

			if err != nil {
				return err
			}
		}
	}
}

func (s *CalculatorService) Sum(ctx context.Context, request *proto.SumRequest) (*proto.SumResponse, error) {
	return &proto.SumResponse{
		Sum: request.Number1 + request.Number2,
	}, nil
}

func (s *CalculatorService) ComputeAverage(stream proto.CalculatorService_ComputeAverageServer) error {
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

func (s *CalculatorService) PrimeNumberDecomposition(request *proto.PrimeNumberDecompositionRequest, stream proto.CalculatorService_PrimeNumberDecompositionServer) error {
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
