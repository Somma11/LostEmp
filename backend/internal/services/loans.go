package services

import (
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"

	"lostemp/internal/apperror"
	"lostemp/internal/db"
	"lostemp/internal/dto"
)

type LoanService struct {
	db *gorm.DB
}

func NewLoanService(d *gorm.DB) *LoanService {
	return &LoanService{db: d}
}

type LoanFilter struct {
	Status     string
	UserID     uint
	EquipmentID uint
	Overdue    bool
	ScopeUserID uint
}

// Create registra uma retirada de equipamento (RF02).
func (s *LoanService) Create(actor *db.User, req dto.CreateLoanRequest) (*db.Loan, error) {
	borrower, equipment, err := s.validateLoanBasics(req.EquipmentID, req.UserID, req.DueAt)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	loan := &db.Loan{
		EquipmentID: equipment.ID,
		UserID:      borrower.ID,
		WithdrawnBy: actor.ID,
		WithdrawnAt: now,
		DueAt:       req.DueAt,
		Status:      db.LoanActive,
		Notes:       req.Notes,
		CreatedAt:   now,
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(loan).Error; err != nil {
			return err
		}
		if err := tx.Model(&db.Equipment{}).Where("id = ?", equipment.ID).Update("status", db.EquipmentLoaned).Error; err != nil {
			return err
		}
		// Reservas ativas do mesmo equipamento que coincidem com o período do
		// empréstimo são marcadas como cumpridas.
		if err := tx.Model(&db.Reservation{}).
			Where("equipment_id = ? AND status = ? AND starts_at < ? AND ends_at > ?",
				equipment.ID, db.ReservationActive, loan.DueAt, now).
			Update("status", db.ReservationFulfilled).Error; err != nil {
			return err
		}
		aid := actor.ID
		return WriteAudit(tx, &aid, "loan.create", "loan", strconv.FormatUint(uint64(loan.ID), 10), map[string]any{
			"equipment_id": equipment.ID,
			"user_id":      borrower.ID,
		})
	}); err != nil {
		return nil, apperror.Internal("falha ao registrar empréstimo")
	}
	return loan, nil
}

// Emergency registra uma retirada de emergência (RF07): cancela reservas ativas
// do equipamento e cria o empréstimo, sobrescrevendo a agenda.
func (s *LoanService) Emergency(actor *db.User, req dto.EmergencyLoanRequest) (*db.Loan, error) {
	borrower, equipment, err := s.validateLoanBasics(req.EquipmentID, req.UserID, req.DueAt)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	loan := &db.Loan{
		EquipmentID: equipment.ID,
		UserID:      borrower.ID,
		WithdrawnBy: actor.ID,
		WithdrawnAt: now,
		DueAt:       req.DueAt,
		Status:      db.LoanActive,
		Notes:       req.Notes,
		CreatedAt:   now,
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		var reservations []db.Reservation
		if err := tx.Where("equipment_id = ? AND status = ?", equipment.ID, db.ReservationActive).
			Find(&reservations).Error; err != nil {
			return err
		}
		aid := actor.ID
		for i := range reservations {
			r := &reservations[i]
			if err := tx.Model(r).Update("status", db.ReservationCancelled).Error; err != nil {
				return err
			}
			if err := WriteAudit(tx, &aid, "reservation.cancel_emergency", "reservation",
				strconv.FormatUint(uint64(r.ID), 10), map[string]any{
					"equipment_id": equipment.ID,
					"user_id":      r.UserID,
					"reason":       "retirada de emergência",
				}); err != nil {
				return err
			}
		}

		if err := tx.Create(loan).Error; err != nil {
			return err
		}
		if err := tx.Model(&db.Equipment{}).Where("id = ?", equipment.ID).Update("status", db.EquipmentLoaned).Error; err != nil {
			return err
		}
		return WriteAudit(tx, &aid, "loan.emergency", "loan", strconv.FormatUint(uint64(loan.ID), 10), map[string]any{
			"equipment_id":      equipment.ID,
			"user_id":           borrower.ID,
			"reservations_affected": len(reservations),
		})
	}); err != nil {
		return nil, apperror.Internal("falha ao registrar retirada de emergência")
	}
	return loan, nil
}

// Return encerra um empréstimo (RF03): devolve o equipamento e o torna
// AVAILABLE novamente.
func (s *LoanService) Return(actor *db.User, loanID uint) (*db.Loan, error) {
	var loan db.Loan
	if err := s.db.First(&loan, loanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("empréstimo não encontrado")
		}
		return nil, apperror.Internal("falha ao consultar empréstimo")
	}
	if loan.Status != db.LoanActive {
		return nil, apperror.Conflict("empréstimo não está ativo")
	}

	now := time.Now()
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&loan).Updates(map[string]any{
			"status":      db.LoanReturned,
			"returned_at": now,
			"returned_by": actor.ID,
		}).Error; err != nil {
			return err
		}
		if err := tx.Model(&db.Equipment{}).Where("id = ?", loan.EquipmentID).
			Update("status", db.EquipmentAvailable).Error; err != nil {
			return err
		}
		aid := actor.ID
		return WriteAudit(tx, &aid, "loan.return", "loan", strconv.FormatUint(uint64(loan.ID), 10), map[string]any{
			"equipment_id": loan.EquipmentID,
			"user_id":      loan.UserID,
		})
	}); err != nil {
		return nil, apperror.Internal("falha ao registrar devolução")
	}

	loan.Status = db.LoanReturned
	loan.ReturnedAt = &now
	loan.ReturnedBy = &actor.ID
	return &loan, nil
}

func (s *LoanService) List(f LoanFilter, page, pageSize int) ([]db.Loan, int64, error) {
	q := s.db.Model(&db.Loan{})

	if f.ScopeUserID > 0 {
		q = q.Where("user_id = ?", f.ScopeUserID)
	}
	if f.UserID > 0 {
		q = q.Where("user_id = ?", f.UserID)
	}
	if f.EquipmentID > 0 {
		q = q.Where("equipment_id = ?", f.EquipmentID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Overdue {
		q = q.Where("status = ? AND due_at < ? AND returned_at IS NULL", db.LoanActive, time.Now())
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, apperror.Internal("falha ao listar empréstimos")
	}

	var loans []db.Loan
	order := "withdrawn_at desc"
	if f.Overdue {
		order = "due_at asc"
	}
	if err := q.Order(order).Offset((page - 1) * pageSize).Limit(pageSize).Find(&loans).Error; err != nil {
		return nil, 0, apperror.Internal("falha ao listar empréstimos")
	}
	return loans, total, nil
}

// validateLoanBasics aplica as regras comuns antes de qualquer retirada:
// equipamento existe e está AVAILABLE; usuário é professor ativo; professor
// sem empréstimo vencido; prazo de devolução no futuro.
func (s *LoanService) validateLoanBasics(equipmentID, userID uint, dueAt time.Time) (*db.User, *db.Equipment, error) {
	var equipment db.Equipment
	if err := s.db.First(&equipment, equipmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, apperror.NotFound("equipamento não encontrado")
		}
		return nil, nil, apperror.Internal("falha ao consultar equipamento")
	}
	if equipment.Status != db.EquipmentAvailable {
		return nil, nil, apperror.Conflict("equipamento não está disponível para retirada")
	}

	var borrower db.User
	if err := s.db.First(&borrower, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, apperror.NotFound("usuário não encontrado")
		}
		return nil, nil, apperror.Internal("falha ao consultar usuário")
	}
	if borrower.Role != db.RoleProfessor {
		return nil, nil, apperror.BadRequest("empréstimo só pode ser feito para um professor")
	}
	if borrower.Status != db.StatusActive {
		return nil, nil, apperror.Conflict("professor está inativo")
	}

	overdue, err := UserHasOverdueLoan(s.db, borrower.ID)
	if err != nil {
		return nil, nil, apperror.Internal("falha ao verificar pendências")
	}
	if overdue {
		return nil, nil, apperror.Conflict("professor possui empréstimo vencido e está bloqueado para novos empréstimos")
	}

	if !dueAt.After(time.Now()) {
		return nil, nil, apperror.BadRequest("o prazo de devolução deve ser no futuro")
	}

	return &borrower, &equipment, nil
}
