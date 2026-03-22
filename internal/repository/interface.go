package repository

import (
	"context"

	"github.com/Bimos6/telegram-service/internal/session"
)

type SessionRepository interface {
	Save(ctx context.Context, s *session.Session) error
	Find(ctx context.Context, id string) (*session.Session, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*session.Session, error)
	Count(ctx context.Context) (int, error)
}
