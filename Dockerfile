FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
# RUN apk add --no-cache curl
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
# COPY --from=builder /app/migrate.linux-amd64 /usr/bin/migrate
COPY .container.env .env
COPY start.sh .
COPY db/migrations ./db/migrations

EXPOSE 8080
# RUN migrate -path /app/db/migrations -database "${DBURL}" -verbose up
# CMD ["/app/main"]
# ENTRYPOINT [ "/app/start.sh" ]