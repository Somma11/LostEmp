package handlers

import (
	"github.com/gofiber/fiber/v2"

	"lostemp/internal/apperror"
	"lostemp/internal/dto"
	"lostemp/internal/middleware"
	"lostemp/internal/services"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}

// Login autentica um usuário e retorna um JWT.
// @Summary Login
// @Description Autentica com e-mail e senha, retornando um JWT.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Credenciais"
// @Success 200 {object} services.LoginResult
// @Failure 400 {object} apperror.AppError
// @Failure 401 {object} apperror.AppError
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.BadRequest("corpo da requisição inválido")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	res, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

// Me retorna os dados do usuário autenticado.
// @Summary Dados do usuário logado
// @Description Retorna o usuário autenticado pelo token.
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} db.User
// @Failure 401 {object} apperror.AppError
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	u, err := h.svc.Me(middleware.UserID(c))
	if err != nil {
		return err
	}
	return c.JSON(u)
}
