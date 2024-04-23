FROM golang:1.21-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app cmd/validator/main.go

FROM alpine:latest

ARG RPC_DIAL_URL
ENV RPC_DIAL_URL=$RPC_DIAL_URL

WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]