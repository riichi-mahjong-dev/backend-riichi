package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	BaseService
	PlayerService *PlayerService
	AdminService  *AdminService
	JWTSecret     string
}

func NewAuthService(db *gorm.DB, playerService *PlayerService, adminService *AdminService, jwtSecret string) *AuthService {
	return &AuthService{
		BaseService:   BaseService{DB: db},
		PlayerService: playerService,
		AdminService:  adminService,
		JWTSecret:     jwtSecret,
	}
}

func (s *AuthService) LoginPlayer(username, password string) (*models.LoginResponse, error) {
	player, err := s.PlayerService.GetPlayerByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Create token claims
	claims := &models.TokenClaims{
		UserID:   player.ID,
		Username: player.Username,
		UserType: models.UserTypePlayer,
	}

	// Generate tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokens(claims)
	if err != nil {
		return nil, err
	}

	// Create auth user response
	authUser := &models.AuthUser{
		ID:       player.ID,
		Username: player.Username,
		UserType: models.UserTypePlayer,
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
		User:         authUser,
	}, nil
}

func (s *AuthService) LoginAdmin(username, password string) (*models.LoginResponse, error) {
	admin, err := s.AdminService.GetAdminByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Create token claims
	claims := &models.TokenClaims{
		UserID:   admin.ID,
		Username: admin.Username,
		UserType: models.UserTypeAdmin,
		Role:     string(admin.Role),
	}

	// Generate tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokens(claims)
	if err != nil {
		return nil, err
	}

	// Create auth user response
	authUser := &models.AuthUser{
		ID:       admin.ID,
		Username: admin.Username,
		UserType: models.UserTypeAdmin,
		Role:     string(admin.Role),
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
		User:         authUser,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*models.LoginResponse, error) {
	// Parse and validate refresh token
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Generate new tokens
	accessToken, newRefreshToken, expiresAt, err := s.generateTokens(claims)
	if err != nil {
		return nil, err
	}

	// Get user info based on user type
	var authUser *models.AuthUser
	if claims.UserType == models.UserTypePlayer {
		player, err := s.PlayerService.GetPlayerByID(claims.UserID)
		if err != nil {
			return nil, errors.New("user not found")
		}
		authUser = &models.AuthUser{
			ID:       player.ID,
			Username: player.Username,
			UserType: models.UserTypePlayer,
		}
	} else {
		admin, err := s.AdminService.GetAdminByID(claims.UserID)
		if err != nil {
			return nil, errors.New("user not found")
		}
		authUser = &models.AuthUser{
			ID:       admin.ID,
			Username: admin.Username,
			UserType: models.UserTypeAdmin,
			Role:     string(admin.Role),
		}
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    expiresAt,
		User:         authUser,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*models.TokenClaims, error) {
	return s.validateToken(tokenString)
}

func (s *AuthService) generateTokens(claims *models.TokenClaims) (string, string, time.Time, error) {
	now := time.Now()
	accessExpiry := now.Add(24 * time.Hour)      // 24 hours
	refreshExpiry := now.Add(7 * 24 * time.Hour) // 7 days

	// Create access token
	accessClaims := jwt.MapClaims{
		"user_id":   claims.UserID,
		"username":  claims.Username,
		"user_type": claims.UserType,
		"role":      claims.Role,
		"exp":       accessExpiry.Unix(),
		"iat":       now.Unix(),
		"type":      "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Create refresh token
	refreshClaims := jwt.MapClaims{
		"user_id":   claims.UserID,
		"username":  claims.Username,
		"user_type": claims.UserType,
		"role":      claims.Role,
		"exp":       refreshExpiry.Unix(),
		"iat":       now.Unix(),
		"type":      "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessTokenString, refreshTokenString, accessExpiry, nil
}

func (s *AuthService) validateToken(tokenString string) (*models.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("invalid username in token")
	}

	userType, ok := claims["user_type"].(string)
	if !ok {
		return nil, errors.New("invalid user_type in token")
	}

	role, _ := claims["role"].(string) // Optional field

	return &models.TokenClaims{
		UserID:   uint64(userID),
		Username: username,
		UserType: models.UserType(userType),
		Role:     role,
	}, nil
}
