package db

import (
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open abre o banco SQLite existente em path.
//
// IMPORTANTE: este projeto NÃO cria nem migra schema. O arquivo e as tabelas
// são de responsabilidade do mantenedor do schema (ver README, seção
// "Pendências de schema"). Aqui apenas abrimos a conexão.
func Open(path string) (*gorm.DB, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("arquivo de banco não encontrado em %q — o arquivo %s já deve existir com o schema correto: %w", path, path, err)
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir banco em %q: %w", path, err)
	}
	return db, nil
}
