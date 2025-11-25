package service

import (
	"fmt"

	"github.com/OzyKleyton/agendamento-api/internal/model"
	"github.com/OzyKleyton/agendamento-api/internal/model/user"
	"github.com/OzyKleyton/agendamento-api/internal/repository"
	jwt "github.com/OzyKleyton/agendamento-api/utils/auth"
	"github.com/OzyKleyton/agendamento-api/utils/security"
)

type AuthService interface {
	Login(loginReq *user.LoginRequest) *model.Response
	RefreshToken(refreshReq *user.RefreshRequest) *model.Response
	ValidateToken(token string) (*jwt.Claims, error)
}

type AuthServiceImpl struct {
	userRepo repository.UserRepository
	jwtUtil  *jwt.JWT
}

func NewAuthService(userRepo repository.UserRepository, jwtUtil *jwt.JWT) AuthService {
	return &AuthServiceImpl{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

func (as *AuthServiceImpl) Login(loginReq *user.LoginRequest) *model.Response {

	if loginReq.Email == "" || loginReq.Password == "" {
		return model.NewErrorResponse(fmt.Errorf("email e senha são obrigatórios"), 400)
	}

	userFound, err := as.userRepo.FindByEmail(loginReq.Email)
	if err != nil {
		return model.NewErrorResponse(fmt.Errorf("credenciais inválidas"), 401)
	}

	err = security.ComparePassword(userFound.Password, loginReq.Password)
	if err != nil {
		return model.NewErrorResponse(fmt.Errorf("credenciais inválidas"), 401)
	}

	tokens, err := as.jwtUtil.GenerateTokensWithID(userFound.ID, userFound.Username, userFound.Email)
	if err != nil {
		return model.NewErrorResponse(err, 500)
	}

	tokenResponse := &user.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
		User:         userFound.ToUserRes(),
	}

	return model.NewSuccessResponse(tokenResponse)
}

func (as *AuthServiceImpl) RefreshToken(refreshReq *user.RefreshRequest) *model.Response {
	if refreshReq.RefreshToken == "" {
		return model.NewErrorResponse(fmt.Errorf("refresh token é obrigatório"), 400)
	}

	tokens, err := as.jwtUtil.RefreshTokens(refreshReq.RefreshToken)
	if err != nil {
		return model.NewErrorResponse(err, 401)
	}

	claims, err := as.jwtUtil.ValidateToken(tokens.AccessToken)
	if err != nil {
		return model.NewErrorResponse(err, 401)
	}

	userFound, err := as.userRepo.FindByEmail(claims.Email)
	if err != nil {
		return model.NewErrorResponse(err, 404)
	}

	tokenResponse := &user.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
		User:         userFound.ToUserRes(),
	}

	return model.NewSuccessResponse(tokenResponse)
}

func (as *AuthServiceImpl) ValidateToken(token string) (*jwt.Claims, error) {
	return as.jwtUtil.ValidateToken(token)
}
