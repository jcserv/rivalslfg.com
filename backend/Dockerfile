FROM golang:1.23.4-alpine as builder
ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build /app/cmd/rivalslfg

FROM alpine:latest as release
WORKDIR /app

# Install migrate binary
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz \
    && tar xvz migrate.linux-amd64.tar.gz \
    && mv migrate.linux-amd64 /usr/local/bin/migrate \
    && chmod +x /usr/local/bin/migrate

COPY --from=builder /app/rivalslfg /rivalslfg
COPY --from=builder /app/db/migrations ./db/migrations

ENTRYPOINT ["/rivalslfg"]