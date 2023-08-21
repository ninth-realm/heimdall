package store

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Session struct {
	Token     string    `json:"token" db:"token"`
	UserId    uuid.UUID `json:"userId" db:"user_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
}

type SessionRepository interface {
	GetSession(token string, opts QueryOptions) (Session, error)
	SaveSession(session Session, opts QueryOptions) error
}
