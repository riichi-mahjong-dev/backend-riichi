package utils

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	// "github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthToken struct {
	secretKey        string
	refreshSecretKey string
}

type DataToken struct {
	userId  uint
	expired int64
}

func InitializeAuth(secretKey, refreshSecretKey string) *AuthToken {
	return &AuthToken{
		secretKey:        secretKey,
		refreshSecretKey: refreshSecretKey,
	}
}

func (authToken *AuthToken) GenerateToken(user *models.User) (string, string, error) {
	userClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	accessTokenString, err := accessToken.SignedString([]byte(authToken.secretKey))

	if err != nil {
		return "", "", err
	}

	refreshUserClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshUserClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(authToken.refreshSecretKey))

	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (authToken *AuthToken) ValidateToken(tokenString string) (float64, float64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}

		return []byte(authToken.secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && err == nil && token.Valid {
		return claims["exp"].(float64), claims["user_id"].(float64), nil
	}

	return 0, 0, err
}

func (authToken *AuthToken) ValidateRefresh(tokenString string) (float64, float64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}

		return []byte(authToken.refreshSecretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && err == nil && token.Valid {
		return claims["exp"].(float64), claims["user_id"].(float64), nil
	}

	return 0, 0, err
}

func (authToken *AuthToken) GeneratePassword(length int) (string, string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var password []byte
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", "", err
		}
		password = append(password, charset[index.Int64()])
	}

	hashPassword, err := HashPasssword(string(password))

	if err != nil {
		return "", "", err
	}

	return string(password), hashPassword, nil
}

func (authToken *AuthToken) ValidatePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func HashPasssword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
