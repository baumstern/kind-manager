FROM golang:1.15-alpine AS builder
WORKDIR /usr/src/app

# Copy project and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -v

# Multi stage build which reduces image size
FROM alpine:3.12

# Copy binary
COPY --from=builder /usr/src/app/kind-manager /usr/local/bin/

EXPOSE 15050
CMD ["kind-manager"]
