package notification

import (
	"ai-assistant-api/internal/config"
	"ai-assistant-api/internal/utils/logger/sl"
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

const (
	chatID = "447006349"
)

type Service struct {
	cfg *config.Config
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *Service {
	return &Service{
		cfg: cfg,
		log: log,
	}
}

func (s *Service) Notify(message string) error {
	s.log.Info("Preparing to send Telegram notification")

	type reqBody struct {
		ChatId string `json:"chat_id"`
		Text   string `json:"text"`
	}

	var reqBodyValue reqBody
	reqBodyValue.ChatId = chatID
	reqBodyValue.Text = message

	reqBodyValueJson, err := json.Marshal(reqBodyValue)
	if err != nil {
		s.log.Error("Failed to send Telegram notification", sl.Err(err))
		return err
	}

	requestBodyJson := bytes.NewReader(reqBodyValueJson)

	apiUrl := "https://api.telegram.org/bot" + s.cfg.TelegramSecret + "/sendMessage"

	req, err := http.NewRequest(http.MethodPost, apiUrl, requestBodyJson)
	if err != nil {
		s.log.Error("Failed to send Telegram notification", sl.Err(err))
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Create client
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	res, err := client.Do(req)
	if err != nil {
		s.log.Error("Failed to send Telegram notification", sl.Err(err))
		return err
	}

	if res.StatusCode/100 != 2 {
		err = errors.New("Failed to send Telegram notification")
		s.log.Error("Failed to send Telegram notification", sl.Err(err))
		return err
	}

	s.log.Info("Notification send")

	return nil
}
