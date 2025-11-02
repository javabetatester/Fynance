package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtService struct {
	secretKey string
	issuer    string
}

func NewJwtService() *JwtService {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "fynance_secure_jwt_secret_key_2024"
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "fynance_api"
	}

	return &JwtService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

type Claim struct {
	Sub string `json:"sub"`
	jwt.StandardClaims
}

func (s *JwtService) GenerateToken(id uuid.UUID) (string, error) {
	claim := &Claim{
		Sub: id.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
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
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("método de assinatura inválido", jwt.ValidationErrorSignatureInvalid)
		}
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

func (s *JwtService) ParseToken(tokenString string) (*Claim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("método de assinatura inválido", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claim)
	if !ok || !token.Valid {
		return nil, jwt.NewValidationError("token inválido", jwt.ValidationErrorClaimsInvalid)
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, jwt.NewValidationError("token expirado", jwt.ValidationErrorExpired)
	}

	return claims, nil
}
