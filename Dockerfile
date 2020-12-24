FROM golang:1.15-alpine AS builder
WORKDIR /usr/src/app

# Copy project and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -v

# Multi stage build which reduces image size
FROM alpine:3.12
WORKDIR /usr/src/app

# Copy binary
COPY --from=builder /usr/src/app/kind-manager .

EXPOSE 15050
CMD ["./kind-manager"]
