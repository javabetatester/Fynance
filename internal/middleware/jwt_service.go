package middleware

import (
	"Fynance/internal/domain/user"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/oklog/ulid/v2"
)

type JwtService struct {
	secretKey   string
	issuer      string
	userService *user.Service
}

func NewJwtService(userService *user.Service) *JwtService {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "fynance_secure_jwt_secret_key_2024"
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "fynance_api"
	}

	return &JwtService{
		secretKey:   secretKey,
		issuer:      issuer,
		userService: userService,
	}
}

type Claim struct {
	Sub  string    `json:"sub"`
	Plan user.Plan `json:"plan"`
	jwt.StandardClaims
}

func (s *JwtService) GenerateToken(id ulid.ULID) (string, error) {
	plan, err := s.userService.GetPlan(id)
	if err != nil {
		return "", err
	}

	claim := &Claim{
		Sub:  id.String(),
		Plan: plan,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(s.secretKey))
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
