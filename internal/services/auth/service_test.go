package auth

import (
	"errors"
	"testing"

	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/mocks/repository/auth"
	"github.com/zeebo/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Register(t *testing.T) {
	t.Run("should register a user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		authRepository := auth.NewMockAuthRepository(ctrl)
		authRepository.EXPECT().GetUserByUsername(gomock.Any()).Return(nil, nil)
		authRepository.EXPECT().CreateUser(gomock.Any()).Return(nil)

		authService := NewAuthService(authRepository)

		user := models.User{
			Username: "testuser",
			Password: "testpassword",
		}

		err := authService.Register(user)
		assert.NoError(t, err)
	})

	t.Run("should return an error if the user already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		authRepository := auth.NewMockAuthRepository(ctrl)
		authRepository.EXPECT().GetUserByUsername(gomock.Any()).Return(&models.User{
			Username: "testuser",
			Password: "testpassword",
		}, nil)

		authService := NewAuthService(authRepository)

		user := models.User{
			Username: "testuser",
			Password: "testpassword",
		}

		err := authService.Register(user)
		assert.Error(t, err)
	})
}

func TestAuthService_Login(t *testing.T) {

	t.Run("should return an error if the user does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		authRepository := auth.NewMockAuthRepository(ctrl)
		authRepository.EXPECT().GetUserByUsername(gomock.Any()).Return(nil, errors.New("user not found"))

		authService := NewAuthService(authRepository)

		user := models.User{
			Username: "testuser",
			Password: "testpassword",
		}

		_, err := authService.Login(user)
		assert.Error(t, err)
	})
}
