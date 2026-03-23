package session

import (
	"context"
	"fmt"

	"github.com/Bimos6/telegram-service/internal/config"
	"github.com/Bimos6/telegram-service/pkg/logger"
)

type Manager struct {
	repo SessionRepository
	cfg  *config.Config
	log  logger.Logger
}

func NewManager(repo SessionRepository, cfg *config.Config, log logger.Logger) *Manager {
	return &Manager{
		repo: repo,
		cfg:  cfg,
		log:  log.WithField("component", "manager"),
	}
}

func (m *Manager) CreateSession(ctx context.Context) (string, error) {
	count, _ := m.repo.Count(ctx)
	if count >= m.cfg.MaxSessions {
		return "", fmt.Errorf("max sessions limit reached: %d", m.cfg.MaxSessions)
	}

	sess := NewSession(m.log)
	if err := sess.Start(ctx); err != nil {
		return "", err
	}

	if err := m.repo.Save(ctx, sess); err != nil {
		return "", err
	}

	m.log.WithField("session_id", sess.ID()).Info("Session created")
	return sess.ID(), nil
}

func (m *Manager) GetSession(ctx context.Context, id string) (*Session, error) {
	return m.repo.Find(ctx, id)
}

func (m *Manager) DeleteSession(ctx context.Context, id string) error {
	sess, err := m.repo.Find(ctx, id)
	if err != nil {
		return err
	}
	sess.Stop()
	return m.repo.Delete(ctx, id)
}

func (m *Manager) ListSessions(ctx context.Context) ([]*Session, error) {
	sessions, err := m.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	return sessions, nil
}

func (m *Manager) StopAll(ctx context.Context) error {
	sessions, _ := m.repo.List(ctx)
	for _, sess := range sessions {
		sess.Stop()
	}
	return nil
}
