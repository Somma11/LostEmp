package dto

import (
	"time"

	"github.com/go-playground/validator/v10"

	"lostemp/internal/apperror"
	"lostemp/internal/db"
)

var validate = validator.New()

// Validate executa a validação de um request DTO.
func Validate(s any) error {
	if err := validate.Struct(s); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return apperror.New(400, "validation_error", err.Error())
		}
		return apperror.BadRequest("corpo da requisição inválido")
	}
	return nil
}

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string    `json:"token"`
	User  *db.User  `json:"user"`
}

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

type CreateUserRequest struct {
	Name         string `json:"name" validate:"required,min=2"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=6"`
	Registration string `json:"registration" validate:"required"`
	Role         string `json:"role" validate:"required,oneof=admin professor technician"`
	Status       string `json:"status" validate:"omitempty,oneof=active inactive"`
}

type UpdateUserRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=2"`
	Email        *string `json:"email" validate:"omitempty,email"`
	Password     *string `json:"password" validate:"omitempty,min=6"`
	Registration *string `json:"registration" validate:"omitempty"`
	Role         *string `json:"role" validate:"omitempty,oneof=admin professor technician"`
	Status       *string `json:"status" validate:"omitempty,oneof=active inactive"`
}

// ---------------------------------------------------------------------------
// Equipments
// ---------------------------------------------------------------------------

type CreateEquipmentRequest struct {
	Identifier string `json:"identifier" validate:"required,min=2"`
	Brand      string `json:"brand" validate:"required,min=1"`
	Model      string `json:"model" validate:"required,min=1"`
	Status     string `json:"status" validate:"omitempty,oneof=AVAILABLE MAINTENANCE"`
}

type UpdateEquipmentRequest struct {
	Identifier *string `json:"identifier" validate:"omitempty,min=2"`
	Brand      *string `json:"brand" validate:"omitempty,min=1"`
	Model      *string `json:"model" validate:"omitempty,min=1"`
	Status     *string `json:"status" validate:"omitempty,oneof=AVAILABLE LOANED MAINTENANCE"`
}

// ---------------------------------------------------------------------------
// Loans
// ---------------------------------------------------------------------------

type CreateLoanRequest struct {
	EquipmentID uint      `json:"equipment_id" validate:"required,gt=0"`
	UserID      uint      `json:"user_id" validate:"required,gt=0"`
	DueAt       time.Time `json:"due_at" validate:"required"`
	Notes       string    `json:"notes"`
}

type EmergencyLoanRequest struct {
	EquipmentID uint      `json:"equipment_id" validate:"required,gt=0"`
	UserID      uint      `json:"user_id" validate:"required,gt=0"`
	DueAt       time.Time `json:"due_at" validate:"required"`
	Notes       string    `json:"notes"`
}

// ---------------------------------------------------------------------------
// Reservations
// ---------------------------------------------------------------------------

type CreateReservationRequest struct {
	EquipmentID uint      `json:"equipment_id" validate:"required,gt=0"`
	UserID      uint      `json:"user_id" validate:"required,gt=0"`
	StartsAt    time.Time `json:"starts_at" validate:"required"`
	EndsAt      time.Time `json:"ends_at" validate:"required"`
}

// ---------------------------------------------------------------------------
// Maintenances
// ---------------------------------------------------------------------------

type CreateMaintenanceRequest struct {
	EquipmentID  uint   `json:"equipment_id" validate:"required,gt=0"`
	Description  string `json:"description" validate:"required,min=3"`
	TechnicianID uint   `json:"technician_id"`
}

type UpdateMaintenanceRequest struct {
	Description *string `json:"description" validate:"omitempty,min=3"`
	Status      *string `json:"status" validate:"omitempty,oneof=OPEN FINISHED"`
}
