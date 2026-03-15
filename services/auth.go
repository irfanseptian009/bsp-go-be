package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/irfanseptian/fims-backend/config"
	"github.com/irfanseptian/fims-backend/database"
	"github.com/irfanseptian/fims-backend/dto"
	"github.com/irfanseptian/fims-backend/middleware"
	"github.com/irfanseptian/fims-backend/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic.
type AuthService struct {
	cfg *config.Config
}

// NewAuthService creates a new AuthService.
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

// Register creates a new user account.
func (s *AuthService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	db := database.GetDB()

	// Check if email already exists
	var existing models.User
	if result := db.Where("email = ?", req.Email).First(&existing); result.Error == nil {
		return nil, errors.New("Email sudah terdaftar")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, errors.New("Gagal memproses password")
	}

	// Create user
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if result := db.Create(&user); result.Error != nil {
		return nil, errors.New("Gagal membuat akun")
	}

	// Generate token
	token, err := s.generateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, errors.New("Gagal membuat token")
	}

	return &dto.AuthResponse{
		User:        user.ToResponse(),
		AccessToken: token,
	}, nil
}

// Login authenticates a user and returns a JWT token.
func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	db := database.GetDB()

	// Find user by email
	var user models.User
	if result := db.Where("email = ?", req.Email).First(&user); result.Error != nil {
		return nil, errors.New("Email atau password salah")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("Email atau password salah")
	}

	// Generate token
	token, err := s.generateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, errors.New("Gagal membuat token")
	}

	return &dto.AuthResponse{
		User:        user.ToResponse(),
		AccessToken: token,
	}, nil
}

// generateToken creates a JWT token with user claims.
func (s *AuthService) generateToken(userID, email, role string) (string, error) {
	claims := middleware.JWTClaims{
		Sub:   userID,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWTExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
