package agent

import (
	"calc_service/pkg/config"
	"calc_service/pkg/logger"
	"time"
)

type Agent struct {
	log *logger.Logger
	cfg *config.Config
}

func NewAgent(log *logger.Logger, cfg *config.Config) *Agent {
	return &Agent{log: log, cfg: cfg}
}

func (a *Agent) Start() {
	for i := 0; i < a.cfg.ComputingPower; i++ {
		go a.worker()
	}
}

func (a *Agent) worker() {
	for {
		task := a.getTask()
		if task != nil {
			result := a.calculate(task)
			a.sendResult(task.ID, result)
		}
		time.Sleep(time.Second)
	}
}

func (a *Agent) getTask() *models.Task {
	// Здесь будет логика получения задачи от оркестратора
	return nil
}

func (a *Agent) calculate(task *models.Task) float64 {
	// Здесь будет логика вычисления задачи
	return 0
}

func (a *Agent) sendResult(taskID int, result float64) {
	// Здесь будет логика отправки результата оркестратору
}
