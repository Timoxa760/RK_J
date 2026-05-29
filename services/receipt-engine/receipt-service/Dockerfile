FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /service ./services/receipt-engine/receipt-service

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /service .
EXPOSE 8002
CMD ["./service"]
