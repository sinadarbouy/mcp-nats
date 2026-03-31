# Build stage
FROM --platform=$BUILDPLATFORM golang:1.25-bookworm AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Set build arguments for cross-compilation
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Build the application with proper cross-compilation
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o mcp-nats ./cmd/mcp-nats

# Install NATS CLI
RUN go install github.com/nats-io/natscli/nats@latest

# Final stage
FROM debian:bookworm-slim

# Create a non-root user
RUN useradd -r -u 1000 -m mcp-nats

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder --chown=1000:1000 /app/mcp-nats /app/
COPY --from=builder --chown=1000:1000 /go/bin/nats /usr/local/bin/

# Use the non-root user
USER mcp-nats

# Expose the port the app runs on
EXPOSE 8000

# Run the application
ENTRYPOINT ["/app/mcp-nats", "--transport", "sse", "--sse-address", "0.0.0.0:8000"]
