# Build stage
FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /build/notification-service .

EXPOSE 4000

CMD ["./notification-service"]