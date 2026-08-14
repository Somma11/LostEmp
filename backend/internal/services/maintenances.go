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

type MaintenanceService struct {
	db *gorm.DB
}

func NewMaintenanceService(d *gorm.DB) *MaintenanceService {
	return &MaintenanceService{db: d}
}

type MaintenanceFilter struct {
	Status      string
	EquipmentID uint
}

// Create abre uma manutenção (RF08): o equipamento passa a MAINTENANCE e só
// pode ser retirado novamente após concluída.
func (s *MaintenanceService) Create(actor *db.User, req dto.CreateMaintenanceRequest) (*db.Maintenance, error) {
	var equipment db.Equipment
	if err := s.db.First(&equipment, req.EquipmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("equipamento não encontrado")
		}
		return nil, apperror.Internal("falha ao consultar equipamento")
	}
	if equipment.Status == db.EquipmentLoaned {
		return nil, apperror.Conflict("equipamento está emprestado e não pode entrar em manutenção")
	}

	technicianID := req.TechnicianID
	if technicianID == 0 {
		if actor.Role != db.RoleTechnician {
			return nil, apperror.BadRequest("informe o technician_id")
		}
		technicianID = actor.ID
	}

	var technician db.User
	if err := s.db.First(&technician, technicianID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("técnico não encontrado")
		}
		return nil, apperror.Internal("falha ao consultar técnico")
	}
	if technician.Role != db.RoleTechnician {
		return nil, apperror.BadRequest("technician_id deve ser um usuário com papel de técnico")
	}

	now := time.Now()
	m := db.Maintenance{
		EquipmentID:  req.EquipmentID,
		TechnicianID: technicianID,
		Description:  req.Description,
		Status:       db.MaintenanceOpen,
		StartedAt:    now,
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&m).Error; err != nil {
			return err
		}
		if err := tx.Model(&db.Equipment{}).Where("id = ?", equipment.ID).
			Update("status", db.EquipmentMaintenance).Error; err != nil {
			return err
		}
		aid := actor.ID
		return WriteAudit(tx, &aid, "maintenance.create", "maintenance", strconv.FormatUint(uint64(m.ID), 10), map[string]any{
			"equipment_id": equipment.ID,
		})
	}); err != nil {
		return nil, apperror.Internal("falha ao abrir manutenção")
	}
	return &m, nil
}

func (s *MaintenanceService) List(f MaintenanceFilter, page, pageSize int) ([]db.Maintenance, int64, error) {
	q := s.db.Model(&db.Maintenance{})
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.EquipmentID > 0 {
		q = q.Where("equipment_id = ?", f.EquipmentID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, apperror.Internal("falha ao listar manutenções")
	}

	var items []db.Maintenance
	if err := q.Order("started_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, apperror.Internal("falha ao listar manutenções")
	}
	return items, total, nil
}

func (s *MaintenanceService) Get(id uint) (*db.Maintenance, error) {
	var m db.Maintenance
	if err := s.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("manutenção não encontrada")
		}
		return nil, apperror.Internal("falha ao consultar manutenção")
	}
	return &m, nil
}

func (s *MaintenanceService) Update(actor *db.User, id uint, req dto.UpdateMaintenanceRequest) (*db.Maintenance, error) {
	m, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if req.Description != nil {
			m.Description = *req.Description
		}
		if req.Status != nil && *req.Status != m.Status {
			switch *req.Status {
			case db.MaintenanceFinished:
				now := time.Now()
				m.FinishedAt = &now
				m.Status = db.MaintenanceFinished
				if err := tx.Model(&db.Equipment{}).Where("id = ?", m.EquipmentID).
					Update("status", db.EquipmentAvailable).Error; err != nil {
					return err
				}
			case db.MaintenanceOpen:
				m.FinishedAt = nil
				m.Status = db.MaintenanceOpen
				if err := tx.Model(&db.Equipment{}).Where("id = ?", m.EquipmentID).
					Update("status", db.EquipmentMaintenance).Error; err != nil {
					return err
				}
			}
		}
		if err := tx.Save(m).Error; err != nil {
			return err
		}
		aid := actor.ID
		return WriteAudit(tx, &aid, "maintenance.update", "maintenance", strconv.FormatUint(uint64(m.ID), 10), map[string]any{
			"status":  m.Status,
			"description": m.Description,
		})
	}); err != nil {
		return nil, apperror.Internal("falha ao atualizar manutenção")
	}
	return m, nil
}

func (s *MaintenanceService) Delete(actor *db.User, id uint) error {
	m, err := s.Get(id)
	if err != nil {
		return err
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&db.Maintenance{}, id).Error; err != nil {
			return err
		}
		if m.Status == db.MaintenanceOpen {
			if err := tx.Model(&db.Equipment{}).Where("id = ?", m.EquipmentID).
				Update("status", db.EquipmentAvailable).Error; err != nil {
				return err
			}
		}
		aid := actor.ID
		return WriteAudit(tx, &aid, "maintenance.delete", "maintenance", strconv.FormatUint(uint64(id), 10), nil)
	}); err != nil {
		return apperror.Internal("falha ao excluir manutenção")
	}
	return nil
}
