# Stage 1: Build the application
FROM golang:1.23.2 AS builder

WORKDIR /app

# Copy source code
COPY etcd-snapshot-to-json.go .

# Initialize Go modules
RUN go mod init etcd-snapshot-to-json
RUN go mod tidy

# Build the Go application as a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o etcd-snapshot-to-json etcd-snapshot-to-json.go
#RUN go build -o etcd-snapshot-to-json etcd-snapshot-to-json.go

# Stage 2: Use a minimal scratch image to run the static binary
FROM scratch

# Copy the built binary from the builder stage
COPY --from=builder /app/etcd-snapshot-to-json /etcd-snapshot-to-json

# Set the entrypoint to the executable
ENTRYPOINT ["/etcd-snapshot-to-json"]
