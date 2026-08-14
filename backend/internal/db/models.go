package db

import "time"

// Papéis de usuário (RBAC).
const (
	RoleAdmin      = "admin"
	RoleProfessor  = "professor"
	RoleTechnician = "technician"
)

// Status de usuário.
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

// Status de equipamento.
const (
	EquipmentAvailable   = "AVAILABLE"
	EquipmentLoaned      = "LOANED"
	EquipmentMaintenance = "MAINTENANCE"
)

// Status de empréstimo.
const (
	LoanActive   = "ACTIVE"
	LoanReturned = "RETURNED"
)

// Status de reserva.
const (
	ReservationActive   = "ACTIVE"
	ReservationCancelled = "CANCELLED"
	ReservationFulfilled = "FULFILLED"
)

// Status de manutenção.
const (
	MaintenanceOpen     = "OPEN"
	MaintenanceFinished = "FINISHED"
)

// User mapeia a tabela `users`.
type User struct {
	ID           uint      `gorm:"column:id;primaryKey" json:"id"`
	Name         string    `gorm:"column:name" json:"name"`
	Email        string    `gorm:"column:email" json:"email"`
	PasswordHash string    `gorm:"column:password_hash" json:"-"`
	Registration string    `gorm:"column:registration" json:"registration"`
	Role         string    `gorm:"column:role" json:"role"`
	Status       string    `gorm:"column:status" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string { return "users" }

// Equipment mapeia a tabela `equipments`.
type Equipment struct {
	ID         uint      `gorm:"column:id;primaryKey" json:"id"`
	Identifier string    `gorm:"column:identifier" json:"identifier"`
	Brand      string    `gorm:"column:brand" json:"brand"`
	Model      string    `gorm:"column:model" json:"model"`
	Status     string    `gorm:"column:status" json:"status"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Equipment) TableName() string { return "equipments" }

// Loan mapeia a tabela `loans`.
type Loan struct {
	ID          uint       `gorm:"column:id;primaryKey" json:"id"`
	EquipmentID uint       `gorm:"column:equipment_id" json:"equipment_id"`
	UserID      uint       `gorm:"column:user_id" json:"user_id"`
	WithdrawnBy uint       `gorm:"column:withdrawn_by" json:"withdrawn_by"`
	WithdrawnAt time.Time  `gorm:"column:withdrawn_at" json:"withdrawn_at"`
	DueAt       time.Time  `gorm:"column:due_at" json:"due_at"`
	ReturnedAt  *time.Time `gorm:"column:returned_at" json:"returned_at"`
	ReturnedBy  *uint      `gorm:"column:returned_by" json:"returned_by"`
	Status      string     `gorm:"column:status" json:"status"`
	Notes       string     `gorm:"column:notes" json:"notes"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
}

func (Loan) TableName() string { return "loans" }

// Reservation mapeia a tabela `reservations`.
type Reservation struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	EquipmentID uint    `gorm:"column:equipment_id" json:"equipment_id"`
	UserID    uint      `gorm:"column:user_id" json:"user_id"`
	StartsAt  time.Time `gorm:"column:starts_at" json:"starts_at"`
	EndsAt    time.Time `gorm:"column:ends_at" json:"ends_at"`
	Status    string    `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Reservation) TableName() string { return "reservations" }

// Maintenance mapeia a tabela `maintenances`.
type Maintenance struct {
	ID           uint       `gorm:"column:id;primaryKey" json:"id"`
	EquipmentID  uint       `gorm:"column:equipment_id" json:"equipment_id"`
	TechnicianID uint       `gorm:"column:technician_id" json:"technician_id"`
	Description  string     `gorm:"column:description" json:"description"`
	Status       string     `gorm:"column:status" json:"status"`
	StartedAt    time.Time  `gorm:"column:started_at" json:"started_at"`
	FinishedAt   *time.Time `gorm:"column:finished_at" json:"finished_at"`
}

func (Maintenance) TableName() string { return "maintenances" }

// AuditLog mapeia a tabela `audit_logs`.
type AuditLog struct {
	ID         uint      `gorm:"column:id;primaryKey" json:"id"`
	UserID     *uint     `gorm:"column:user_id" json:"user_id"`
	Action     string    `gorm:"column:action" json:"action"`
	EntityType string    `gorm:"column:entity_type" json:"entity_type"`
	EntityID   string    `gorm:"column:entity_id" json:"entity_id"`
	Payload    string    `gorm:"column:payload" json:"payload"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_logs" }
