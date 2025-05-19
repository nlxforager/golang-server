ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app ./cmd/product/makanplace/main.go


FROM debian:bookworm
RUN apk add --no-cache ca-certificates
COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]
