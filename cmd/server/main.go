package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"vocab-app/internal/config"
	"vocab-app/internal/handler"
	"vocab-app/internal/repository"
	"vocab-app/internal/service"
	"vocab-app/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)
	log.Info("starting vocabulary builder", "port", cfg.Port)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Error("failed to open db", "error", err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		log.Error("failed to connect to db, waiting...", "error", err)
		time.Sleep(5 * time.Second)
		if err := db.Ping(); err != nil {
			log.Error("db ping failed", "error", err)
			os.Exit(1)
		}
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS words (
		id SERIAL PRIMARY KEY,
		word TEXT NOT NULL,
		translation TEXT NOT NULL,
		example TEXT,
		difficulty INTEGER DEFAULT 1,
		next_review TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		status TEXT DEFAULT 'new',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`)
	if err != nil {
		log.Error("failed to init table", "error", err)
		os.Exit(1)
	}
	log.Info("database connected and initialized")

	repo := repository.NewWordRepository(db)
	svc := service.NewWordService(repo, cfg.ReviewIntervalHours)
	h := handler.NewHandler(svc, log)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	mux.Handle("/", http.FileServer(http.Dir("web")))

	server := &http.Server{Addr: ":" + cfg.Port, Handler: mux}

	go func() {
		log.Info("server listening", "addr", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", "error", err)
		os.Exit(1)
	}
	log.Info("server stopped gracefully")
}
