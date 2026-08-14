package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"lostemp/internal/config"
	"lostemp/internal/db"
	"lostemp/internal/middleware"
	"lostemp/internal/services"
)

// RegisterRoutes monta todas as rotas com prefixo /api/v1.
func RegisterRoutes(app *fiber.App, cfg *config.Config, database *gorm.DB) {
	authSvc := services.NewAuthService(database, cfg)
	userSvc := services.NewUserService(database)
	equipmentSvc := services.NewEquipmentService(database)
	loanSvc := services.NewLoanService(database)
	reservationSvc := services.NewReservationService(database)
	maintenanceSvc := services.NewMaintenanceService(database)

	authH := NewAuthHandler(authSvc)
	userH := NewUserHandler(userSvc)
	equipmentH := NewEquipmentHandler(equipmentSvc)
	loanH := NewLoanHandler(loanSvc)
	reservationH := NewReservationHandler(reservationSvc)
	maintenanceH := NewMaintenanceHandler(maintenanceSvc)
	healthH := NewHealthHandler(database)

	api := app.Group("/api/v1")

	api.Get("/health", healthH.Health)
	api.Post("/auth/login", authH.Login)

	auth := api.Group("", middleware.JWT(cfg.JWTSecret))
	auth.Get("/auth/me", authH.Me)

	// users — admin
	auth.Post("/users", middleware.RequireRoles(db.RoleAdmin), userH.Create)
	auth.Get("/users", middleware.RequireRoles(db.RoleAdmin), userH.List)
	auth.Get("/users/:id", middleware.RequireRoles(db.RoleAdmin), userH.Get)
	auth.Patch("/users/:id", middleware.RequireRoles(db.RoleAdmin), userH.Update)
	auth.Delete("/users/:id", middleware.RequireRoles(db.RoleAdmin), userH.Delete)

	// equipments — consulta para todos autenticados; mutação apenas admin
	auth.Get("/equipments", equipmentH.List)
	auth.Get("/equipments/:id", equipmentH.Get)
	auth.Post("/equipments", middleware.RequireRoles(db.RoleAdmin), equipmentH.Create)
	auth.Patch("/equipments/:id", middleware.RequireRoles(db.RoleAdmin), equipmentH.Update)
	auth.Delete("/equipments/:id", middleware.RequireRoles(db.RoleAdmin), equipmentH.Delete)

	// loans — admin e professor criam; devolução e emergência apenas admin
	auth.Post("/loans", middleware.RequireRoles(db.RoleAdmin, db.RoleProfessor), loanH.Create)
	auth.Post("/loans/emergency", middleware.RequireRoles(db.RoleAdmin), loanH.Emergency)
	auth.Post("/loans/:id/return", middleware.RequireRoles(db.RoleAdmin), loanH.Return)
	auth.Get("/loans", loanH.List)
	auth.Get("/loans/overdue", loanH.ListOverdue)

	// reservations — professor reserva; admin/professor cancelam; lista agenda
	auth.Post("/reservations", middleware.RequireRoles(db.RoleAdmin, db.RoleProfessor), reservationH.Create)
	auth.Get("/reservations", reservationH.List)
	auth.Patch("/reservations/:id/cancel", reservationH.Cancel)

	// maintenances — técnico e admin
	auth.Post("/maintenances", middleware.RequireRoles(db.RoleAdmin, db.RoleTechnician), maintenanceH.Create)
	auth.Get("/maintenances", middleware.RequireRoles(db.RoleAdmin, db.RoleTechnician), maintenanceH.List)
	auth.Get("/maintenances/:id", middleware.RequireRoles(db.RoleAdmin, db.RoleTechnician), maintenanceH.Get)
	auth.Patch("/maintenances/:id", middleware.RequireRoles(db.RoleAdmin, db.RoleTechnician), maintenanceH.Update)
	auth.Delete("/maintenances/:id", middleware.RequireRoles(db.RoleAdmin, db.RoleTechnician), maintenanceH.Delete)
}
