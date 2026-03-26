package main

import (
	"errors"
	"log"
	"net/http"

	"packplanner/internal/application/packapp"
	"packplanner/internal/config"
	"packplanner/internal/domain/pack"
	"packplanner/internal/infrastructure/repository/memory"
	"packplanner/internal/transport/httpapi"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// Wire the application from the outside in so each layer stays focused on a single concern.
	repo, err := memory.NewPackSizeRepository(cfg.DefaultPackSizes)
	if err != nil {
		log.Fatalf("create repository: %v", err)
	}

	planner := pack.NewOptimalPlanner()
	service := packapp.NewService(repo, planner)
	server := httpapi.NewServer(service, cfg.AllowedOrigins)

	if err := server.Start(":" + cfg.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("start server: %v", err)
	}
}
