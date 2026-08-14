package handlers

import (
	"github.com/gofiber/fiber/v2"

	"lostemp/internal/apperror"
	"lostemp/internal/db"
	"lostemp/internal/dto"
	"lostemp/internal/middleware"
	"lostemp/internal/services"
)

type LoanHandler struct {
	svc *services.LoanService
}

func NewLoanHandler(s *services.LoanService) *LoanHandler {
	return &LoanHandler{svc: s}
}

// Create registra uma retirada de equipamento (RF02).
// @Summary Registrar retirada
// @Description Cria um empréstimo vinculando equipamento a um professor.
// @Tags loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateLoanRequest true "Dados da retirada"
// @Success 201 {object} db.Loan
// @Failure 409 {object} apperror.AppError
// @Router /loans [post]
func (h *LoanHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateLoanRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.BadRequest("corpo da requisição inválido")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	loan, err := h.svc.Create(currentUser(c), req)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(loan)
}

// Emergency registra uma retirada de emergência, sobrescrevendo a agenda (RF07).
// @Summary Retirada de emergência
// @Description Cancela reservas ativas/futuras do equipamento e cria o empréstimo.
// @Tags loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.EmergencyLoanRequest true "Dados da retirada de emergência"
// @Success 201 {object} db.Loan
// @Failure 409 {object} apperror.AppError
// @Router /loans/emergency [post]
func (h *LoanHandler) Emergency(c *fiber.Ctx) error {
	var req dto.EmergencyLoanRequest
	if err := c.BodyParser(&req); err != nil {
		return apperror.BadRequest("corpo da requisição inválido")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	loan, err := h.svc.Emergency(currentUser(c), req)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(loan)
}

// Return registra a devolução de um equipamento (RF03).
// @Summary Registrar devolução
// @Description Encerra o empréstimo e devolve o equipamento para AVAILABLE.
// @Tags loans
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do empréstimo"
// @Success 200 {object} db.Loan
// @Failure 409 {object} apperror.AppError
// @Router /loans/{id}/return [post]
func (h *LoanHandler) Return(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}
	loan, err := h.svc.Return(currentUser(c), id)
	if err != nil {
		return err
	}
	return c.JSON(loan)
}

// List lista empréstimos paginados (admin vê todos; professor vê os próprios).
// @Summary Listar empréstimos
// @Tags loans
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filtro por status (ACTIVE, RETURNED)"
// @Param user_id query int false "Filtro por usuário (apenas admin)"
// @Param equipment_id query int false "Filtro por equipamento"
// @Param page query int false "Página (default 1)"
// @Param page_size query int false "Itens por página (default 20, máx 100)"
// @Success 200 {object} map[string]any
// @Router /loans [get]
func (h *LoanHandler) List(c *fiber.Ctx) error {
	page, pageSize := pagination(c)
	filter := h.scopeFilter(c)
	filter.Status = c.Query("status")
	filter.EquipmentID = uint(c.QueryInt("equipment_id"))
	if middleware.Role(c) == db.RoleAdmin {
		filter.UserID = uint(c.QueryInt("user_id"))
	}
	loans, total, err := h.svc.List(filter, page, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(listResponse(loans, page, pageSize, total))
}

// ListOverdue lista empréstimos em atraso (RF04).
// @Summary Listar empréstimos em atraso
// @Description Lista empréstimos ativos com prazo vencido.
// @Tags loans
// @Produce json
// @Security BearerAuth
// @Param page query int false "Página (default 1)"
// @Param page_size query int false "Itens por página (default 20, máx 100)"
// @Success 200 {object} map[string]any
// @Router /loans/overdue [get]
func (h *LoanHandler) ListOverdue(c *fiber.Ctx) error {
	page, pageSize := pagination(c)
	filter := h.scopeFilter(c)
	filter.Overdue = true
	loans, total, err := h.svc.List(filter, page, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(listResponse(loans, page, pageSize, total))
}

func (h *LoanHandler) scopeFilter(c *fiber.Ctx) services.LoanFilter {
	if middleware.Role(c) == db.RoleAdmin {
		return services.LoanFilter{}
	}
	return services.LoanFilter{ScopeUserID: middleware.UserID(c)}
}
