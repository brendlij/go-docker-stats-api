# Docker Status API

A lightweight, production-ready Go API service that provides real-time Docker container status monitoring. Runs as a containerized microservice with secure access to the Docker daemon.

**Perfect for**: Container management systems like Dockge, monitoring dashboards, orchestration platforms, and CI/CD pipelines.

---

## Features

- **Container Discovery** - List all containers (running, stopped, paused)
- **Real-time Status** - Get current state and uptime for each container
- **Health Monitoring** - Track container health check status
- **Detailed Insights** - Access image details, state information, and metadata
- **Lightweight** - ~15MB Docker image with minimal resource footprint
- **Secure** - Read-only Docker socket access
- **REST API** - Clean, JSON-based endpoints for easy integration
- **Health Checks** - Built-in endpoint for orchestration health verification

---

## Quick Start

### Using Docker Compose (Recommended)

```bash
docker-compose up -d
```

The API will be available at `http://localhost:8911`

### Using Docker CLI

```bash
docker run -d \
  --name docker-status-api \
  -p 8911:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  go-docker-stats-api
```

### Local Development

```bash
# Install dependencies
go mod download

# Run the application
go run main.go
```

---

## API Endpoints

### Get All Containers

**Request:**

```http
GET /api/containers
```

**Response:** JSON array of all containers

```json
[
  {
    "id": "abc123def456",
    "name": "dockge",
    "image": "louislam/dockge:latest",
    "state": "running",
    "status": "Up 2 hours",
    "health": "healthy"
  },
  {
    "id": "xyz789abc123",
    "name": "nginx",
    "image": "nginx:latest",
    "state": "running",
    "status": "Up 1 day",
    "health": "unknown"
  }
]
```

### Get Specific Container Status

**Request:**

```http
GET /api/containers/{container_id_or_name}
```

**Response:** Detailed container information

```json
{
  "id": "abc123def456",
  "name": "dockge",
  "image": "louislam/dockge:latest",
  "state": "running",
  "status": "running",
  "health": "healthy"
}
```

### Health Check

**Request:**

```http
GET /health
```

**Response:** API health status

```json
{
  "status": "ok"
}
```

---

## Usage Examples

### cURL

```bash
# Get all containers
curl http://localhost:8911/api/containers

# Get specific container
curl http://localhost:8911/api/containers/dockge

# Check API health
curl http://localhost:8911/health
```

### JavaScript/TypeScript

```javascript
// Get all containers
const response = await fetch("http://localhost:8911/api/containers");
const containers = await response.json();
console.log(containers);

// Get specific container
const containerResponse = await fetch(
  "http://localhost:8911/api/containers/dockge"
);
const container = await containerResponse.json();
console.log(container);
```

### Python

```python
import requests

# Get all containers
response = requests.get('http://localhost:8911/api/containers')
containers = response.json()
print(containers)

# Get specific container
response = requests.get('http://localhost:8911/api/containers/dockge')
container = response.json()
print(container)
```

---

## Configuration

### Port Configuration

The API runs on port **8911** by default. To change it:

1. **Dockerfile** - Update the `EXPOSE` directive:

   ```dockerfile
   EXPOSE 9000
   ```

2. **main.go** - Update the listen address:

   ```go
   log.Fatal(http.ListenAndServe(":9000", nil))
   ```

3. **compose.yaml** - Update the port mapping:
   ```yaml
   ports:
     - "9000:8080"
   ```

### Docker Socket Access

For the container to access the Docker daemon, it must:

1. Have the Docker socket mounted as a volume: `/var/run/docker.sock:/var/run/docker.sock:ro`
2. The socket should be mounted read-only (`:ro`) for security

## Requirements

- Docker & Docker Compose (for running the container)
- Go 1.21+ (for local development)
- Docker daemon accessible via `/var/run/docker.sock`

## Example Usage with curl

```bash
# Get all containers
curl http://localhost:8080/api/containers

# Get specific container
curl http://localhost:8080/api/containers/dockge

# Health check
curl http://localhost:8080/health
```

## Example Usage in JavaScript/TypeScript

```javascript
// Get all containers
const response = await fetch("http://localhost:8080/api/containers");
const containers = await response.json();
console.log(containers);

// Get specific container
const containerResponse = await fetch(
  "http://localhost:8080/api/containers/dockge"
);
const container = await containerResponse.json();
console.log(container);
```

## Configuration

The API runs on port `8080` by default. You can change this in:

- `main.go`: Line where `http.ListenAndServe(":8080", nil)` is called
- `Dockerfile`: Adjust the `EXPOSE` port
- `compose.yaml`: Adjust the `ports` mapping

## Security Considerations

- The Docker socket is mounted as read-only, preventing accidental modifications
- The API has no authentication - consider adding an API key or reverse proxy if exposed externally
- Only request container information from trusted sources

## License

MIT
