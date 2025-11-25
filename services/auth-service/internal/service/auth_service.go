package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shinshark/azure-url-shortener/services/auth-service/internal/config"
	"github.com/shinshark/azure-url-shortener/services/auth-service/internal/models"
	"github.com/shinshark/azure-url-shortener/services/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo   *repository.UserRepository
	Config *config.Config
}

func NewAuthService(repo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{Repo: repo, Config: cfg}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if user exists
	existing, err := s.Repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: string(hashed),
		Role:         "User", // Default role
	}

	if err := s.Repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (string, *models.User, error) {
	user, err := s.Repo.GetUserByUsername(req.Username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	})

	tokenString, err := token.SignedString([]byte(s.Config.JWTSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}
