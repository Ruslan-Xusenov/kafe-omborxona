package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"kafe-omborxona/internal/domain"
)

type AuthService struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewAuthService(repo domain.UserRepository, secret string) *AuthService {
	return &AuthService{userRepo: repo, jwtSecret: secret}
}

func (s *AuthService) Login(req domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("noto'g'ri login yoki parol")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("noto'g'ri login yoki parol")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{Token: token, User: *user}, nil
}

func (s *AuthService) generateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

func (s *AuthService) CreateUser(req domain.CreateUserRequest) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Role:         req.Role,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) GetAllUsers() ([]domain.User, error) {
	return s.userRepo.GetAll()
}

func (s *AuthService) GetUserByID(id int) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *AuthService) UpdateUser(id int, req domain.UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	user.Username = req.Username
	user.FullName = req.FullName
	user.Role = req.Role
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hash)
	} else {
		user.PasswordHash = ""
	}
	return s.userRepo.Update(user)
}

func (s *AuthService) DeleteUser(id int) error {
	return s.userRepo.Delete(id)
}
