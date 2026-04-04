# syntax=docker/dockerfile:1

# --- Build stage ---
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# --- Runtime stage ---
FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
RUN mkdir -p /data && chown 65534:65534 /data
EXPOSE 9092
VOLUME ["/data"]
ENV DB_PATH=/data/dfx.db
USER 65534
CMD ["./server"]