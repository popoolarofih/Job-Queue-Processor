package api

import (
	"jobqueue/internal/tasks"
	"log"
	"net/http"
	"os" // Added import for os
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

// JobSubmissionRequest is the request body for submitting a new job.
type JobSubmissionRequest struct {
	Type    string                 `json:"type" binding:"required"`
	Payload map[string]interface{} `json:"payload" binding:"required"` // Flexible payload
}

// SubmitJobHandler handles requests to submit a new job.
func SubmitJobHandler(c *gin.Context) {
	var req JobSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type == tasks.TypeGreeting {
		name, ok := req.Payload["name"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "payload for greeting task must contain a 'name' field of type string"})
			return
		}

		taskInfo, err := tasks.EnqueueGreetingTask(name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue task", "details": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "job submitted successfully", "job_id": taskInfo.ID, "queue": taskInfo.Queue, "type": taskInfo.Type})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported job type"})
	}
}

// GetJobsStatusHandler provides information about jobs and queues.
func GetJobsStatusHandler(c *gin.Context) {
	// Use the same mechanism to get Redis address
	redisAddr := tasks.RedisAddr // Default
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		redisAddr = addr
	}

	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: redisAddr})
	defer inspector.Close()

	queues, err := inspector.Queues()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get queues from inspector", "details": err.Error()})
		return
	}

	type JobInfo struct {
		ID      string    `json:"id"`
		Type    string    `json:"type"`
		Payload string    `json:"payload"`
		Queue   string    `json:"queue"`
		State   string    `json:"state"`
		NextTry time.Time `json:"next_try,omitempty"`
		LastErr string    `json:"last_err,omitempty"`
	}

	type QueueInfo struct {
		Name         string    `json:"name"`
		Size         int       `json:"size"`
		Active       int       `json:"active"`
		Pending      int       `json:"pending"`
		Aggregating  int       `json:"aggregating"`
		Scheduled    int       `json:"scheduled"`
		Retry        int       `json:"retry"`
		Archived     int       `json:"archived"`
		Completed    int       `json:"completed"`
		Processed    int64     `json:"processed"`
		Failed       int64     `json:"failed"`
		LastFullScan time.Time `json:"last_full_scan"`
		Tasks        []JobInfo `json:"tasks"`
	}

	var responseQueues []QueueInfo

	for _, qname := range queues {
		qinfo, err := inspector.GetQueueInfo(qname)
		if err != nil {
			log.Printf("Error getting info for queue %s: %v", qname, err)
			continue // Skip this queue on error
		}

		currentQueueInfo := QueueInfo{
			Name:         qinfo.Queue,
			Size:         qinfo.Size,
			Active:       qinfo.Active,
			Pending:      qinfo.Pending,
			Aggregating:  qinfo.Aggregating,
			Scheduled:    qinfo.Scheduled,
			Retry:        qinfo.Retry,
			Archived:     qinfo.Archived,
			Completed:    qinfo.Completed,
			Processed:    qinfo.Processed,
			Failed:       qinfo.Failed,
			LastFullScan: qinfo.LastFullScan,
			Tasks:        []JobInfo{},
		}

		// Fetch a few tasks from different states for a snapshot
		taskFetchers := map[string]func(string, ...asynq.ListOption) ([]*asynq.TaskInfo, error){
			"active":    inspector.ListActiveTasks,
			"pending":   inspector.ListPendingTasks,
			"retry":     inspector.ListRetryTasks,
			"archived":  inspector.ListArchivedTasks,
			"completed": inspector.ListCompletedTasks, // Added completed
		}

		opts := []asynq.ListOption{asynq.PageSize(5)} // Get up to 5 tasks per category

		for state, fetcher := range taskFetchers {
			tasks, err := fetcher(qname, opts...)
			if err != nil {
				log.Printf("Error listing %s tasks for queue %s: %v", state, qname, err)
				continue
			}
			for _, t := range tasks {
				currentQueueInfo.Tasks = append(currentQueueInfo.Tasks, JobInfo{
					ID:      t.ID,
					Type:    t.Type,
					Payload: string(t.Payload),
					Queue:   t.Queue,
					State:   state, // Use the map key as state
					NextTry: t.NextRetry,
					LastErr: t.LastErr,
				})
			}
		}
		responseQueues = append(responseQueues, currentQueueInfo)
	}

	c.JSON(http.StatusOK, responseQueues)
}
