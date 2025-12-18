# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy workspace files
COPY go.work go.work.sum ./
COPY shared ./shared
COPY services/file-service ./services/file-service

# Build the application
WORKDIR /app/services/file-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Create uploads directory
RUN mkdir -p /app/uploads && chown appuser:appuser /app/uploads

# Copy binary from builder
COPY --from=builder /app/services/file-service/server .
COPY --from=builder /app/services/file-service/migrations ./migrations

# Use non-root user
USER appuser

# Expose port
EXPOSE 9098

# Run application
CMD ["./server"]
