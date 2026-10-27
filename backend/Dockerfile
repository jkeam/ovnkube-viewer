# Build stage
FROM registry.redhat.io/ubi9/go-toolset:9.8-1793107762 AS builder

# Set working directory
WORKDIR /opt/app-root/src

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 creates a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage - use UBI9 micro for minimal footprint
FROM registry.redhat.io/ubi9/ubi-micro:9.8-1791895366

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /opt/app-root/src/server /app/server

# Run as non-root user
USER 1001

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/server"]
