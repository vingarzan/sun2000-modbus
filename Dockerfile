# syntax=docker/dockerfile:1

# Build
FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /sun2000-modbus

# Deploy
FROM alpine:latest

WORKDIR /

COPY --from=builder /sun2000-modbus /sun2000-modbus

USER nobody:nogroup

ENTRYPOINT ["/sun2000-modbus"]
