package handlers

import (
	"ai-assistant-api/internal/model"
	"ai-assistant-api/internal/service"
	"ai-assistant-api/internal/utils/logger/sl"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type TaskGetter interface {
	GetTask(ctx context.Context, taskID uuid.UUID) (*model.Task, error)
}

func HandleGetTask(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.HandleGetTask"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		strTaskID := chi.URLParam(r, "id")
		taskID, err := uuid.Parse(strTaskID)
		if err != nil {
			log.Info("task id is invalid")
			http.Error(w, "invalud request", http.StatusBadRequest)
			return
		}

		task, err := taskGetter.GetTask(context.Background(), taskID)
		if err != nil {
			log.Error("failed to get task", sl.Err(err))
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		taskEncoded, err := json.Marshal(task)
		if err != nil {
			log.Error("response marshalling failed", sl.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(taskEncoded)
		w.WriteHeader(http.StatusOK)
	}
}
