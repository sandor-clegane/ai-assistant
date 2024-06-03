package sender

type TaskMessage struct {
	Instruction string `json:"instruction"`
	Code        string `json:"code"`
	TaskID      string `json:"task_id"`
}
