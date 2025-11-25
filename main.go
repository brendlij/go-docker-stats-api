package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// ContainerStatus represents the simplified data we want to return
type ContainerStatus struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	State   string `json:"state"`  // e.g., "running", "exited"
	Status  string `json:"status"` // e.g., "Up 2 hours"
	Health  string `json:"health,omitempty"`
}

func main() {
	// Initialize the Docker client
	// client.FromEnv will automatically pick up the socket mounted at /var/run/docker.sock
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create docker client: %v", err)
	}
	defer cli.Close()

	// Handler for the status endpoint
	http.HandleFunc("/api/containers", func(w http.ResponseWriter, r *http.Request) {
		// Set JSON content type
		w.Header().Set("Content-Type", "application/json")

		// List all containers (running and stopped)
		containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error listing containers: %v", err)
			return
		}

		// Transform the raw docker data into our simplified struct
		var stats []ContainerStatus
		for _, ctr := range containers {
			// Get the first name (docker names start with /)
			name := "unknown"
			if len(ctr.Names) > 0 {
				name = strings.TrimPrefix(ctr.Names[0], "/")
			}

			// Check health status if available
			health := ""
			// The status string often contains health info like "Up 2 hours (healthy)"
			if strings.Contains(ctr.Status, "(healthy)") {
				health = "healthy"
			} else if strings.Contains(ctr.Status, "(unhealthy)") {
				health = "unhealthy"
			} else if strings.Contains(ctr.Status, "(starting)") {
				health = "starting"
			}

			stats = append(stats, ContainerStatus{
				ID:     ctr.ID[:12], // Short ID
				Name:   name,
				Image:  ctr.Image,
				State:  ctr.State,
				Status: ctr.Status,
				Health: health,
			})
		}

		// Encode to JSON
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	})

	// simple health check for the app itself
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}