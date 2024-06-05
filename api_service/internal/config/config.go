package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TelegramSecret string `yaml:"tg_secret" env-required:"true"`
	HttpRequestsMax int `yaml:"http_requests_max" env-default:"20"`
	Env          string `yaml:"env" env-default:"local"`
	StoragePath  string `yaml:"storage_path" env-required:"true"`
	GithubSecret string `yaml:"github_secret" env-required:"true"`

	HTTPServer `yaml:"http_server"`
	Kafka      `yaml:"kafka"`
	Background `yaml:"background"`
}

type HTTPServer struct {
	Address         string        `yaml:"address" env-default:"localhost:8080"`
	Timeout         time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
}

type Kafka struct {
	Topic         string `yaml:"topic" env-default:"ai_assistant_tasks"`
	ServerAddress string `yaml:"server_address" env-default:"localhost:9092"`
	ConsumerTopic string `yaml:"consumer_topic" env-default:"ai_assistant_responses"`
	ConsumerGroup string `yaml:"consumer_group" env-default:"ai_assistant_responses_group"`
}

type Background struct {
	SendWorkerTimeout time.Duration `yaml:"send_worker_timeout" env-default:"10s"`
}

func MustLoad(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
