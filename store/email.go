package store

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Email struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewEmail struct {
	UserID uuid.UUID
	Email  string
}

type EmailRepository interface {
	InsertEmail(email NewEmail, opts QueryOptions) (uuid.UUID, error)
}
