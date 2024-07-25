# Use the official Golang image to create a build artifact.
FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Change directory to cms
WORKDIR /app/cms

# Build the Go app
RUN go build -o /app/main

FROM ubuntu:latest

RUN apt-get update --fix-missing && apt-get install -y software-properties-common ffmpeg curl

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /app/

COPY .env /app/
RUN mkdir /app/tmp

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
