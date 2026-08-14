package services

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"lostemp/internal/db"
)

// WriteAudit registra um AuditLog. Toda ação crítica do sistema passa por aqui
// (auditoria: usuário, data/hora e entidade afetada).
func WriteAudit(d *gorm.DB, userID *uint, action, entityType, entityID string, payload any) error {
	raw := ""
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			raw = ""
		} else {
			raw = string(b)
		}
	}
	return d.Create(&db.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Payload:    raw,
		CreatedAt:  time.Now(),
	}).Error
}
