package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type Claims struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
	jwtv5.RegisteredClaims
}

func New(secret string, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (m *Manager) Generate(userID int64, username string) (string, string, string, error) {
	access, _, err := m.sign(userID, username, "access", m.accessTTL)
	if err != nil {
		return "", "", "", err
	}
	refresh, tokenID, err := m.sign(userID, username, "refresh", m.refreshTTL)
	return access, refresh, tokenID, err
}

func (m *Manager) Parse(tokenString, expectedType string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwtv5.ParseWithClaims(tokenString, claims, func(token *jwtv5.Token) (any, error) {
		if token.Method.Alg() != jwtv5.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid || claims.TokenType != expectedType {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (m *Manager) sign(userID int64, username, tokenType string, ttl time.Duration) (string, string, error) {
	now := time.Now()
	tokenID := uuid.NewString()
	claims := Claims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ID:        tokenID,
			Issuer:    "go-waste-routes",
			Subject:   username,
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(now.Add(ttl)),
		},
	}
	token, err := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", "", err
	}
	return token, tokenID, nil
}
