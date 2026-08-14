package handlers

import (
	"github.com/gofiber/fiber/v2"

	"lostemp/internal/apperror"
	"lostemp/internal/dto"
	"lostemp/internal/services"
)

type UserHandler struct {
	svc *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{svc: s}
}

// Create cria um usuário (apenas admin).
// @Summary Criar usuário
// @Description Cadastra um usuário (admin, professor ou técnico).
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateUserRequest true "Dados do usuário"
// @Success 201 {object} db.User
// @Failure 400 {object} apperror.AppError
// @Failure 409 {object} apperror.AppError
// @Router /users [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.BadRequest("corpo da requisição inválido")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	u, err := h.svc.Create(currentUser(c), req)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(u)
}

// List lista usuários paginados (apenas admin).
// @Summary Listar usuários
// @Description Lista usuários com filtros opcionais e paginação.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param role query string false "Filtro por papel"
// @Param status query string false "Filtro por status"
// @Param page query int false "Página (default 1)"
// @Param page_size query int false "Itens por página (default 20, máx 100)"
// @Success 200 {object} map[string]any
// @Router /users [get]
func (h *UserHandler) List(c *fiber.Ctx) error {
	page, pageSize := pagination(c)
	filter := struct {
		Role   string
		Status string
	}{
		Role:   c.Query("role"),
		Status: c.Query("status"),
	}
	users, total, err := h.svc.List(filter, page, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(listResponse(users, page, pageSize, total))
}

// Get retorna um usuário (apenas admin).
// @Summary Buscar usuário
// @Tags users
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Success 200 {object} db.User
// @Router /users/{id} [get]
func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}
	u, err := h.svc.Get(id)
	if err != nil {
		return err
	}
	return c.JSON(u)
}

// Update atualiza um usuário (apenas admin).
// @Summary Atualizar usuário
// @Tags users
// @Accept json
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Param body body dto.UpdateUserRequest true "Campos a atualizar"
// @Success 200 {object} db.User
// @Router /users/{id} [patch]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}
	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.BadRequest("corpo da requisição inválido")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	u, err := h.svc.Update(currentUser(c), id, req)
	if err != nil {
		return err
	}
	return c.JSON(u)
}

// Delete remove um usuário (apenas admin).
// @Summary Excluir usuário
// @Tags users
// @Security BearerAuth
// @Param id path int true "ID do usuário"
// @Success 204
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.Delete(currentUser(c), id); err != nil {
		return err
	}
	return c.SendStatus(204)
}
