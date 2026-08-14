package services

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"lostemp/internal/apperror"
	"lostemp/internal/db"
	"lostemp/internal/dto"
)

type EquipmentService struct {
	db *gorm.DB
}

func NewEquipmentService(d *gorm.DB) *EquipmentService {
	return &EquipmentService{db: d}
}

type EquipmentFilter struct {
	Identifier string
	Brand      string
	Model      string
	Status     string
}

func (s *EquipmentService) Create(actor *db.User, req dto.CreateEquipmentRequest) (*db.Equipment, error) {
	if req.Status == "" {
		req.Status = db.EquipmentAvailable
	}
	if err := s.ensureIdentifierAvailable(req.Identifier, 0); err != nil {
		return nil, err
	}

	now := time.Now()
	e := db.Equipment{
		Identifier: strings.TrimSpace(req.Identifier),
		Brand:      strings.TrimSpace(req.Brand),
		Model:      strings.TrimSpace(req.Model),
		Status:     req.Status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.db.Create(&e).Error; err != nil {
		return nil, apperror.Internal("falha ao criar equipamento")
	}

	aid := actor.ID
	_ = WriteAudit(s.db, &aid, "equipment.create", "equipment", strconv.FormatUint(uint64(e.ID), 10), map[string]any{
		"identifier": e.Identifier,
	})
	return &e, nil
}

func (s *EquipmentService) List(f EquipmentFilter, page, pageSize int) ([]db.Equipment, int64, error) {
	q := s.db.Model(&db.Equipment{})
	if f.Identifier != "" {
		q = q.Where("identifier LIKE ?", "%"+f.Identifier+"%")
	}
	if f.Brand != "" {
		q = q.Where("brand LIKE ?", "%"+f.Brand+"%")
	}
	if f.Model != "" {
		q = q.Where("model LIKE ?", "%"+f.Model+"%")
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, apperror.Internal("falha ao listar equipamentos")
	}

	var items []db.Equipment
	if err := q.Order("id asc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, apperror.Internal("falha ao listar equipamentos")
	}
	return items, total, nil
}

func (s *EquipmentService) Get(id uint) (*db.Equipment, error) {
	var e db.Equipment
	if err := s.db.First(&e, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("equipamento não encontrado")
		}
		return nil, apperror.Internal("falha ao consultar equipamento")
	}
	return &e, nil
}

func (s *EquipmentService) Update(actor *db.User, id uint, req dto.UpdateEquipmentRequest) (*db.Equipment, error) {
	e, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if req.Identifier != nil && *req.Identifier != e.Identifier {
		if err := s.ensureIdentifierAvailable(*req.Identifier, id); err != nil {
			return nil, err
		}
		e.Identifier = strings.TrimSpace(*req.Identifier)
	}
	if req.Brand != nil {
		e.Brand = strings.TrimSpace(*req.Brand)
	}
	if req.Model != nil {
		e.Model = strings.TrimSpace(*req.Model)
	}
	if req.Status != nil {
		e.Status = *req.Status
	}

	e.UpdatedAt = time.Now()
	if err := s.db.Save(e).Error; err != nil {
		return nil, apperror.Internal("falha ao atualizar equipamento")
	}

	aid := actor.ID
	_ = WriteAudit(s.db, &aid, "equipment.update", "equipment", strconv.FormatUint(uint64(e.ID), 10), map[string]any{
		"identifier": e.Identifier,
		"status":     e.Status,
	})
	return e, nil
}

func (s *EquipmentService) Delete(actor *db.User, id uint) error {
	e, err := s.Get(id)
	if err != nil {
		return err
	}

	var activeLoans int64
	if err := s.db.Model(&db.Loan{}).Where("equipment_id = ? AND status = ?", id, db.LoanActive).Count(&activeLoans).Error; err != nil {
		return apperror.Internal("falha ao verificar empréstimos do equipamento")
	}
	if activeLoans > 0 {
		return apperror.Conflict("equipamento possui empréstimo ativo e não pode ser excluído")
	}

	if err := s.db.Delete(&db.Equipment{}, id).Error; err != nil {
		return apperror.Internal("falha ao excluir equipamento")
	}

	aid := actor.ID
	_ = WriteAudit(s.db, &aid, "equipment.delete", "equipment", strconv.FormatUint(uint64(id), 10), map[string]any{
		"identifier": e.Identifier,
	})
	return nil
}

func (s *EquipmentService) ensureIdentifierAvailable(identifier string, exceptID uint) error {
	var count int64
	q := s.db.Model(&db.Equipment{}).Where("identifier = ?", strings.TrimSpace(identifier))
	if exceptID > 0 {
		q = q.Where("id <> ?", exceptID)
	}
	if err := q.Count(&count).Error; err != nil {
		return apperror.Internal("falha ao verificar identificador")
	}
	if count > 0 {
		return apperror.Conflict("identificador já cadastrado")
	}
	return nil
}
