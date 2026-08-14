package services

import (
	"testing"
	"time"
)

func TestOverlaps(t *testing.T) {
	base := time.Date(2026, 8, 13, 12, 0, 0, 0, time.UTC)

	day := func(h int) time.Time { return base.Add(time.Duration(h) * time.Hour) }

	tests := []struct {
		name        string
		aStart, aEnd time.Time
		bStart, bEnd time.Time
		want        bool
	}{
		{"intervalos iguais", day(8), day(10), day(8), day(10), true},
		{"B dentro de A", day(8), day(12), day(9), day(11), true},
		{"A dentro de B", day(9), day(11), day(8), day(12), true},
		{"B termina no início de A (contíguo)", day(8), day(10), day(6), day(8), false},
		{"B começa no fim de A (contíguo)", day(8), day(10), day(10), day(12), false},
		{"disjuntos", day(8), day(10), day(11), day(12), false},
		{"sobreposição parcial", day(8), day(10), day(9), day(12), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Overlaps(tt.aStart, tt.aEnd, tt.bStart, tt.bEnd); got != tt.want {
				t.Errorf("Overlaps() = %v, want %v", got, tt.want)
			}
		})
	}
}
