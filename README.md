# Go Job Queue Processor with React Frontend

This project implements a job queue processing system using Go on the backend and React on the frontend. Background jobs are managed using Asynq, a Redis-backed task queue. The Gin framework is used for RESTful APIs, and the frontend provides a dashboard for submitting jobs and viewing their status.

## Features

*   **Backend (Go):**
    *   Gin for RESTful APIs.
    *   Asynq for background job processing.
    *   Redis as the message broker for Asynq.
    *   API endpoint to submit jobs (`POST /api/jobs`).
    *   API endpoint to view job/queue statuses (`GET /api/jobs/status`) using Asynq Inspector.
*   **Frontend (React):**
    *   Dashboard to submit new "greeting" jobs.
    *   Real-time (polling) display of queue and task statuses.
    *   Basic UI for job submission and status monitoring.
*   **Containerization:**
    *   Docker and Docker Compose for easy setup and deployment of the application and Redis.
    *   Multi-stage Dockerfile for optimized Go backend and React frontend build.
    *   Go application serves the built React frontend.

## Project Structure

```
.
├── cmd/server/main.go      # Go application entry point (Gin server, Asynq worker)
├── Dockerfile              # Dockerfile for building the application
├── docker-compose.yml      # Docker Compose for running app and Redis
├── frontend/               # React frontend application
│   ├── public/
│   ├── src/
│   │   ├── App.js          # Main React component
│   │   └── ...
│   ├── package.json
│   └── ...
├── internal/
│   ├── api/                # Gin API handlers
│   │   └── handlers.go
│   └── tasks/              # Asynq task definitions, handlers, client, server
│       ├── client.go
│       ├── handlers.go
│       ├── server.go
│       └── tasks.go
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── README.md               # This file
└── .dockerignore           # Specifies intentionally untracked files for Docker
```

## Prerequisites

*   Docker
*   Docker Compose

## Getting Started

1.  **Clone the repository (if applicable) or ensure all files are in place.**

2.  **Build and run the application using Docker Compose:**
    Open a terminal in the project root directory and run:
    ```bash
    docker-compose up --build
    ```
    This command will:
    *   Build the React frontend.
    *   Build the Go backend application (which includes the frontend assets).
    *   Start the application container and a Redis container.
    *   The `-d` flag can be added (`docker-compose up --build -d`) to run in detached mode.

3.  **Access the application:**
    Open your web browser and navigate to `http://localhost:8080`.

    You should see the Job Queue Dashboard.

4.  **Using the Dashboard:**
    *   **Submit Jobs:** Enter a name in the "Submit New Greeting Job" form and click "Submit Job". The backend worker will process this job (logging a greeting message to the container logs).
    *   **View Statuses:** The "Current Job Statuses" section will display information about the queues and the tasks within them. This section polls the backend every 5 seconds for updates. You can also use the "Refresh Status" button.

5.  **Viewing Logs:**
    *   To view logs for the application container (Go backend, Gin, Asynq worker):
        ```bash
        docker-compose logs -f app
        ```
    *   To view logs for the Redis container:
        ```bash
        docker-compose logs -f redis
        ```

6.  **Stopping the application:**
    Press `Ctrl+C` in the terminal where `docker-compose up` is running. If running in detached mode, use:
    ```bash
    docker-compose down
    ```
    To remove the Redis data volume as well (for a clean restart):
    ```bash
    docker-compose down -v
    ```

## How It Works

1.  The **React frontend** allows users to submit a job type (currently "greeting:sayhello") and a payload (a name).
2.  The job submission request is sent to the **Go backend's** `/api/jobs` endpoint.
3.  The Gin API handler uses the **Asynq client** to enqueue the task into a Redis queue (default queue for greetings).
4.  An **Asynq server (worker)**, running as a goroutine within the same Go application, picks up tasks from the Redis queue.
5.  The worker executes the corresponding task handler (e.g., `HandleGreetingTask`), which in this case logs a greeting message.
6.  The frontend polls the `/api/jobs/status` endpoint. This endpoint uses **Asynq Inspector** to query Redis for current queue statistics and task details, which are then displayed on the dashboard.

## Future Enhancements (Potential)

*   **PostgreSQL Integration:** Store detailed job metadata, logs, and history in PostgreSQL for more robust tracking and querying.
*   **WebSockets:** Implement WebSockets for true real-time updates on the frontend instead of polling.
*   **More Job Types:** Add support for different types of background jobs.
*   **Authentication/Authorization:** Secure the API and frontend.
*   **More Detailed Job View:** Allow users to click on a job to see more details or logs (if stored).
*   **Asynqmon Integration:** Provide instructions or a way to easily run Asynqmon (Asynq's official web UI) for administrative monitoring of queues.
*   **Graceful Shutdown:** Implement graceful shutdown for the Go server and Asynq workers.
*   **Unit and Integration Tests:** Add comprehensive automated tests.
*   **Configuration Management:** More sophisticated configuration management beyond environment variables for Redis.