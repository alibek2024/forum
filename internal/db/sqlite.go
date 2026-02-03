package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // Используем этот драйвер вместо mattn
)

func InitSQLite(dbPath string, schemaPath string) (*sqlx.DB, error) {
	// Добавляем параметр _pragma=foreign_keys(1) прямо в путь к файлу
	// Для драйвера "sqlite" (modernc.org) это делается так:
	dsn := dbPath + "?_pragma=foreign_keys(1)"

	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// Читаем и выполняем схему
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return nil, err
	}

	return db, nil
}
