package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wu-clan/lykn/api"
	"github.com/wu-clan/lykn/config"
	"github.com/wu-clan/lykn/database"
	"github.com/wu-clan/lykn/internal/model"
	"github.com/wu-clan/lykn/internal/service"
)

func main() {
	cfgPath := "config/config.yaml"
	if p := os.Getenv("LYKN_CONFIG"); p != "" {
		cfgPath = p
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get database handle: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	if err := db.AutoMigrate(&model.User{}, &model.Project{}, &model.Feature{}, &model.Plan{}, &model.PlanFeature{}, &model.License{}); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	if err := service.Init(db, cfg); err != nil {
		log.Fatalf("ensure default user: %v", err)
	}

	router := api.RegisterRoutes(cfg)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Printf("http server starting port=%d", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server start failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)
	<-quit

	log.Print("http server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("http server shutdown failed: %v", err)
	}

	log.Print("http server stopped")
}
