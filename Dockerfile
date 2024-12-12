# Use the specific Go image version for building the artifact
FROM golang:1.23.4-alpine as builder

# Create and change to the app directory.
WORKDIR /app

# Copy local package files to the container's workspace.
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy all directories at this level into the container
COPY . .

# Build the binary. This assumes main.go is still the entry point in cmd.
RUN go build -o /todo-list ./cmd

# Use a minimal alpine image for the runtime.
FROM alpine:latest  

# Install CA certificates.
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage.
COPY --from=builder /todo-list /todo-list

# Run the binary.
CMD ["/todo-list"]