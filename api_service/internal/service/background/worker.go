package background

import (
	"ai-assistant-api/internal/model"
	"ai-assistant-api/internal/utils/logger/sl"
	"context"
	"log/slog"
	"time"
)

func (s *Service) startWorker(ctx context.Context, done <-chan struct{}) {
	const op = "service.background.startWorker"
	log := s.log.With(
		slog.String("op", op),
	)

	sendTicker := time.NewTicker(s.cfg.SendWorkerTimeout)
	go func() {
		for {
			select {
			case <-done:
				log.Info("shutdown task send worker")
				return
			case <-sendTicker.C:
				tasks, err := s.taskService.TaskListByStatus(ctx, model.StatusSending)
				if err != nil {
					log.Error("failed to get task list", sl.Err(err))
					continue
				}

				for _, task := range tasks {
					err = s.taskService.SendTask(ctx, task)
					if err != nil {
						log.Error("failde to send task", sl.Err(err))
						continue
					}
				}
			}
		}
	}()
}
