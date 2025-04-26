package auth

import (
	"errors"

	"github.com/syned13/flight-prices-api/internal/middleware"
	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/internal/repository/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user models.User) error
	Login(user models.User) (string, error)
}

type authService struct {
	authRepository auth.AuthRepository
}

func NewAuthService(authRepository auth.AuthRepository) AuthService {
	return &authService{
		authRepository: authRepository,
	}
}

func (s *authService) Register(user models.User) error {
	existingUser, err := s.authRepository.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.authRepository.CreateUser(user)
}

func (s *authService) Login(user models.User) (string, error) {
	storedUser, err := s.authRepository.GetUserByUsername(user.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := middleware.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}
