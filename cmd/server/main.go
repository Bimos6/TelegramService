package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	telegrpc "github.com/Bimos6/telegram-service/internal/app/grpc"
	pb "github.com/Bimos6/telegram-service/internal/app/grpc/proto"
	"github.com/Bimos6/telegram-service/internal/config"
	"github.com/Bimos6/telegram-service/internal/repository"
	"github.com/Bimos6/telegram-service/internal/session"
	"github.com/Bimos6/telegram-service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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

	grpcServer := grpc.NewServer()
	pb.RegisterTelegramServiceServer(grpcServer, telegrpc.NewServer(manager, log))

	reflection.Register(grpcServer)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.WithField("port", cfg.GRPCPort).Info("Server starting")
		if err := grpcServer.Serve(lis); err != nil {
			log.WithField("error", err).Error("Server failed")
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	manager.StopAll(shutdownCtx)

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Info("Server stopped")
	case <-time.After(5 * time.Second):
		log.Warn("Force stopping")
		grpcServer.Stop()
	}
}
