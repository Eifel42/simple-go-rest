# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY ../go.mod ../farmApp/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Install swag for generating OpenAPI docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy all entries from farmApp to the Working Directory inside the container
COPY .. .

# Generate the OpenAPI docs
RUN swag init -g main.go

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]