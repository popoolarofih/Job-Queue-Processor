package tasks

import "encoding/json"

// TypeGreeting is the type name for the greeting task.
const TypeGreeting = "greeting:sayhello"

// GreetingPayload is the payload for the greeting task.
type GreetingPayload struct {
	Name string
}

// NewGreetingTask creates a new greeting task payload.
func NewGreetingTask(name string) ([]byte, error) {
	payload := GreetingPayload{Name: name}
	return json.Marshal(payload)
}
