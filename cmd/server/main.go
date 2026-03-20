package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/Bimos6/telegram-service/internal/app/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedTelegramServiceServer
}

func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterTelegramServiceServer(s, &server{})
	reflection.Register(s)

	fmt.Println("Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
