FROM golang:1.23.4-alpine as builder
ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build /app/cmd/rivalslfg

FROM scratch as release
COPY --from=builder /app/rivalslfg /rivalslfg
ENTRYPOINT ["/rivalslfg"]