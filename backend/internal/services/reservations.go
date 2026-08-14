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

type ReservationService struct {
	db *gorm.DB
}

func NewReservationService(d *gorm.DB) *ReservationService {
	return &ReservationService{db: d}
}

type ReservationFilter struct {
	Status      string
	UserID      uint
	EquipmentID uint
	ScopeUserID uint
}

// Create registra uma reserva de equipamento (RF06).
func (s *ReservationService) Create(actor *db.User, req dto.CreateReservationRequest) (*db.Reservation, error) {
	var equipment db.Equipment
	if err := s.db.First(&equipment, req.EquipmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("equipamento não encontrado")
		}
		return nil, apperror.Internal("falha ao consultar equipamento")
	}

	var borrower db.User
	if err := s.db.First(&borrower, req.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("usuário não encontrado")
		}
		return nil, apperror.Internal("falha ao consultar usuário")
	}
	if borrower.Role != db.RoleProfessor {
		return nil, apperror.BadRequest("reserva só pode ser feita para um professor")
	}
	if borrower.Status != db.StatusActive {
		return nil, apperror.Conflict("professor está inativo")
	}

	overdue, err := UserHasOverdueLoan(s.db, borrower.ID)
	if err != nil {
		return nil, apperror.Internal("falha ao verificar pendências")
	}
	if overdue {
		return nil, apperror.Conflict("professor possui empréstimo vencido e está bloqueado para novas reservas")
	}

	if !req.EndsAt.After(req.StartsAt) {
		return nil, apperror.BadRequest("o término da reserva deve ser após o início")
	}
	if !req.StartsAt.After(time.Now()) {
		return nil, apperror.BadRequest("a reserva deve iniciar no futuro")
	}

	var conflicts int64
	if err := s.db.Model(&db.Reservation{}).
		Where("equipment_id = ? AND status = ? AND starts_at < ? AND ends_at > ?",
			req.EquipmentID, db.ReservationActive, req.EndsAt, req.StartsAt).
		Count(&conflicts).Error; err != nil {
		return nil, apperror.Internal("falha ao verificar conflitos de agenda")
	}
	if conflicts > 0 {
		return nil, apperror.Conflict("já existe reserva ativa conflitante para este equipamento no período")
	}

	now := time.Now()
	r := db.Reservation{
		EquipmentID: req.EquipmentID,
		UserID:      req.UserID,
		StartsAt:    req.StartsAt,
		EndsAt:      req.EndsAt,
		Status:      db.ReservationActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.db.Create(&r).Error; err != nil {
		return nil, apperror.Internal("falha ao criar reserva")
	}

	aid := actor.ID
	_ = WriteAudit(s.db, &aid, "reservation.create", "reservation", strconv.FormatUint(uint64(r.ID), 10), map[string]any{
		"equipment_id": req.EquipmentID,
		"user_id":      req.UserID,
	})
	return &r, nil
}

func (s *ReservationService) List(f ReservationFilter, page, pageSize int) ([]db.Reservation, int64, error) {
	q := s.db.Model(&db.Reservation{})

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

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, apperror.Internal("falha ao listar reservas")
	}

	var items []db.Reservation
	if err := q.Order("starts_at asc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, apperror.Internal("falha ao listar reservas")
	}
	return items, total, nil
}

func (s *ReservationService) Cancel(actor *db.User, id uint) error {
	var r db.Reservation
	if err := s.db.First(&r, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("reserva não encontrada")
		}
		return apperror.Internal("falha ao consultar reserva")
	}
	if actor.Role != db.RoleAdmin && actor.ID != r.UserID {
		return apperror.Forbidden("só o dono da reserva ou um administrador pode cancelá-la")
	}
	if r.Status != db.ReservationActive {
		return apperror.Conflict("reserva não está ativa")
	}

	if err := s.db.Model(&r).Update("status", db.ReservationCancelled).Error; err != nil {
		return apperror.Internal("falha ao cancelar reserva")
	}

	aid := actor.ID
	_ = WriteAudit(s.db, &aid, "reservation.cancel", "reservation", strconv.FormatUint(uint64(r.ID), 10), map[string]any{
		"equipment_id": r.EquipmentID,
		"user_id":      r.UserID,
	})
	return nil
}
