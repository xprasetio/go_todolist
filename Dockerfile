# Use the official Golang image as the base image
FROM golang:alpine

# Update the package index and install git
RUN apk update && apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Download necessary Go modules
RUN go mod tidy

# Build the Go application
RUN go build -o binary

# Command to run the executable
ENTRYPOINT ["/app/binary"]
