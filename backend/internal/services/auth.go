package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"lostemp/internal/apperror"
	"lostemp/internal/config"
	"lostemp/internal/db"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(d *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: d, cfg: cfg}
}

type LoginResult struct {
	Token string  `json:"token"`
	User  *db.User `json:"user"`
}

// Login autentica por e-mail/senha (bcrypt) e emite um JWT.
func (s *AuthService) Login(email, password string) (*LoginResult, error) {
	var u db.User
	err := s.db.Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.Unauthorized("credenciais inválidas")
	}
	if err != nil {
		return nil, apperror.Internal("falha ao consultar usuário")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, apperror.Unauthorized("credenciais inválidas")
	}

	if u.Status != db.StatusActive {
		return nil, apperror.Forbidden("usuário inativo")
	}

	token, err := s.generateToken(&u)
	if err != nil {
		return nil, apperror.Internal("falha ao gerar token")
	}

	uid := u.ID
	_ = WriteAudit(s.db, &uid, "auth.login", "user", strconv.FormatUint(uint64(u.ID), 10), map[string]any{
		"email": u.Email,
	})

	return &LoginResult{Token: token, User: &u}, nil
}

func (s *AuthService) Me(userID uint) (*db.User, error) {
	var u db.User
	if err := s.db.First(&u, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("usuário não encontrado")
		}
		return nil, apperror.Internal("falha ao consultar usuário")
	}
	return &u, nil
}

func (s *AuthService) generateToken(u *db.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   strconv.FormatUint(uint64(u.ID), 10),
		"role":  u.Role,
		"email": u.Email,
		"name":  u.Name,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(s.cfg.JWTTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
