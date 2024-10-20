# Use the official Golang image as a base image
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Command to run the executable
RUN go install github.com/githubnemo/CompileDaemon@latest

ENTRYPOINT ["/go/bin/CompileDaemon", "--build=go build -o binary .", "--command=./binary"]
