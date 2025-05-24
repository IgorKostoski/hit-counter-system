FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /hit-counter-api .


FROM golang:1.23-alpine AS development

WORKDIR /app

RUN apk add --no-cache git

RUN go install github.com/cosmtrek/air@v1.51.0

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]

FROM alpine:latest

WORKDIR /app

COPY --from=builder /hit-counter-api /app/hit-counter-api

EXPOSE 8080

ENV APP_PORT=8080
ENV GIN_MODE=release

CMD ["/app/hit-counter-api"]