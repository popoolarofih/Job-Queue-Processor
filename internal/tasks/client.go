package tasks

import (
	"log"
	"os"

	"github.com/hibiken/asynq"
)

var RedisAddr = "127.0.0.1:6379" // Default Redis address

var Client *asynq.Client

func getRedisAddr() string {
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		return addr
	}
	return RedisAddr
}

func InitAsynqClient() {
	currentRedisAddr := getRedisAddr()
	Client = asynq.NewClient(asynq.RedisClientOpt{Addr: currentRedisAddr})
	log.Printf("Asynq client initialized with Redis address: %s", currentRedisAddr)
}

// EnqueueGreetingTask enqueues a new greeting task.
func EnqueueGreetingTask(name string) (*asynq.TaskInfo, error) {
	payload, err := NewGreetingTask(name)
	if err != nil {
		return nil, err
	}
	task := asynq.NewTask(TypeGreeting, payload)
	info, err := Client.Enqueue(task)
	if err != nil {
		return nil, err
	}
	log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return info, nil
}
