# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary dari path cmd/api/main.go
# Kita beri nama outputnya 'binary-app'
RUN CGO_ENABLED=0 GOOS=linux go build -o /binary-app ./cmd/api

# Stage 2: Final Image (Kecil & Aman)
FROM alpine:latest
WORKDIR /root/

# Ambil binary dari stage builder
COPY --from=builder /binary-app .

EXPOSE 8080

# Jalankan binary
CMD ["./binary-app"]