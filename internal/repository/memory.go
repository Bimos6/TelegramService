package repository

import (
	"context"
	"fmt"

	//"fmt"
	"sync"

	"github.com/Bimos6/telegram-service/internal/session"
	"github.com/Bimos6/telegram-service/pkg/logger"
)

type Repository struct {
	sessions map[string]*session.Session
	mu       sync.RWMutex
	log      logger.Logger
}

func NewRepository(log logger.Logger) *Repository {
	return &Repository{
		sessions: make(map[string]*session.Session),
		log:      log.WithField("component", "repository"),
	}
}

func (r *Repository) Save(ctx context.Context, s *session.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[s.ID()] = s
	r.log.WithField("session_id", s.ID()).Debug("Session saved")
	return nil
}

func (r *Repository) Find(ctx context.Context, id string) (*session.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session not found: %s", id)
	}
	return s, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, id)
	r.log.WithField("session_id", id).Debug("Session deleted")
	return nil
}

func (r *Repository) List(ctx context.Context) ([]*session.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sessions := make([]*session.Session, 0, len(r.sessions))
	for _, s := range r.sessions {
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.sessions), nil
}
