# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git needed for fetching dependencies
RUN apk add --no-cache git

# Copy module files first to cache dependencies
COPY go.mod ./
# Note: In a real scenario, you would run 'go mod tidy' to generate go.sum
# creating a dummy go.sum here to satisfy build if it's missing, 
# strictly for this self-contained example.
RUN touch go.sum && go mod download

# Copy source code
COPY main.go .

# Build the binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -o docker-status-api .

# Final Stage (Tiny image)
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/docker-status-api .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./docker-status-api"]