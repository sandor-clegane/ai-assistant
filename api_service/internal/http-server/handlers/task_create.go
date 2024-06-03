package handlers

import (
	"ai-assistant-api/internal/utils/logger/sl"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type TaskCreator interface {
	CreateTask(ctx context.Context, instruction, code string) (uuid.UUID, error)
}

type taskCreateRequest struct {
	Instruction string `json:"instruction"`
	Code        string `json:"code"`
}

func HandleCreateTask(log *slog.Logger, taskCreator TaskCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.HandleCreateTask"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		taskData := taskCreateRequest{}
		err := json.NewDecoder(r.Body).Decode(&taskData)
		if err != nil {
			log.Error("failed to unmarshall request", sl.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		createdTask, err := taskCreator.CreateTask(r.Context(), taskData.Instruction, taskData.Code)
		if err != nil {
			log.Error("failed to save pull request", sl.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedResponse, err := json.Marshal(createdTask)
		if err != nil {
			log.Error("response marshalling failed", sl.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(encodedResponse)
		w.WriteHeader(http.StatusOK)
	}
}
