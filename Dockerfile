ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-s -w' -o=./bin/api ./cmd/api

FROM gcr.io/distroless/static

WORKDIR /app
COPY --from=builder /build/bin/api ./api
ENTRYPOINT ["/app/api"]