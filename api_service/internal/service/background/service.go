package background

import (
	"ai-assistant-api/internal/config"
	"ai-assistant-api/internal/model"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type TaskService interface {
	TaskListByStatus(ctx context.Context, status model.TaskStatus) ([]model.Task, error)
	SendTask(ctx context.Context, task model.Task) error
	FinishTask(ctx context.Context, taskID uuid.UUID, response string) error
}

type Service struct {
	done        chan struct{}
	cfg         *config.Config
	taskService TaskService
	log         *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger, taskService TaskService) *Service {
	return &Service{
		taskService: taskService,
		log:         log,
		cfg:         cfg,
	}
}

func (s *Service) Start() {
	s.done = make(chan struct{})
	ctx := context.Background()

	s.startWorker(ctx, s.done)
	s.startConsumer(ctx, s.done)
}

func (s *Service) Shutdown() {
	close(s.done)
}
