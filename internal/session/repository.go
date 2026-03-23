package session

import "context"

type SessionRepository interface {
	Save(ctx context.Context, s *Session) error
	Find(ctx context.Context, id string) (*Session, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*Session, error)
	Count(ctx context.Context) (int, error)
}
