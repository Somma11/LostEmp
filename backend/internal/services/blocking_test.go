package services

import (
	"testing"
	"time"

	"lostemp/internal/db"
)

func TestHasOverdueLoan(t *testing.T) {
	now := time.Date(2026, 8, 13, 12, 0, 0, 0, time.UTC)

	base := func(due time.Time, status string, returned bool) db.Loan {
		l := db.Loan{
			ID:         1,
			Status:     status,
			DueAt:      due,
			ReturnedAt: nil,
		}
		if returned {
			r := time.Now()
			l.ReturnedAt = &r
		}
		return l
	}

	tests := []struct {
		name  string
		loans []db.Loan
		want  bool
	}{
		{
			name: "sem empréstimos não bloqueia",
			want: false,
		},
		{
			name:  "empréstimo ativo dentro do prazo não bloqueia",
			loans: []db.Loan{base(now.Add(2*24*time.Hour), db.LoanActive, false)},
			want:  false,
		},
		{
			name:  "empréstimo ativo vencido bloqueia",
			loans: []db.Loan{base(now.Add(-24*time.Hour), db.LoanActive, false)},
			want:  true,
		},
		{
			name:  "empréstimo vencido mas devolvido não bloqueia",
			loans: []db.Loan{base(now.Add(-24*time.Hour), db.LoanReturned, true)},
			want:  false,
		},
		{
			name: "devolvido em dia + ativo vencido bloqueia",
			loans: []db.Loan{
				base(now.Add(-24*time.Hour), db.LoanReturned, true),
				base(now.Add(-3*24*time.Hour), db.LoanActive, false),
			},
			want: true,
		},
		{
			name:  "vencendo exatamente agora não bloqueia",
			loans: []db.Loan{base(now, db.LoanActive, false)},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasOverdueLoan(tt.loans, now); got != tt.want {
				t.Errorf("HasOverdueLoan() = %v, want %v", got, tt.want)
			}
		})
	}
}
