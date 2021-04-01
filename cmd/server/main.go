package main

import (
	"context"
	"fmt"
	"grpc-api/proto"
	"io"
	"net"

	"google.golang.org/grpc"
)

// UserServerImpl Implement proto.PingPongServer
type UserServerImpl struct {
}

// GetProfile get user service
func (s *UserServerImpl) GetProfile(ctx context.Context, request *proto.UserRequest) (*proto.User, error) {
	fmt.Println("Ping Received")

	resp := proto.User{
		ApplicantId: request.Id,
		ThName:      "ธวัชชัย",
		ThSurname:   "รัตนมงคล",
		ThNickname:  "วัช",
	}

	return &resp, nil
}

type CalculatorService struct {
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

func main() {
	// args := os.Args
	// fmt.Printf("args: %v", args[1:])

	userSrv := UserServerImpl{}
	calculatorSrv := CalculatorService{}

	lis, err := net.Listen("tcp", ":8000")

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, &userSrv)
	proto.RegisterCalculatorServiceServer(grpcServer, &calculatorSrv)

	// Start grpcServer
	if err = grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
