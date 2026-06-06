FROM golang:1.23-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /build/conduit ./cmd/gateway/

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /build/conduit /usr/local/bin/conduit
COPY config.example.yaml /etc/conduit/config.yaml

EXPOSE 8080
ENTRYPOINT ["conduit"]
CMD ["--config", "/etc/conduit/config.yaml"]
