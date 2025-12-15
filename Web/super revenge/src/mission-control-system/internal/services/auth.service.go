package services

import (
	"app-api/internal/dtos"
	"app-api/internal/repositories"
	"app-api/internal/utils"
	"context"

	"go.uber.org/zap"
)

type AuthService struct {
	logger *zap.Logger
	repo   *repositories.AuthRepository
}

func NewAuthService(logger *zap.Logger) *AuthService {
	return &AuthService{
		logger: logger,
		repo:   repositories.NewAuthRepository(logger),
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, registerUserDto dtos.RegisterUserDto) (dtos.StructuredResponse, error) {
	return s.repo.RegisterUser(ctx, registerUserDto)
}

func (s *AuthService) LoginUser(ctx context.Context, loginUserDto dtos.LoginUserDto) (dtos.StructuredResponse, error) {
	return s.repo.LoginUser(ctx, loginUserDto)
}

func (s *AuthService) CheckToken(ctx context.Context, tokenDto dtos.TokenDto) (dtos.StructuredResponse, error) {
	token, err := utils.ValidateToken(tokenDto.Token)
	if err != nil {
		s.logger.Error("Error validating token", zap.Error(err))
		return dtos.StructuredResponse{
			Success: false,
			Status:  500,
			Message: "Internal server error",
			Payload: nil,
		}, err
	}

	if token == nil {
		return dtos.StructuredResponse{
			Success: false,
			Status:  401,
			Message: "Invalid token",
			Payload: nil,
		}, nil
	}

	return dtos.StructuredResponse{
		Success: true,
		Status:  200,
		Message: "Token is valid",
		Payload: token,
	}, nil

}
