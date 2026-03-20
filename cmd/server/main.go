package main

import (
	"fmt"
	"net"
	"os"

	pb "github.com/Bimos6/telegram-service/internal/app/grpc/proto"

	"github.com/Bimos6/telegram-service/internal/config"
	"github.com/Bimos6/telegram-service/pkg/logger"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTelegramServiceServer
}

func main() {
	cfg := config.Load()

	logger.Init("debug")
	log := logger.Get()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterTelegramServiceServer(s, &server{})

	fmt.Println("Server listening on :", cfg.GRPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
