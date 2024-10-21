# Use the official Golang image as a base image
FROM golang:1.22-alpine

RUN apk add --no-cache curl

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Install CompileDaemon for hot-reloading
RUN go install github.com/githubnemo/CompileDaemon@latest

# Install golangci-lint
# binary will be $(go env GOPATH)/bin/golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0

# Command to run the executable with CompileDaemon
ENTRYPOINT ["/go/bin/CompileDaemon", "--build=go build -o binary .", "--command=./binary"]
