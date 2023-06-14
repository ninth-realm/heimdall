package store

import (
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Client struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Enabled   bool      `json:"enabled" db:"enabled"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type NewClient struct {
	Name    string
	Enabled bool
}

type ClientPatch struct {
	Name    *string
	Enabled *bool
}

func (p ClientPatch) ApplyTo(client Client) Client {
	if p.Name != nil {
		client.Name = strings.TrimSpace(*p.Name)
	}

	if p.Enabled != nil {
		client.Enabled = *p.Enabled
	}

	return client
}

type ClientRepository interface {
	ListClients(opts QueryOptions) ([]Client, error)
	GetClientById(id uuid.UUID, opts QueryOptions) (Client, error)
	InsertClient(user NewClient, opts QueryOptions) (uuid.UUID, error)
	SaveClient(user Client, opts QueryOptions) error
	DeleteClient(id uuid.UUID, opts QueryOptions) error
}
