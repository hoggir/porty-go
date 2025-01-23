# Use the official Golang image as the base image
FROM golang:1.23.2

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/main.go
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go
# RUN swag init -g cmd/main.go update swagger

# Expose the port the application will run on
EXPOSE 8000

# Set the entry point for the container
CMD ["./main"]