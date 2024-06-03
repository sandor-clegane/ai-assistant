package task

import (
	"ai-assistant-api/internal/model"
	"ai-assistant-api/internal/service"
	"ai-assistant-api/internal/storage"
	"ai-assistant-api/internal/utils/logger/sl"
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

type TaskStorage interface {
	SaveTask(ctx context.Context, task model.Task) (*model.Task, error)
	GetTask(ctx context.Context, taskID uuid.UUID) (*model.Task, error)
	TaskListByStatus(ctx context.Context, status model.TaskStatus) ([]model.Task, error)
}

type KafkaService interface {
	SendTask(ctx context.Context, task model.Task) error
}

type Service struct {
	log          *slog.Logger
	taskStorage  TaskStorage
	kafkaService KafkaService
}

func New(log *slog.Logger, storage TaskStorage, kafkaService KafkaService) *Service {
	return &Service{
		log:          log,
		taskStorage:  storage,
		kafkaService: kafkaService,
	}
}

func (s *Service) CreateTask(ctx context.Context, instruction, code string) (uuid.UUID, error) {
	const op = "service.task.SaveTask"
	log := s.log.With(
		slog.String("op", op),
	)

	task := model.Task{
		Instruction: instruction,
		Code:        code,
		Status:      model.StatusSending,
	}

	createdTask, err := s.taskStorage.SaveTask(ctx, task)
	if err != nil {
		log.Error("failed to save pull reqest", sl.Err(err))
		return uuid.UUID{}, err
	}

	return createdTask.ID, nil
}

func (s *Service) GetTask(ctx context.Context, taskID uuid.UUID) (*model.Task, error) {
	const op = "service.task.SaveTask"
	log := s.log.With(
		slog.String("op", op),
	)

	task, err := s.taskStorage.GetTask(ctx, taskID)
	if err != nil {
		log.Error("failed to find task by ID", sl.Err(err), slog.Any("id", taskID))
		if errors.Is(err, storage.ErrNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, err
	}

	return task, nil
}

func (s *Service) TaskListByStatus(ctx context.Context, status model.TaskStatus) ([]model.Task, error) {
	const op = "service.task.TaskListByStatus"
	log := s.log.With(
		slog.String("op", op),
	)

	taskList, err := s.taskStorage.TaskListByStatus(ctx, status)
	if err != nil {
		log.Error("failed to get task list", sl.Err(err))
		return nil, err
	}

	return taskList, nil
}

func (s *Service) SendTask(ctx context.Context, task model.Task) error {
	const op = "service.task.SendTask"
	log := s.log.With(
		slog.String("op", op),
	)

	err := s.kafkaService.SendTask(ctx, task)
	if err != nil {
		log.Error("send task failed", sl.Err(err))
		return err
	}

	task.Status = model.StatusProcessing

	_, err = s.taskStorage.SaveTask(ctx, task)
	if err != nil {
		log.Error("failed to save updated task", sl.Err(err))
		return err
	}

	return nil
}

func (s *Service) FinishTask(ctx context.Context, taskID uuid.UUID, response string) error {
	const op = "service.task.FinishTask"
	log := s.log.With(
		slog.String("op", op),
	)

	task := model.Task{
		ID:       taskID,
		Status:   model.StatusCompleted,
		Response: response,
	}

	_, err := s.taskStorage.SaveTask(ctx, task)
	if err != nil {
		log.Error("failed to save updated task", sl.Err(err))
		return err
	}

	return nil
}
