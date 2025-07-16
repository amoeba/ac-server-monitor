ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o monitor .

FROM debian:bookworm

WORKDIR /app
COPY --from=builder /build/monitor .
COPY --from=builder /build/static ./static
COPY --from=builder /build/templates ./templates

# Needed for fetching servers list from GitHub
RUN apt-get update && apt-get install -y ca-certificates

CMD ["/app/monitor"]
