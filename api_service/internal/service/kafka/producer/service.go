package sender

import (
	"ai-assistant-api/internal/config"
	"ai-assistant-api/internal/model"
	"ai-assistant-api/internal/utils/logger/sl"
	"context"
	"encoding/json"
	"log/slog"
)

type Service struct {
	cfg    *config.Config
	client *client
	log    *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) (*Service, error) {
	client, err := newClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Service{
		cfg:    cfg,
		client: client,
		log:    log,
	}, nil
}

func (s *Service) SendTask(ctx context.Context, task model.Task) error {
	const op = "service.kafka.SendTask"
	log := s.log.With(
		slog.String("op", op),
	)

	taskMessage := TaskMessage{
		Instruction: task.Instruction,
		Code:        task.Code,
		TaskID:      task.ID.String(),
	}

	taskMessageEncoded, err := json.Marshal(taskMessage)
	if err != nil {
		log.Error("marshalling failed", sl.Err(err))
		return err
	}

	err = s.client.SendMessage(ctx, taskMessageEncoded, s.cfg.Topic)
	if err != nil {
		log.Error("failed to send message", sl.Err(err))
		return err
	}

	return nil
}
