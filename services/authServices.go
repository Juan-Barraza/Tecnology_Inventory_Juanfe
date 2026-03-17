package services

import (
	"errors"
	"time"

	dtos "inventory-juanfe/dtos/request"
	response "inventory-juanfe/dtos/response"
	repository "inventory-juanfe/repositories"
	"inventory-juanfe/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(req dtos.LoginRequest) (*response.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("internal error")
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("could not generate token")
	}

	go func() {
		_ = s.userRepo.UpdateLastLogin(user.ID, time.Now())
	}()

	return &response.LoginResponse{
		Token: token,
		User: response.UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
