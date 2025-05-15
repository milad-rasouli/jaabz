# sudo docker build -t jaabz:0.0.1 .
# Stage 1: Build the Go binary
FROM golang:1.24.2-alpine AS build

# Install only essential build dependencies
RUN apk add --no-cache \
    build-base \
    gcompat \
    librdkafka-dev \
    libc6-compat

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -buildvcs=false -trimpath \
    -ldflags="-w -s" -o /app/build ./cmd/.

# Stage 2: Create minimal production image
FROM alpine:latest AS prod

# Install only runtime dependencies
RUN apk add --no-cache libc6-compat

# Set working directory
WORKDIR /app

# Copy built binary from build stage
COPY --from=build /app/build .

# Ensure binary is executable
RUN chmod +x /app/build

# Set entrypoint
ENTRYPOINT ["/app/build"]
