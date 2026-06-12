FROM golang:1.26.2 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY api ./api
COPY cmd ./cmd
COPY config ./config
COPY database ./database
COPY internal ./internal
COPY migrations ./migrations
COPY pkg ./pkg

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/lykn-server ./cmd/server

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata wget \
    && addgroup -S lykn \
    && adduser -S -D -H -h /app -s /sbin/nologin -G lykn lykn \
    && mkdir -p /app/config \
    && chown -R lykn:lykn /app

COPY --from=builder --chown=lykn:lykn /out/lykn-server /app/lykn-server

ENV GIN_MODE=release
ENV LYKN_CONFIG=/app/config/config.yaml

USER lykn

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8080/health >/dev/null || exit 1

ENTRYPOINT ["/app/lykn-server"]
