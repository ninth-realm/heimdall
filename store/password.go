package store

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Password struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewPassword struct {
	UserID uuid.UUID
	Hash   string
}

type PasswordRepository interface {
	InsertPassword(password NewPassword, opts QueryOptions) (uuid.UUID, error)
}
