# STAGE 1
FROM golang:alpine as builder

LABEL maintainer="humanbelnik"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && \
    apk --update add ca-certificates git

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/auth/

# STAGE 2
FROM alpine:latest

WORKDIR /root

RUN apk --no-cache add ca-certificates git && \
    mkdir ./config

COPY --from=builder ./app/app . 
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ./config/local.yaml ./config
COPY .env .

EXPOSE 8080

CMD ./app -config_path=./config/local.yaml
