package handlers

import (
	"github.com/gofiber/fiber/v2"

	"lostemp/internal/apperror"
	"lostemp/internal/db"
	"lostemp/internal/dto"
	"lostemp/internal/middleware"
	"lostemp/internal/services"
)

type ReservationHandler struct {
	svc *services.ReservationService
}

func NewReservationHandler(s *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{svc: s}
}

// Create cria uma reserva/agendamento (RF06).
// @Summary Criar reserva
// @Description Reserva um equipamento para um período futuro.
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateReservationRequest true "Dados da reserva"
// @Success 201 {object} db.Reservation
// @Failure 409 {object} apperror.AppError
// @Router /reservations [post]
func (h *ReservationHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateReservationRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.BadRequest("corpo da requisição inválido")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	r, err := h.svc.Create(currentUser(c), req)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(r)
}

// List lista reservas paginadas (admin vê a agenda; professor vê as próprias).
// @Summary Listar reservas
// @Tags reservations
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filtro por status (ACTIVE, CANCELLED, FULFILLED)"
// @Param user_id query int false "Filtro por usuário (apenas admin)"
// @Param equipment_id query int false "Filtro por equipamento"
// @Param page query int false "Página (default 1)"
// @Param page_size query int false "Itens por página (default 20, máx 100)"
// @Success 200 {object} map[string]any
// @Router /reservations [get]
func (h *ReservationHandler) List(c *fiber.Ctx) error {
	page, pageSize := pagination(c)
	filter := services.ReservationFilter{
		Status:      c.Query("status"),
		EquipmentID: uint(c.QueryInt("equipment_id")),
	}
	if middleware.Role(c) == db.RoleAdmin {
		filter.UserID = uint(c.QueryInt("user_id"))
	} else {
		filter.ScopeUserID = middleware.UserID(c)
	}
	items, total, err := h.svc.List(filter, page, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(listResponse(items, page, pageSize, total))
}

// Cancel cancela uma reserva (dono ou admin).
// @Summary Cancelar reserva
// @Tags reservations
// @Security BearerAuth
// @Param id path int true "ID da reserva"
// @Success 200
// @Failure 409 {object} apperror.AppError
// @Router /reservations/{id}/cancel [patch]
func (h *ReservationHandler) Cancel(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}
	if err := h.svc.Cancel(currentUser(c), id); err != nil {
		return err
	}
	return c.SendStatus(200)
}
