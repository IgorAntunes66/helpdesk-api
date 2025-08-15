FROM golang:1.24-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/tickets-server ./tickets-service/cmd/api

FROM alpine:latest

WORKDIR /app

COPY db/migrations ./db/migrations

COPY --from=builder /app/tickets-server .

EXPOSE 8080

CMD [ "./tickets-server" ]