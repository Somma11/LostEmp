package services

import (
	"errors"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"lostemp/internal/apperror"
	"lostemp/internal/db"
	"lostemp/internal/dto"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(d *gorm.DB) *UserService {
	return &UserService{db: d}
}

func (s *UserService) Create(actor *db.User, req dto.CreateUserRequest) (*db.User, error) {
	if req.Status == "" {
		req.Status = db.StatusActive
	}
	if err := s.ensureEmailAvailable(req.Email, 0); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.Internal("falha ao gerar hash de senha")
	}

	now := time.Now()
	u := db.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Registration: req.Registration,
		Role:         req.Role,
		Status:       req.Status,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.db.Create(&u).Error; err != nil {
		return nil, apperror.Internal("falha ao criar usuário")
	}

	aid := actor.ID
	_ = WriteAudit(s.db, &aid, "user.create", "user", strconv.FormatUint(uint64(u.ID), 10), map[string]any{
		"email": u.Email,
		"role":  u.Role,
	})
	return &u, nil
}

func (s *UserService) List(filter struct {
	Role   string
	Status string
}, page, pageSize int) ([]db.User, int64, error) {
	q := s.db.Model(&db.User{})
	if filter.Role != "" {
		q = q.Where("role = ?", filter.Role)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, apperror.Internal("falha ao listar usuários")
	}

	var users []db.User
	if err := q.Order("id asc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, apperror.Internal("falha ao listar usuários")
	}
	return users, total, nil
}

func (s *UserService) Get(id uint) (*db.User, error) {
	var u db.User
	if err := s.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("usuário não encontrado")
		}
		return nil, apperror.Internal("falha ao consultar usuário")
	}
	return &u, nil
}

func (s *UserService) Update(actor *db.User, id uint, req dto.UpdateUserRequest) (*db.User, error) {
	u, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	if req.Email != nil && *req.Email != u.Email {
		if err := s.ensureEmailAvailable(*req.Email, id); err != nil {
			return nil, err
		}
		u.Email = *req.Email
	}
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Registration != nil {
		u.Registration = *req.Registration
	}
	if req.Role != nil {
		u.Role = *req.Role
	}
	if req.Status != nil {
		u.Status = *req.Status
	}
	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, apperror.Internal("falha ao gerar hash de senha")
		}
		u.PasswordHash = string(hash)
	}

	u.UpdatedAt = time.Now()
	if err := s.db.Save(u).Error; err != nil {
		return nil, apperror.Internal("falha ao atualizar usuário")
	}

	aid := actor.ID
	_ = WriteAudit(s.db, &aid, "user.update", "user", strconv.FormatUint(uint64(u.ID), 10), map[string]any{
		"email": u.Email,
		"role":  u.Role,
	})
	return u, nil
}

func (s *UserService) Delete(actor *db.User, id uint) error {
	u, err := s.Get(id)
	if err != nil {
		return err
	}

	if actor.ID == id {
		return apperror.Conflict("não é possível excluir o próprio usuário")
	}

	var activeLoans int64
	if err := s.db.Model(&db.Loan{}).Where("user_id = ? AND status = ?", id, db.LoanActive).Count(&activeLoans).Error; err != nil {
		return apperror.Internal("falha ao verificar empréstimos do usuário")
	}
	if activeLoans > 0 {
		return apperror.Conflict("usuário possui empréstimo ativo e não pode ser excluído")
	}

	if err := s.db.Delete(&db.User{}, id).Error; err != nil {
		return apperror.Internal("falha ao excluir usuário")
	}

	aid := actor.ID
	_ = WriteAudit(s.db, &aid, "user.delete", "user", strconv.FormatUint(uint64(id), 10), map[string]any{
		"email": u.Email,
	})
	return nil
}

func (s *UserService) ensureEmailAvailable(email string, exceptID uint) error {
	var count int64
	q := s.db.Model(&db.User{}).Where("email = ?", email)
	if exceptID > 0 {
		q = q.Where("id <> ?", exceptID)
	}
	if err := q.Count(&count).Error; err != nil {
		return apperror.Internal("falha ao verificar e-mail")
	}
	if count > 0 {
		return apperror.Conflict("e-mail já cadastrado")
	}
	return nil
}

