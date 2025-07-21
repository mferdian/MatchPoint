package service

import (
	"fieldreserve/constants"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	InterfaceJWTService interface {
		GenerateToken(userID string, role string) (string, string, error)
		ValidateToken(token string) (*jwt.Token, error)
		GetUserIDByToken(tokenString string) (string, error)
		GetRoleByToken(tokenString string) (string, error)
	}

	jwtCustomClaims struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
		jwt.RegisteredClaims
	}

	JWTService struct {
		secretKey string
		issuer    string
	}
)

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "Template"
	}

	return secretKey
}

func NewJWTService() *JWTService {
	return &JWTService{
		secretKey: getSecretKey(),
		issuer:    "Template",
	}
}


func (j *JWTService) GenerateToken(userID string, role string) (string, string, error) {
	accessClaims := jwtCustomClaims{
		userID,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 300)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", constants.ErrGenerateAccessToken
	}

	refreshClaims := jwtCustomClaims{
		userID,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 3600 * 24 * 7)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", constants.ErrGenerateRefreshToken
	}

	return accessTokenString, refreshTokenString, nil
}

func (j *JWTService) parseToken(t_ *jwt.Token) (any, error) {
	if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, constants.ErrUnexpectedSigningMethod
	}

	return []byte(j.secretKey), nil
}

func (j *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, j.parseToken)
}

func (j *JWTService) GetUserIDByToken(tokenString string) (string, error) {
	token, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", constants.ErrValidateToken
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return "", constants.ErrTokenInvalid
	}

	return claims.UserID, nil
}

func (j *JWTService) GetRoleByToken(tokenString string) (string, error) {
	token, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", constants.ErrValidateToken
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return "", constants.ErrTokenInvalid
	}

	return claims.Role, nil
}