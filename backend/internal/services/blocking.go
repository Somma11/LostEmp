package services

import (
	"time"

	"gorm.io/gorm"

	"lostemp/internal/db"
)

// HasOverdueLoan é a regra pura de bloqueio por atraso: o usuário está
// bloqueado se possui ao menos um empréstimo ativo cujo prazo (due_at)
// já passou e que ainda não foi devolvido.
func HasOverdueLoan(loans []db.Loan, now time.Time) bool {
	for i := range loans {
		l := &loans[i]
		if l.Status == db.LoanActive && l.ReturnedAt == nil && l.DueAt.Before(now) {
			return true
		}
	}
	return false
}

// UserHasOverdueLoan consulta os empréstimos do usuário e aplica a regra de
// bloqueio por atraso. É usada antes de novos empréstimos e reservas.
func UserHasOverdueLoan(d *gorm.DB, userID uint) (bool, error) {
	var loans []db.Loan
	if err := d.Where("user_id = ? AND status = ?", userID, db.LoanActive).Find(&loans).Error; err != nil {
		return false, err
	}
	return HasOverdueLoan(loans, time.Now()), nil
}
