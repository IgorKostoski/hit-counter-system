FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /hit-counter-api .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /hit-counter-api /app/hit-counter-api

EXPOSE 8080

ENV APP_PORT=8080
ENV GIN_MODE=release

CMD ["/app/hit-counter-api"]