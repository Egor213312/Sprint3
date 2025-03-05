package main

import (
	"calc_service/internal/orchestrator"
	"calc_service/pkg/config"
	"calc_service/pkg/logger"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.NewLogger(cfg.LogLevel)

	orchestrator := orchestrator.NewOrchestrator(log, cfg)
	http.HandleFunc("/api/v1/calculate", orchestrator.HandleCalculate)
	http.HandleFunc("/api/v1/expressions", orchestrator.HandleGetExpressions)
	http.HandleFunc("/api/v1/expressions/", orchestrator.HandleGetExpressionByID)
	http.HandleFunc("/internal/task", orchestrator.HandleTask)

	log.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
