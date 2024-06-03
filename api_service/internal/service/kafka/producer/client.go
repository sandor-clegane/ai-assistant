package sender

import (
	"ai-assistant-api/internal/config"
	"context"

	"github.com/IBM/sarama"
)

type client struct {
	producer sarama.SyncProducer
}

func newClient(cfg *config.Config) (*client, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{cfg.ServerAddress}, config)
	if err != nil {
		return nil, err
	}

	return &client{
		producer: producer,
	}, nil
}

func (c *client) SendMessage(ctx context.Context, msg []byte, topic string) error {
	saramaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	_, _, err := c.producer.SendMessage(saramaMsg)
	return err
}
