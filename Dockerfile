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

# HTTP_IP - The IP address to listen on. Defaults to 0.0.0.0.
# HTTP_PORT - The port to listen on. Defaults to 8080.
# MODBUS_IP - The IP address of the modbus device. This is required.
# MODBUS_PORT - The port of the modbus device. Defaults to 502.
# MODBUS_TIMEOUT - The timeout for modbus requests. Defaults to 5 seconds.
# MODBUS_SLEEP - The sleep time between modbus requests. Defaults to 5 seconds.
# MODBUS_SLAVE_ID - The slave ID of the modbus device. Defaults to 1.
WORKDIR /

COPY --from=builder /sun2000-modbus /sun2000-modbus

USER nobody:nogroup

EXPOSE 8080

ENTRYPOINT ["/sun2000-modbus"]
