# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git needed for fetching dependencies
RUN apk add --no-cache git

# Copy both go.mod AND main.go immediately
# This allows 'go mod tidy' to see the imports in main.go
COPY go.mod main.go ./

# Run tidy to automatically find and download the correct dependencies
RUN go mod tidy

# Build the binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -o docker-status-api .

# Final Stage (Tiny image)
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/docker-status-api .

# Expose port
EXPOSE 8682

# Run the binary
CMD ["./docker-status-api"]