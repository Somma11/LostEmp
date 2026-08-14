package handlers

import (
	"github.com/gofiber/fiber/v2"

	"lostemp/internal/apperror"
	"lostemp/internal/dto"
	"lostemp/internal/services"
)

type EquipmentHandler struct {
	svc *services.EquipmentService
}

func NewEquipmentHandler(s *services.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{svc: s}
}

// Create cadastra um equipamento (apenas admin).
// @Summary Cadastrar equipamento
// @Description Registra equipamento com identificador único, marca e modelo (RF01).
// @Tags equipments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateEquipmentRequest true "Dados do equipamento"
// @Success 201 {object} db.Equipment
// @Failure 409 {object} apperror.AppError
// @Router /equipments [post]
func (h *EquipmentHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateEquipmentRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.BadRequest("corpo da requisição inválido")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	e, err := h.svc.Create(currentUser(c), req)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(e)
}

// List lista equipamentos com filtros e paginação.
// @Summary Listar equipamentos
// @Description Lista o catálogo de equipamentos com filtros por identificador, marca, modelo e status.
// @Tags equipments
// @Produce json
// @Security BearerAuth
// @Param identifier query string false "Filtro por identificador (parcial)"
// @Param brand query string false "Filtro por marca (parcial)"
// @Param model query string false "Filtro por modelo (parcial)"
// @Param status query string false "Filtro por status (AVAILABLE, LOANED, MAINTENANCE)"
// @Param page query int false "Página (default 1)"
// @Param page_size query int false "Itens por página (default 20, máx 100)"
// @Success 200 {object} map[string]any
// @Router /equipments [get]
func (h *EquipmentHandler) List(c *fiber.Ctx) error {
	page, pageSize := pagination(c)
	filter := services.EquipmentFilter{
		Identifier: c.Query("identifier"),
		Brand:      c.Query("brand"),
		Model:      c.Query("model"),
		Status:     c.Query("status"),
	}
	items, total, err := h.svc.List(filter, page, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(listResponse(items, page, pageSize, total))
}

// Get retorna um equipamento.
// @Summary Buscar equipamento
// @Tags equipments
// @Security BearerAuth
// @Param id path int true "ID do equipamento"
// @Success 200 {object} db.Equipment
// @Router /equipments/{id} [get]
func (h *EquipmentHandler) Get(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}
	e, err := h.svc.Get(id)
	if err != nil {
		return err
	}
	return c.JSON(e)
}

// Update atualiza um equipamento (apenas admin).
// @Summary Atualizar equipamento
// @Tags equipments
// @Accept json
// @Security BearerAuth
// @Param id path int true "ID do equipamento"
// @Param body body dto.UpdateEquipmentRequest true "Campos a atualizar"
// @Success 200 {object} db.Equipment
// @Router /equipments/{id} [patch]
func (h *EquipmentHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}
	var req dto.UpdateEquipmentRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.BadRequest("corpo da requisição inválido")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	e, err := h.svc.Update(currentUser(c), id, req)
	if err != nil {
		return err
	}
	return c.JSON(e)
}

// Delete remove um equipamento (apenas admin).
// @Summary Excluir equipamento
// @Tags equipments
// @Security BearerAuth
// @Param id path int true "ID do equipamento"
// @Success 204
// @Router /equipments/{id} [delete]
func (h *EquipmentHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.Delete(currentUser(c), id); err != nil {
		return err
	}
	return c.SendStatus(204)
}
