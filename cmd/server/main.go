package main

import (
	"context"
	"fmt"
	"net"
	"os"

	pb "github.com/Bimos6/telegram-service/internal/app/grpc/proto"
	"github.com/Bimos6/telegram-service/internal/config"
	"github.com/Bimos6/telegram-service/internal/repository"
	"github.com/Bimos6/telegram-service/internal/session"
	"github.com/Bimos6/telegram-service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedTelegramServiceServer
	log     logger.Logger
	manager *session.Manager
}

func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	s.log.Info("Ping received")
	return &pb.PingResponse{Message: "pong"}, nil
}

func (s *server) CreateTestSession(ctx context.Context, req *pb.CreateTestSessionRequest) (*pb.CreateTestSessionResponse, error) {
	sessionID, err := s.manager.CreateSession(ctx)
	if err != nil {
		return nil, err
	}

	sess, _ := s.manager.GetSession(ctx, sessionID)

	go func() {
		for msg := range sess.Messages() {
			s.log.WithField("message", msg.Text).Info("Test message received")
		}
	}()

	return &pb.CreateTestSessionResponse{
		SessionId: sessionID,
		Message:   "Session created",
	}, nil
}

func (s *server) ListSessions(ctx context.Context, req *pb.ListSessionsRequest) (*pb.ListSessionsResponse, error) {
	sessions, _ := s.manager.ListSessions(ctx)
	ids := make([]string, len(sessions))
	for i, sess := range sessions {
		ids[i] = sess.ID()
	}
	return &pb.ListSessionsResponse{SessionIds: ids}, nil
}

func main() {
	cfg := config.Load()
	log := logger.Init(cfg.LogLevel)

	repo := repository.NewRepository(log)
	manager := session.NewManager(repo, cfg, log)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.WithField("error", err).Error("Failed to listen")
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterTelegramServiceServer(s, &server{
		log:     log,
		manager: manager,
	})

	reflection.Register(s)

	log.WithField("port", cfg.GRPCPort).Info("Server starting")
	if err := s.Serve(lis); err != nil {
		log.WithField("error", err).Error("Server failed")
		os.Exit(1)
	}
}
