# Use the official Golang image as the base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY farmApp/go.mod farmApp/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Install swag for generating OpenAPI docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy all entries from farmApp to the Working Directory inside the container
COPY farmApp /app/farmApp

# Change working directory to farmApp
WORKDIR /app/farmApp

# Generate the OpenAPI docs with dependencies
RUN swag init -g main.go --parseDependency

# Expose port 8080 to the outside world
EXPOSE 8080