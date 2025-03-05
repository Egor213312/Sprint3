package orchestrator

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Egor213312/Sprint3/internal/models" // Импорт для использования Expression и Task
	"github.com/Egor213312/Sprint3/pkg/config"
	"github.com/Egor213312/Sprint3/pkg/logger"
	"github.com/google/uuid"
)

var (
	expressions = make(map[string]models.Expression) // Используем models.Expression
	tasks       = make(map[string]models.Task)       // Используем models.Task
	mu          sync.Mutex
)

type Orchestrator struct {
	log *logger.Logger
	cfg *config.Config
}

func NewOrchestrator(log *logger.Logger, cfg *config.Config) *Orchestrator {
	return &Orchestrator{log: log, cfg: cfg}
}

func (o *Orchestrator) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	id := uuid.New().String()
	mu.Lock()
	expressions[id] = models.Expression{ // Используем models.Expression
		ID:     id,
		Status: "pending",
		Result: 0,
	}
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (o *Orchestrator) HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var exprs []models.Expression // Используем models.Expression
	for _, expr := range expressions {
		exprs = append(exprs, expr)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": exprs})
}

func (o *Orchestrator) HandleGetExpressionByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/expressions/"):]
	mu.Lock()
	expr, exists := expressions[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
}

func (o *Orchestrator) HandleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		mu.Lock()
		defer mu.Unlock()

		for _, task := range tasks {
			if task.Status == "pending" {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(task)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	} else if r.Method == http.MethodPost {
		var result struct {
			ID     string  `json:"id"`
			Result float64 `json:"result"`
		}
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		if task, exists := tasks[result.ID]; exists {
			task.Status = "completed"
			task.Result = result.Result // Теперь поле Result доступно
			tasks[result.ID] = task
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
