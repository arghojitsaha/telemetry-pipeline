# Stage 1: Build the Go binary using Alpine for a lightweight footprint
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency files first to leverage Docker layer caching
COPY go.mod ./
RUN go mod download

# Copy source code and build a statically linked binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o telemetry-service main.go

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/telemetry-service .

EXPOSE 8080

CMD ["./telemetry-service"]