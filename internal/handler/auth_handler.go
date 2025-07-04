package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend/configs"
	"github.com/riichi-mahjong-dev/backend/internal/models"
	"github.com/riichi-mahjong-dev/backend/internal/services"
	"github.com/riichi-mahjong-dev/backend/utils"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	userService   *services.UserService
	authGenerator *utils.AuthToken
	oauth2Config  *configs.ConfigOauth
}

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *models.User `json:"user"`
}

func InitializeAuthHandler(userService *services.UserService, authGenerator *utils.AuthToken, oauth2Config *configs.ConfigOauth) *AuthHandler {
	return &AuthHandler{
		userService:   userService,
		authGenerator: authGenerator,
		oauth2Config:  oauth2Config,
	}
}

func (authHandler *AuthHandler) LoginUser(c *fiber.Ctx) error {
	loginInput := new(LoginInput)

	if err := c.BodyParser(loginInput); err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	user, err := authHandler.userService.UserLogin(loginInput.Username, loginInput.Password)

	if err != nil {
		return utils.ResponseError(fiber.StatusUnauthorized, "Unauthorized", err)(c)
	}

	accessToken, refreshToken, err := authHandler.authGenerator.GenerateToken(user)

	if err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusOK, "Success authorization", TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})(c)
}

func (authHandler *AuthHandler) LoginAdmin(c *fiber.Ctx) error {
	loginInput := new(LoginInput)

	if err := c.BodyParser(loginInput); err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	user, err := authHandler.userService.AdminLogin(loginInput.Username, loginInput.Password)

	if err != nil {
		return utils.ResponseError(fiber.StatusUnauthorized, "Unauthorized", err)(c)
	}

	accessToken, refreshToken, err := authHandler.authGenerator.GenerateToken(user)

	if err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusOK, "Success authorization", TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})(c)
}

func (authHandler *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	refreshTokenInput := new(RefreshTokenInput)

	if err := c.BodyParser(refreshTokenInput); err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	exp, userId, err := authHandler.authGenerator.ValidateRefresh(refreshTokenInput.RefreshToken)

	if err != nil {
		return utils.ResponseError(fiber.StatusUnauthorized, "Unauthorized", err)(c)
	}

	now := time.Now().Unix()

	if int64(exp) < now {
		return utils.ResponseError(fiber.StatusUnauthorized, "Unauthorized", fmt.Errorf("refresh token expired"))(c)
	}

	user, err := authHandler.userService.GetUser(uint(userId))

	if err != nil {
		return utils.ResponseError(fiber.StatusUnauthorized, "Unauthorized", err)(c)
	}

	accessToken, refreshToken, err := authHandler.authGenerator.GenerateToken(user)

	if err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusOK, "Success authorization", TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})(c)
}

func (authHandler *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	key := authHandler.oauth2Config.Key
	url := authHandler.oauth2Config.Config.AuthCodeURL(key, oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

func (authHandler *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	token, err := authHandler.oauth2Config.Config.Exchange(context.Background(), code)

	if err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	client := authHandler.oauth2Config.Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	name := userInfo["name"].(string)
	email := userInfo["email"].(string)
	phoneNumber := userInfo["phone"].(string)

	user := authHandler.userService.RegisterFromGoogle(name, email, phoneNumber)

	accessToken, refreshToken, err := authHandler.authGenerator.GenerateToken(user)

	if err != nil {
		return utils.ResponseError(fiber.StatusInternalServerError, "Something went wrong", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusOK, "Success authorization", TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})(c)
}

func (authHandler *AuthHandler) RegisterUser(c *fiber.Ctx) error {
	userInput := new(models.User)

	if err := c.BodyParser(userInput); err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	if err := authHandler.userService.CreateUser(userInput); err != nil {
		return utils.ResponseError(fiber.StatusBadRequest, "Bad request", err)(c)
	}

	return utils.ResponseSuccess(fiber.StatusCreated, "Success register user", nil)(c)
}
