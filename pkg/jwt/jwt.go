package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	Secret []byte
	Issuer string
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewManager(secret, issuer string) *Manager {
	return &Manager{
		Secret: []byte(secret),
		Issuer: issuer,
	}
}

func (m *Manager) GenerateAccessToken(
	userID uint,
	role string,
	ttl time.Duration,
) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Issuer,
			Subject:   "access_token",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(m.Secret)
}

func (m *Manager) ParseAccessToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return m.Secret, nil
		},
		jwt.WithIssuer(m.Issuer),
		jwt.WithSubject("access_token"),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
