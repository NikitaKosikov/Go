package auth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	IdNameURL           = "id"
	UserRole            = "user"
	AdminRole           = "admin"
	PrefixToken         = "Bearer "
)

type Claims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

type TokenManager interface {
	GenerateAccessToken(id string, ttl time.Duration) (string, error)
	VerifyJWTMiddleware(roles ...string) gin.HandlerFunc
	Parse(token string, claims *Claims) (string, error)
	GenerateRefreshToken() (string, error)
	GetTokenFromString(token string, claims *Claims) (*jwt.Token, error)
	ValidateToken(token *jwt.Token, claims *Claims) error
}

type Manager struct {
	secretKey string
}

func NewManager(secretKey string) (*Manager, error) {
	if secretKey == "" {
		return nil, fmt.Errorf("empty secret key")
	}
	return &Manager{secretKey: secretKey}, nil
}

func (m *Manager) GenerateAccessToken(id string, ttl time.Duration) (string, error) {
	claims := &Claims{
		Role: UserRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			Subject:   id,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := accessToken.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", fmt.Errorf("can't signed jwt")
	}

	return PrefixToken + token, nil
}

func (m *Manager) Parse(token string, claims *Claims) (string, error) {
	jwt, _ := m.GetTokenFromString(token, claims)

	if err := m.ValidateToken(jwt, claims); err != nil {
		return "", err
	}
	return jwt.SignedString([]byte(m.secretKey))
}

func (m *Manager) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (m *Manager) GetTokenFromString(token string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})
}

func (m *Manager) ValidateToken(token *jwt.Token, claims *Claims) error {
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return fmt.Errorf("token is not valid")
	}
	return nil
}
