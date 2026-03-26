package session

import (
	"context"
	"fmt"
	"sync"

	"github.com/Bimos6/telegram-service/internal/models"
	"github.com/Bimos6/telegram-service/internal/telegram"
	"github.com/Bimos6/telegram-service/pkg/logger"
	"github.com/google/uuid"
)

type Session struct {
	id       string
	client   *telegram.MockClient
	state    models.SessionState
	stateMu  sync.RWMutex
	messages chan models.Message
	log      logger.Logger
	cancel   context.CancelFunc
}

func NewSession(appID int, appHash string, log logger.Logger) *Session {
	id := uuid.New().String()
	return &Session{
		id:       id,
		client:   telegram.NewMockClient(id, log),
		state:    models.SessionStateInitializing,
		messages: make(chan models.Message, 100),
		log:      log.WithField("component", "session").WithField("session_id", id),
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
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	s.state = state
}

func (s *Session) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	s.log.Info("Starting session")
	s.log.WithField("state_before", s.State()).Debug("State before Connect")

	if err := s.client.Connect(ctx); err != nil {
		s.log.WithField("error", err).Error("Failed to connect")
		return err
	}

	s.log.Debug("Setting state to Ready")
	s.setState(models.SessionStateReady)

	currentState := s.State()
	s.log.WithField("state_after", currentState).Info("Session started")

	go func() {
		for msg := range s.client.Messages() {
			select {
			case s.messages <- msg:
			default:
				s.log.Warn("Message channel full")
			}
		}
	}()

	return nil
}

func (s *Session) Stop() error {
	s.log.Info("Stopping session")
	s.setState(models.SessionStateClosed)
	if s.cancel != nil {
		s.cancel()
	}
	s.client.Disconnect()
	close(s.messages)
	return nil
}

func (s *Session) SendMessage(ctx context.Context, text string) (int64, error) {
	if s.State() != models.SessionStateReady {
		return 0, fmt.Errorf("session not ready")
	}
	return s.client.SendMessage(ctx, "@me", text)
}

func (s *Session) Messages() <-chan models.Message {
	return s.messages
}

func (s *Session) GetQRCode() (string, error) {
	//s.setState(models.SessionStateWaitingQR)
	return s.client.GetQRCode(context.Background())
}
