# Start with a base Go image
FROM golang:1.23-alpine as builder

# Set the Current Working Directory inside the container to /app
WORKDIR /app

# Copy the Go Modules manifests (go.mod and go.sum)
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod tidy

# Copy the entire Go project into the container
COPY . .

# Build the Go app (specify the path to main.go in cmd/api/)
RUN go build -o main ./cmd/api/

# Start a new stage to create a smaller image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary from the builder stage
COPY --from=builder /app/main .

# Copy your .env file (if you have one) for environment variables
COPY .env ./

# Expose the port that your app will run on
EXPOSE 8080

# Run the Go application
CMD ["./main"]
