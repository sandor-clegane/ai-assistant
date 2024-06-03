package background

import (
	"ai-assistant-api/internal/config"
	"ai-assistant-api/internal/utils/logger/sl"
	"context"
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

func (s *Service) startConsumer(ctx context.Context, done chan struct{}) error {
	config := sarama.NewConfig()
	consumerGroup, err := sarama.NewConsumerGroup([]string{s.cfg.ServerAddress}, s.cfg.ConsumerGroup, config)
	if err != nil {
		return err
	}

	consumer, err := newConsumer(s.cfg, s.taskService, s.log)
	if err != nil {
		return err
	}

	consumerCtx, cancel := context.WithCancel(ctx)
	go func() {
		<-done
		consumerGroup.Close()
		cancel()
	}()

	go func() {
		for {
			err = consumerGroup.Consume(consumerCtx, []string{s.cfg.ConsumerTopic}, consumer)
			if err != nil {
				s.log.Error("consume failed", sl.Err(err))
			}
			if consumerCtx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

type consumer struct {
	taskService   TaskService
	consumerGroup sarama.ConsumerGroup
	log           *slog.Logger
}

func newConsumer(cfg *config.Config, taskService TaskService, log *slog.Logger) (*consumer, error) {
	config := sarama.NewConfig()
	consumerGroup, err := sarama.NewConsumerGroup([]string{cfg.ServerAddress}, cfg.ConsumerGroup, config)
	if err != nil {
		return nil, err
	}

	return &consumer{
		consumerGroup: consumerGroup,
		taskService:   taskService,
		log:           log,
	}, nil
}

func (*consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

type taskResponseMessage struct {
	TaskID   string `json:"task_id"`
	Response string `json:"response"`
}

func (c *consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var data taskResponseMessage
		err := json.Unmarshal(msg.Value, &data)
		if err != nil {
			c.log.Error("unkarshal failed", sl.Err(err))
			continue
		}

		taskID, err := uuid.Parse(data.TaskID)
		if err != nil {
			c.log.Error("failed to parse task id", sl.Err(err))
			continue
		}

		c.taskService.FinishTask(context.Background(), taskID, data.Response)
		sess.MarkMessage(msg, "")
	}

	return nil
}
