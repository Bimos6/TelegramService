package session

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Bimos6/telegram-service/internal/models"
	"github.com/Bimos6/telegram-service/pkg/logger"
	"github.com/google/uuid"
)

type Session struct {
	id       string
	state    models.SessionState
	stateMu  sync.RWMutex
	messages chan models.Message
	log      logger.Logger
	cancel   context.CancelFunc
}

func NewSession(log logger.Logger) *Session {
	return &Session{
		id:       uuid.New().String(),
		state:    models.SessionStateInitializing,
		messages: make(chan models.Message, 100),
		log:      log.WithField("component", "session"),
	}
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) State() models.SessionState {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.state
}

func (s *Session) setState(state models.SessionState) {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	s.state = state
}

func (s *Session) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	s.log.WithField("session_id", s.id).Info("Starting session")
	s.setState(models.SessionStateReady)

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		msgID := int64(1)
		for {
			select {
			case <-ctx.Done():
				s.log.WithField("session_id", s.id).Debug("Message generator stopped")
				return
			case <-ticker.C:
				s.messages <- models.Message{
					ID:        msgID,
					From:      "test_bot",
					Text:      fmt.Sprintf("Test message #%d", msgID),
					Timestamp: time.Now(),
					SessionID: s.id,
				}
				msgID++
			}
		}
	}()

	return nil
}

func (s *Session) Stop() error {
	s.log.WithField("session_id", s.id).Info("Stopping session")
	s.setState(models.SessionStateClosed)
	if s.cancel != nil {
		s.cancel()
	}
	close(s.messages)
	return nil
}

func (s *Session) SendMessage(ctx context.Context, text string) (int64, error) {
	if s.State() != models.SessionStateReady {
		return 0, fmt.Errorf("session not ready")
	}

	msgID := time.Now().UnixNano()
	s.log.WithField("session_id", s.id).WithField("text", text).Info("Sending message")
	return msgID, nil
}

func (s *Session) Messages() <-chan models.Message {
	return s.messages
}

func (s *Session) GetQRCode() (string, error) {
	s.setState(models.SessionStateWaitingQR)
	return "tg://login?token=mock_token", nil
}
