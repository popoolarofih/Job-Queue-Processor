package tasks

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

// HandleGreetingTask handles the greeting task.
// It logs a greeting message.
func HandleGreetingTask(ctx context.Context, t *asynq.Task) error {
	var p GreetingPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	log.Printf("Hello, %s!", p.Name)
	return nil
}
