package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/ninth-realm/heimdall/crypto"
	"github.com/ninth-realm/heimdall/store"
)

type Service struct {
	Repo store.Repository

	JWTSettings JWTSettings
}

func (s Service) Login(ctx context.Context, username, password string) (Token, error) {
	return store.RunUnitOfWork(ctx, s.Repo, func(tx *sqlx.Tx) (Token, error) {
		opts := store.QueryOptions{Ctx: ctx, Txn: tx}

		user, err := s.Repo.GetUserByEmail(username, opts)
		if err != nil {
			return Token{}, err
		}

		hash, err := s.Repo.GetUserPasswordHash(user.ID, opts)
		if err != nil {
			return Token{}, err
		}

		correctPassword, err := crypto.ValidatePassword(password, hash)
		if err != nil {
			return Token{}, err
		} else if !correctPassword {
			return Token{}, errors.New("incorrect password")
		}

		token, err := crypto.GenerateRandBase64String(32)
		if err != nil {
			return Token{}, err
		}

		now := time.Now()
		lifespan := 24 * time.Hour
		err = s.Repo.SaveSession(store.Session{
			Token:     token,
			UserId:    user.ID,
			CreatedAt: now,
			ExpiresAt: now.Add(lifespan),
		}, opts)
		if err != nil {
			return Token{}, err
		}

		return Token{AccessToken: token, Lifespan: int(lifespan.Seconds())}, nil
	})
}

func (s Service) IntrospectToken(ctx context.Context, token string) (TokenInfo, error) {
	return store.RunUnitOfWork(ctx, s.Repo, func(tx *sqlx.Tx) (TokenInfo, error) {
		opts := store.QueryOptions{Ctx: ctx, Txn: tx}

		session, err := s.Repo.GetSession(token, opts)
		if err != nil {
			return TokenInfo{}, err
		}

		return TokenInfo{
			Active:    true,
			ExpiresAt: int(session.ExpiresAt.Sub(session.CreatedAt).Seconds()),
			UserID:    session.UserId.String(),
		}, nil
	})
}

func (s Service) ValidateAPIKey(ctx context.Context, key string) error {
	clientIDStr, token, found := strings.Cut(key, ":")
	if !found {
		return errors.New("malformed API key")
	}

	clientID, err := uuid.FromString(clientIDStr)
	if !found {
		return errors.New("invalid client ID")
	}

	prefix, suffix, found := strings.Cut(token, ".")
	if !found {
		return errors.New("malformed API key")
	}

	k, err := s.Repo.GetClientAPIKey(clientID, prefix, store.QueryOptions{Ctx: ctx})
	if err != nil {
		return err
	}

	ok, err := crypto.ValidatePassword(suffix, k.Hash)
	if err != nil {
		return err
	} else if !ok {
		return errors.New("invalid API key")
	}

	return nil
}
