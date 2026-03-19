# Stage 1: Builder
FROM golang:1.26-alpine AS builder

# Install build tools dasar (gcc, make, git jika perlu)
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod dan go.sum duluan untuk caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build API binary
RUN go build -ldflags="-w -s" -o api-server ./cmd/api/main.go

# Build Worker binary
RUN go build -ldflags="-w -s" -o email-worker ./cmd/worker/main.go


# Stage 2: Final Image (Alpine yang ringan)
FROM alpine:latest

# Install tzdata untuk zona waktu dan ca-certificates untuk kirim email
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

# Copy binary dari stage builder ke image final
# COPY .env .
COPY --from=builder /app/api-server .
COPY --from=builder /app/email-worker .

# Secara default, image ini tidak menjalankan apa-apa.
# Kita akan override COMMAND-nya di docker-compose.yml
CMD ["./api-server"]