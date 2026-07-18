# =========================
# Stage 1 - Build
# =========================
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd

# =========================
# Stage 2 - Runtime
# =========================
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata wget

COPY --from=builder /app/backend .

EXPOSE 8080

CMD ["./backend"]
