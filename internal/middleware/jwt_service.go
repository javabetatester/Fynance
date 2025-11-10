package middleware

import (
	"context"
	"errors"

	"Fynance/config"
	"Fynance/internal/domain/user"

	"github.com/golang-jwt/jwt"
	"github.com/oklog/ulid/v2"
)

type JwtService struct {
	secretKey   string
	issuer      string
	tokenTTL    int64
	parser      *jwt.Parser
	userService *user.Service
}

func NewJwtService(settings config.JWTConfig, userService *user.Service) (*JwtService, error) {
	if settings.ExpiresIn <= 0 {
		return nil, errors.New("JWT_EXPIRES_IN inválido")
	}
	parser := &jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}
	return &JwtService{
		secretKey:   settings.SecretKey,
		issuer:      settings.Issuer,
		tokenTTL:    int64(settings.ExpiresIn.Seconds()),
		parser:      parser,
		userService: userService,
	}, nil
}

type Claim struct {
	Sub  string    `json:"sub"`
	Plan user.Plan `json:"plan"`
	jwt.StandardClaims
}

func (s *JwtService) GenerateToken(ctx context.Context, id ulid.ULID) (string, error) {
	plan, err := s.userService.GetPlan(ctx, id)
	if err != nil {
		return "", err
	}
	now := jwt.TimeFunc().Unix()
	claim := &Claim{
		Sub:  id.String(),
		Plan: plan,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now + s.tokenTTL,
			IssuedAt:  now,
			Issuer:    s.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JwtService) ValidateToken(tokenString string) bool {
	claims, err := s.parse(tokenString)
	if err != nil {
		return false
	}
	now := jwt.TimeFunc().Unix()
	if !claims.VerifyIssuer(s.issuer, true) {
		return false
	}
	if !claims.VerifyExpiresAt(now, true) {
		return false
	}
	return true
}

func (s *JwtService) ParseToken(tokenString string) (*Claim, error) {
	claims, err := s.parse(tokenString)
	if err != nil {
		return nil, err
	}
	now := jwt.TimeFunc().Unix()
	if !claims.VerifyIssuer(s.issuer, true) {
		return nil, jwt.NewValidationError("emissor inválido", jwt.ValidationErrorIssuer)
	}
	if !claims.VerifyExpiresAt(now, true) {
		return nil, jwt.NewValidationError("token expirado", jwt.ValidationErrorExpired)
	}
	return claims, nil
}

func (s *JwtService) parse(tokenString string) (*Claim, error) {
	token, err := s.parser.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claim)
	if !ok || !token.Valid {
		return nil, jwt.NewValidationError("token inválido", jwt.ValidationErrorClaimsInvalid)
	}
	return claims, nil
}
