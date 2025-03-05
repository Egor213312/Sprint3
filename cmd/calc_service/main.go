package main

import (
	"fmt"
	"net/http"

	"github.com/Egor213312/Sprint3/internal/orchestrator"
	"github.com/Egor213312/Sprint3/pkg/config"
	"github.com/Egor213312/Sprint3/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.NewLogger(cfg.LogLevel)

	orchestrator := orchestrator.NewOrchestrator(log, cfg)
	http.HandleFunc("/api/v1/calculate", orchestrator.HandleCalculate)
	http.HandleFunc("/api/v1/expressions", orchestrator.HandleGetExpressions)
	http.HandleFunc("/api/v1/expressions/", orchestrator.HandleGetExpressionByID)

	log.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(fmt.Sprintf("Server failed to start: %s", err))
	}
}
