// Command seed insere APENAS linhas de exemplo (usuários e equipamentos).
// NÃO cria tabelas nem migra schema: o arquivo de banco e as tabelas devem
// já existir (responsabilidade do mantenedor do schema — ver README).
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"lostemp/internal/db"
)

const defaultPassword = "123456"

func main() {
	database, err := db.Open(dbPath())
	if err != nil {
		log.Fatal(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	users := []db.User{
		{Name: "Administrador", Email: "admin@escola.edu", Registration: "ADM-0001", Role: db.RoleAdmin, Status: db.StatusActive},
		{Name: "Professora Maria", Email: "professor@escola.edu", Registration: "PROF-0001", Role: db.RoleProfessor, Status: db.StatusActive},
		{Name: "Técnico João", Email: "ti@escola.edu", Registration: "TEC-0001", Role: db.RoleTechnician, Status: db.StatusActive},
	}

	createdUsers := 0
	for _, u := range users {
		var count int64
		if err := database.Model(&db.User{}).Where("email = ?", u.Email).Count(&count).Error; err != nil {
			log.Fatalf("falha ao consultar usuário %s: %v", u.Email, err)
		}
		if count > 0 {
			fmt.Printf("seed: usuário %s já existe, ignorado\n", u.Email)
			continue
		}
		now := time.Now()
		u.PasswordHash = string(hash)
		u.CreatedAt = now
		u.UpdatedAt = now
		if err := database.Create(&u).Error; err != nil {
			log.Fatalf("falha ao inserir usuário %s: %v", u.Email, err)
		}
		fmt.Printf("seed: usuário criado %s (%s) — senha: %s\n", u.Email, u.Role, defaultPassword)
		createdUsers++
	}

	equipments := []db.Equipment{
		{Identifier: "NT-0001", Brand: "Dell", Model: "Latitude 3420", Status: db.EquipmentAvailable},
		{Identifier: "NT-0002", Brand: "Lenovo", Model: "ThinkPad T14", Status: db.EquipmentAvailable},
		{Identifier: "TB-0001", Brand: "Samsung", Model: "Galaxy Tab S7", Status: db.EquipmentAvailable},
		{Identifier: "TB-0002", Brand: "Apple", Model: "iPad 10", Status: db.EquipmentAvailable},
		{Identifier: "NT-0003", Brand: "Positivo", Model: "Motion Q232A", Status: db.EquipmentAvailable},
	}

	createdEquipments := 0
	for _, e := range equipments {
		var count int64
		if err := database.Model(&db.Equipment{}).Where("identifier = ?", e.Identifier).Count(&count).Error; err != nil {
			log.Fatalf("falha ao consultar equipamento %s: %v", e.Identifier, err)
		}
		if count > 0 {
			fmt.Printf("seed: equipamento %s já existe, ignorado\n", e.Identifier)
			continue
		}
		now := time.Now()
		e.CreatedAt = now
		e.UpdatedAt = now
		if err := database.Create(&e).Error; err != nil {
			log.Fatalf("falha ao inserir equipamento %s: %v", e.Identifier, err)
		}
		fmt.Printf("seed: equipamento criado %s (%s %s)\n", e.Identifier, e.Brand, e.Model)
		createdEquipments++
	}

	fmt.Printf("seed concluído: %d usuários e %d equipamentos inseridos\n", createdUsers, createdEquipments)
}

func dbPath() string {
	// DB_PATH permite apontar para outro arquivo; default ./data/lostemp.db
	if p := os.Getenv("DB_PATH"); p != "" {
		return p
	}
	return "./data/lostemp.db"
}
