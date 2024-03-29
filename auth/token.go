package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ninth-realm/heimdall/store"
)

// Token holds the information required for transmitting the JWT to the client.
type Token struct {
	AccessToken string `json:"accessToken"`
	Lifespan    int    `json:"lifespan"`
}

type TokenInfo struct {
	Active    bool   `json:"active"`
	UserID    string `json:"sub,omitempty"`
	Username  string `json:"username,omitempty"`
	ExpiresAt int    `json:"exp,omitempty"`
}

type signingAlgorithm string

// The valid JWT hashing function algorithms.
const (
	HMAC256Algorithm signingAlgorithm = "HS256"
)

func (a signingAlgorithm) isValid() bool {
	return a == HMAC256Algorithm
}

// JWTSettings are the available configuration values for generating JWTs.
type JWTSettings struct {
	Issuer     string
	Lifespan   int
	SigningKey string
	Algorithm  signingAlgorithm
}

func (s JWTSettings) validate() error {
	if strings.TrimSpace(s.Issuer) == "" {
		return errors.New("JWT issuer cannot be empty")
	}

	if s.Lifespan <= 0 {
		return errors.New("JWT lifetime must be a positive integer")
	}

	if strings.TrimSpace(s.SigningKey) == "" {
		return errors.New("JWT signing key required")
	}

	if !s.Algorithm.isValid() {
		return errors.New("unknown signing algorithm")
	}

	return nil
}

func generateJWT(user store.User, settings JWTSettings) (Token, error) {
	if err := settings.validate(); err != nil {
		return Token{}, err
	}

	t := jwt.New(jwt.GetSigningMethod(string(settings.Algorithm)))

	now := time.Now()
	t.Claims = &jwt.MapClaims{
		"iss": settings.Issuer,
		"iat": jwt.NewNumericDate(now),
		"exp": jwt.NewNumericDate(now.Add(time.Second * time.Duration(settings.Lifespan))),
		"sub": user.ID,
	}

	signed, err := t.SignedString([]byte(settings.SigningKey))
	if err != nil {
		return Token{}, err
	}

	return Token{
		AccessToken: signed,
		Lifespan:    settings.Lifespan,
	}, nil
}

func validateJWT(token string, settings JWTSettings) (*jwt.Token, error) {
	return jwt.Parse(
		token,
		func(t *jwt.Token) (interface{}, error) { return []byte(settings.SigningKey), nil },
	)
}
