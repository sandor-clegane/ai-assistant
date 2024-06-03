package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

//======|| TaskStatus ||========================================

type TaskStatus string

const (
	StatusSending    TaskStatus = "SENDING"
	StatusProcessing TaskStatus = "PROCESSING"
	StatusCompleted  TaskStatus = "COMPLETED"
)

func (st *TaskStatus) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("faile type assertion")
	}

	switch str {
	case string(StatusCompleted):
		*st = StatusCompleted
	case string(StatusProcessing):
		*st = StatusProcessing
	case string(StatusSending):
		*st = StatusSending
	default:
		return errors.New(fmt.Sprintf("unknown enum value %s", str))
	}

	return nil
}

func (st TaskStatus) Value() (driver.Value, error) {
	return string(st), nil
}
