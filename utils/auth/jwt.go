package jwt

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWT struct {
	config Config
}

type Config struct {
	SecretKey     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func NewJWT() *JWT {
	return &JWT{
		config: Config{
			SecretKey:     viper.GetString("JWT_SECRET"),
			AccessExpiry:  time.Hour * time.Duration(getEnvInt("JWT_ACCESS_EXPIRY_HOURS", 1)),
			RefreshExpiry: time.Hour * 24 * time.Duration(getEnvInt("JWT_REFRESH_EXPIRY_DAYS", 7)),
		},
	}
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (j *JWT) GenerateTokensWithID(userID uint, username, email string) (*TokenPair, error) {
	if j.config.SecretKey == "" {
		return nil, errors.New("secret key não pode ser vazia")
	}

	accessExpiresAt := time.Now().Add(j.config.AccessExpiry)
	accessClaims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshExpiresAt := time.Now().Add(j.config.RefreshExpiry)
	refreshClaims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}

func (j *JWT) GenerateTokens(username, email string) (*TokenPair, error) {
	return j.GenerateTokensWithID(0, username, email)
}

func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	if j.config.SecretKey == "" {
		return nil, errors.New("secret key não pode ser vazia")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}

func (j *JWT) RefreshTokens(refreshToken string) (*TokenPair, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("refresh token inválido")
	}

	return j.GenerateTokensWithID(claims.UserID, claims.Username, claims.Email)
}
