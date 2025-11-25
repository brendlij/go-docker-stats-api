package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
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

var cli *client.Client

func init() {
	var err error
	// Initialize the Docker client
	// client.FromEnv will automatically pick up the socket mounted at /var/run/docker.sock
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create docker client: %v", err)
	}
}

func main() {
	defer cli.Close()

	// Handler for the container status endpoint
	http.HandleFunc("/api/containers", handleContainers)
	http.HandleFunc("/api/containers/", handleContainerDetail)
	http.HandleFunc("/health", handleHealth)

	log.Println("Starting Docker Status API on :8911")
	log.Fatal(http.ListenAndServe(":8911", nil))
}

// handleHealth returns a simple health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleContainers returns all container statuses
func handleContainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// List all containers (running and stopped)
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		http.Error(w, `{"error": "Failed to list containers"}`, http.StatusInternalServerError)
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

		// Get health status
		health := "unknown"
		if ctr.State == "running" {
			// Try to get detailed health info
			inspect, err := cli.ContainerInspect(context.Background(), ctr.ID)
			if err == nil && inspect.State.Health != nil {
				health = inspect.State.Health.Status
			}
		}

		stats = append(stats, ContainerStatus{
			ID:     ctr.ID[:12],
			Name:   name,
			Image:  ctr.Image,
			State:  ctr.State,
			Status: ctr.Status,
			Health: health,
		})
	}

	json.NewEncoder(w).Encode(stats)
}

// handleContainerDetail returns status for a specific container
func handleContainerDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract container ID from path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, `{"error": "Invalid container ID"}`, http.StatusBadRequest)
		return
	}

	containerID := parts[3]
	if containerID == "" {
		http.Error(w, `{"error": "Container ID required"}`, http.StatusBadRequest)
		return
	}

	// Get container info
	inspect, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		http.Error(w, `{"error": "Container not found"}`, http.StatusNotFound)
		log.Printf("Error inspecting container: %v", err)
		return
	}

	name := strings.TrimPrefix(inspect.Name, "/")
	health := "unknown"
	if inspect.State.Health != nil {
		health = inspect.State.Health.Status
	}

	status := ContainerStatus{
		ID:     inspect.ID[:12],
		Name:   name,
		Image:  inspect.Config.Image,
		State:  inspect.State.Status,
		Status: getStatusString(inspect.State),
		Health: health,
	}

	json.NewEncoder(w).Encode(status)
}

// getStatusString returns a human-readable status string
func getStatusString(state *types.ContainerState) string {
	if state.Running {
		return "running"
	} else if state.Paused {
		return "paused"
	} else if state.ExitCode == 0 {
		return "exited (0)"
	}
	return "exited (" + string(rune(state.ExitCode)) + ")"
}
