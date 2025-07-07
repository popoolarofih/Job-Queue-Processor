package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"jobqueue/internal/api"
	"jobqueue/internal/tasks"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Asynq client
	tasks.InitAsynqClient()

	// Start the Asynq server (worker) in a separate goroutine
	log.Println("Starting Asynq worker server in a goroutine...")
	go tasks.RunServer()

	// Set up Gin router
	router := gin.Default()

	// API routes
	apiRoutes := router.Group("/api")
	{
		apiRoutes.POST("/jobs", api.SubmitJobHandler)
		apiRoutes.GET("/jobs/status", api.GetJobsStatusHandler)
	}

	// Serve frontend static files
	// Get the directory of the executable
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	staticFilesPath := filepath.Join(exePath, "frontend/build")

	log.Printf("Serving static files from: %s", staticFilesPath)

	// Serve static files from the "frontend/build" directory
	router.StaticFS("/static", http.Dir(filepath.Join(staticFilesPath, "static")))
	router.StaticFile("/", filepath.Join(staticFilesPath, "index.html"))
	router.StaticFile("/favicon.ico", filepath.Join(staticFilesPath, "favicon.ico"))
	router.StaticFile("/logo192.png", filepath.Join(staticFilesPath, "logo192.png"))
    router.StaticFile("/logo512.png", filepath.Join(staticFilesPath, "logo512.png"))
	router.StaticFile("/manifest.json", filepath.Join(staticFilesPath, "manifest.json"))
    router.StaticFile("/asset-manifest.json", filepath.Join(staticFilesPath, "asset-manifest.json"))


	// Handle other routes by serving index.html for client-side routing
	router.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(staticFilesPath, "index.html"))
	})


	// Start Gin server
	log.Println("Starting Gin API server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run Gin server: %v", err)
	}

	// Note: tasks.Client.Close() should be called on graceful shutdown.
}
