package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"vocab-app/internal/config"
	"vocab-app/pkg/logger"
)

func main() {
	// Парсим флаги командной строки
	mode := flag.String("mode", "migrate", "Режим: migrate, create-admin, clear-cache")
	email := flag.String("email", "admin@example.com", "Email администратора")
	password := flag.String("password", "admin123", "Пароль администратора")
	flag.Parse()

	cfg := config.Load()
	log := logger.New(cfg.LogLevel) // ← Возвращает *slog.Logger

	// Подключение к БД
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Error("failed to open db", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}

	// Выполняем команду в зависимости от режима
	switch *mode {
	case "migrate":
		if err := runMigrations(db, log); err != nil {
			log.Error("migration failed", "error", err)
			os.Exit(1)
		}
		log.Info("migrations completed successfully")

	case "create-admin":
		if err := createAdminUser(db, log, *email, *password); err != nil {
			log.Error("failed to create admin", "error", err)
			os.Exit(1)
		}
		log.Info("admin user created", "email", *email)

	case "clear-cache":
		log.Info("cache cleared (no-op in this implementation)")

	default:
		log.Error("unknown mode", "mode", *mode)
		os.Exit(1)
	}
}

// runMigrations применяет миграции к базе данных
func runMigrations(db *sql.DB, log *slog.Logger) error { // ← *slog.Logger, не *logger.Logger!
	log.Info("applying migrations...")

	// Миграция 1: таблица users (если ещё не создана)
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT DEFAULT 'user',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	log.Info("migration 1: users table created/verified")

	// Миграция 2: индекс для поиска по email
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)
	`)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	log.Info("migration 2: index on users.email created")

	// Миграция 3: поле last_login в users (идемпотентно)
	_, err = db.Exec(`
		DO $$ 
		BEGIN 
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name='users' AND column_name='last_login'
			) THEN
				ALTER TABLE users ADD COLUMN last_login TIMESTAMP WITH TIME ZONE;
			END IF;
		END $$;
	`)
	if err != nil {
		return fmt.Errorf("failed to add last_login column: %w", err)
	}
	log.Info("migration 3: last_login column added (if not exists)")

	return nil
}

// createAdminUser создаёт пользователя-администратора
func createAdminUser(db *sql.DB, log *slog.Logger, email, password string) error { // ← *slog.Logger
	// В реальном проекте здесь нужно хешировать пароль!
	// Для демонстрации используем простой подход
	_, err := db.Exec(`
		INSERT INTO users (email, password_hash, role) 
		VALUES ($1, $2, 'admin')
		ON CONFLICT (email) DO UPDATE SET role = 'admin'
	`, email, password)

	if err != nil {
		return fmt.Errorf("failed to insert admin: %w", err)
	}
	log.Info("admin user ensured", "email", email)
	return nil
}
