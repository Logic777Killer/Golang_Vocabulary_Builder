package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	repo := repository.NewWordRepository()
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
