package tasks

import (
	"log"
	// os is already imported in client.go, which is part of the same package.
	// We'll use the getRedisAddr from there.

	"github.com/hibiken/asynq"
)

// RunServer starts an Asynq server (worker).
func RunServer() {
	currentRedisAddr := getRedisAddr() // Uses the function from client.go
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: currentRedisAddr},
		asynq.Config{
			// Specify how many concurrent workers to run.
			Concurrency: 10,
			// Optionally specify basic priority mapping.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeGreeting, HandleGreetingTask)

	log.Println("Asynq server (worker) starting...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run asynq server: %v", err)
	}
}
