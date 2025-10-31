package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtService struct {
	secretKey string
	issure    string
}

func NewJwtService() *JwtService {
	return &JwtService{
		secretKey: "mockSecret",
		issure:    "mockIssure",
	}
}

type Claim struct {
	sub string `json:"sub"`
	jwt.StandardClaims
}

func (s *JwtService) GenerateToken(id uuid.UUID) (string, error) {
	claim := &Claim{
		sub: id.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    s.issure,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	t, err := token.SignedString([]byte(s.secretKey))

	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *JwtService) ValidateToken(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return false
	}

	claims, ok := token.Claims.(*Claim)
	if !ok || !token.Valid {
		return false
	}

	return claims.ExpiresAt > time.Now().Unix()
}
