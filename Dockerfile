FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /build/conduit github.com/muhammad6535/conduit/cmd/gateway

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /build/conduit /usr/local/bin/conduit
COPY config.example.yaml /etc/conduit/config.example.yaml
COPY config.railway.yaml /etc/conduit/config.yaml

EXPOSE 8080
ENTRYPOINT ["conduit"]
CMD ["--config", "/etc/conduit/config.yaml"]
