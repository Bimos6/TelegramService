package models

import "time"

type SessionState int

const (
	SessionStateInitializing SessionState = iota
	SessionStateWaitingQR
	SessionStateReady
	SessionStateError
	SessionStateClosed
)

type Message struct {
	ID        int64
	From      string
	Text      string
	Timestamp time.Time
	SessionID string
}
