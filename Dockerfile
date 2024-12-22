ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-s' -o=./api ./cmd/api

FROM gcr.io/distroless/static

WORKDIR /app
COPY --from=builder /build/api ./api
CMD ["/app/api"]