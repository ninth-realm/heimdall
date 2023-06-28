package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ninth-realm/heimdall/crypto"
	"github.com/ninth-realm/heimdall/store"
)

type Service struct {
	Repo store.Repository

	JWTSettings JWTSettings
}

func (s Service) Login(ctx context.Context, username, password string) (Token, error) {
	user, err := s.Repo.GetUserByEmail(username, store.QueryOptions{Ctx: ctx})
	if err != nil {
		return Token{}, err
	}

	hash, err := s.Repo.GetUserPasswordHash(user.ID, store.QueryOptions{Ctx: ctx})
	if err != nil {
		return Token{}, err
	}

	correctPassword, err := crypto.ValidatePassword(password, hash)
	if err != nil {
		return Token{}, err
	} else if !correctPassword {
		return Token{}, errors.New("incorrect password")
	}

	token, err := generateJWT(user, s.JWTSettings)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func (s Service) IntrospectToken(ctx context.Context, token string) (TokenInfo, error) {
	t, err := validateJWT(token, s.JWTSettings)
	if err != nil {
		return TokenInfo{}, err
	}

	claims := t.Claims.(jwt.MapClaims)

	return TokenInfo{
		Active:    true,
		ExpiresAt: int(claims["exp"].(float64)),
		UserID:    claims["sub"].(string),
	}, nil
}
