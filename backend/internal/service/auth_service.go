package service

import (
	"context"
	"time"

	"github.com/wu-clan/lykn/backend/internal/common"
	"github.com/wu-clan/lykn/backend/internal/dao"
	"github.com/wu-clan/lykn/backend/internal/dto"
	"github.com/wu-clan/lykn/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userDAO   *dao.UserDAO
	jwtSecret string
	jwtTTL    time.Duration
}

func NewAuthService(userDAO *dao.UserDAO, jwtSecret string, jwtTTL time.Duration) *AuthService {
	return &AuthService{userDAO: userDAO, jwtSecret: jwtSecret, jwtTTL: jwtTTL}
}

func (s *AuthService) EnsureDefaultUser(ctx context.Context) error {
	count, err := s.userDAO.Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.userDAO.Create(ctx, &model.User{Username: "admin", Password: string(hashed)})
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*dto.LoginResponse, error) {
	user, err := s.userDAO.FindByUsername(ctx, username)
	if err != nil {
		return nil, common.NewUnauthorized("login failed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, common.NewUnauthorized("login failed")
	}
	token, expireAt, err := common.GenerateToken(user.ID, user.Username, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return nil, common.NewInternal(common.CodeLoginFailed, "generate login token failed")
	}
	return &dto.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expireAt,
		User:        dto.UserResponse{ID: user.ID, Username: user.Username},
	}, nil
}

func (s *AuthService) GetCurrentUser(ctx context.Context) (*dto.UserResponse, error) {
	userID, ok := common.UserIDFromContext(ctx)
	if !ok {
		return nil, common.NewUnauthorized("missing user context")
	}
	user, err := s.userDAO.GetByID(ctx, userID)
	if err != nil {
		return nil, common.NewUnauthorized("user not found")
	}
	return &dto.UserResponse{ID: user.ID, Username: user.Username}, nil
}

func Login(ctx context.Context, username, password string) (*dto.LoginResponse, error) {
	if authRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return authRuntime.Login(ctx, username, password)
}

func GetCurrentUser(ctx context.Context) (*dto.UserResponse, error) {
	if authRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return authRuntime.GetCurrentUser(ctx)
}
