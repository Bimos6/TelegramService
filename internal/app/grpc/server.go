package grpc

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Bimos6/telegram-service/internal/app/grpc/proto"
	"github.com/Bimos6/telegram-service/internal/session"
	"github.com/Bimos6/telegram-service/pkg/logger"
)

type Server struct {
	pb.UnimplementedTelegramServiceServer
	manager *session.Manager
	log     logger.Logger
}

func NewServer(manager *session.Manager, log logger.Logger) *Server {
	return &Server{
		manager: manager,
		log:     log.WithField("component", "grpc"),
	}
}

func (s *Server) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	s.log.Info("CreateSession called")

	sessionID, err := s.manager.CreateSession(ctx)
	if err != nil {
		s.log.WithField("error", err).Error("Failed to create session")
		return nil, status.Errorf(codes.Internal, "failed: %v", err)
	}

	sess, _ := s.manager.GetSession(ctx, sessionID)
	qrCode, _ := sess.GetQRCode()

	return &pb.CreateSessionResponse{
		SessionId: sessionID,
		QrCode:    qrCode,
	}, nil
}

func (s *Server) DeleteSession(ctx context.Context, req *pb.DeleteSessionRequest) (*pb.DeleteSessionResponse, error) {
	s.log.WithField("session_id", req.SessionId).Info("DeleteSession called")

	if err := s.manager.DeleteSession(ctx, req.SessionId); err != nil {
		return nil, status.Errorf(codes.NotFound, "session not found")
	}

	return &pb.DeleteSessionResponse{}, nil
}

func (s *Server) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	s.log.WithField("session_id", req.SessionId).WithField("peer", req.Peer).Info("SendMessage called")

	sess, err := s.manager.GetSession(ctx, req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "session not found")
	}

	msgID, err := sess.SendMessage(ctx, req.Text)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed: %v", err)
	}

	return &pb.SendMessageResponse{MessageId: strconv.FormatInt(msgID, 10)}, nil
}

func (s *Server) SubscribeMessages(req *pb.SubscribeMessagesRequest, stream pb.TelegramService_SubscribeMessagesServer) error {
	s.log.WithField("session_id", req.SessionId).Info("SubscribeMessages called")

	sess, err := s.manager.GetSession(stream.Context(), req.SessionId)
	if err != nil {
		return status.Errorf(codes.NotFound, "session not found")
	}

	for {
		select {
		case <-stream.Context().Done():
			s.log.WithField("session_id", req.SessionId).Info("Subscription ended")
			return nil
		case msg, ok := <-sess.Messages():
			if !ok {
				return nil
			}
			if err := stream.Send(&pb.MessageUpdate{
				MessageId: msg.ID,
				From:      msg.From,
				Text:      msg.Text,
				Timestamp: msg.Timestamp.Unix(),
			}); err != nil {
				return err
			}
		}
	}
}
