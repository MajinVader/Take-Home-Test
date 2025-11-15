# Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod/sum dan download dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

# Run
FROM alpine:latest

WORKDIR /app

# Copy binary dari stage builder
COPY --from=builder /app/app .

# Port yang dipakai Fiber
EXPOSE 3000

# Jalankan app
CMD ["./app"]
