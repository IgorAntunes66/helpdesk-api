FROM golang:1.24-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/users-server ./users-service/cmd/api

FROM alpine:latest

WORKDIR /app

COPY db/migrations ./db/migrations

COPY --from=builder /app/users-server .

EXPOSE 8082

CMD [ "./users-server" ]