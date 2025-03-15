# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install necessary dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Copy the vendor directory to use with -mod=vendor
COPY vendor ./vendor

# Copy the source code
COPY . .

# Build the application using the vendor directory
RUN go build -mod=vendor -o app ./cmd/rest/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app .

# Run the application
CMD ["/app/app"]

# sudo docker tag <imageID> ewanlav/diploma:latest