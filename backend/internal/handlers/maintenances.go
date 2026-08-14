package handlers

import (
	"github.com/gofiber/fiber/v2"

	"lostemp/internal/apperror"
	"lostemp/internal/dto"
	"lostemp/internal/services"
)

type MaintenanceHandler struct {
	svc *services.MaintenanceService
}

func NewMaintenanceHandler(s *services.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{svc: s}
}

// Create abre uma manutenção (RF08).
// @Summary Abrir manutenção
// @Description Registra manutenção e coloca o equipamento como MAINTENANCE.
// @Tags maintenances
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateMaintenanceRequest true "Dados da manutenção"
// @Success 201 {object} db.Maintenance
// @Failure 409 {object} apperror.AppError
// @Router /maintenances [post]
func (h *MaintenanceHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateMaintenanceRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.BadRequest("corpo da requisição inválido")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	m, err := h.svc.Create(currentUser(c), req)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(m)
}

// List lista manutenções paginadas.
// @Summary Listar manutenções
// @Tags maintenances
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filtro por status (OPEN, FINISHED)"
// @Param equipment_id query int false "Filtro por equipamento"
// @Param page query int false "Página (default 1)"
// @Param page_size query int false "Itens por página (default 20, máx 100)"
// @Success 200 {object} map[string]any
// @Router /maintenances [get]
func (h *MaintenanceHandler) List(c *fiber.Ctx) error {
	page, pageSize := pagination(c)
	filter := services.MaintenanceFilter{
		Status:      c.Query("status"),
		EquipmentID: uint(c.QueryInt("equipment_id")),
	}
	items, total, err := h.svc.List(filter, page, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(listResponse(items, page, pageSize, total))
}

// Get retorna uma manutenção.
// @Summary Buscar manutenção
// @Tags maintenances
// @Security BearerAuth
// @Param id path int true "ID da manutenção"
// @Success 200 {object} db.Maintenance
// @Router /maintenances/{id} [get]
func (h *MaintenanceHandler) Get(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}
	m, err := h.svc.Get(id)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Update atualiza uma manutenção (ex.: concluir, status FINISHED).
// @Summary Atualizar manutenção
// @Tags maintenances
// @Accept json
// @Security BearerAuth
// @Param id path int true "ID da manutenção"
// @Param body body dto.UpdateMaintenanceRequest true "Campos a atualizar"
// @Success 200 {object} db.Maintenance
// @Router /maintenances/{id} [patch]
func (h *MaintenanceHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}
	var req dto.UpdateMaintenanceRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.BadRequest("corpo da requisição inválido")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	m, err := h.svc.Update(currentUser(c), id, req)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Delete remove uma manutenção.
// @Summary Excluir manutenção
// @Tags maintenances
// @Security BearerAuth
// @Param id path int true "ID da manutenção"
// @Success 204
// @Router /maintenances/{id} [delete]
func (h *MaintenanceHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.Delete(currentUser(c), id); err != nil {
		return err
	}
	return c.SendStatus(204)
}
