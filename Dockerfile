# syntax=docker/dockerfile:1.7

ARG NODE_IMAGE=node:24-alpine
ARG GO_IMAGE=golang:1.26.2
ARG RUNTIME_IMAGE=alpine:3.22
ARG PNPM_VERSION=10.33.2
ARG VITE_API_BASE=/api/v1

FROM --platform=$BUILDPLATFORM ${NODE_IMAGE} AS frontend-deps

ARG PNPM_VERSION
ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG ALL_PROXY
ARG NO_PROXY
ARG http_proxy
ARG https_proxy
ARG all_proxy
ARG no_proxy

WORKDIR /src/frontend

ENV PNPM_HOME=/pnpm
ENV PATH=${PNPM_HOME}:${PATH}

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack enable && corepack prepare "pnpm@${PNPM_VERSION}" --activate
RUN --mount=type=cache,id=lykn-pnpm-store,target=/pnpm/store \
    pnpm config set store-dir /pnpm/store && \
    pnpm install --frozen-lockfile

FROM frontend-deps AS frontend-builder

ARG VITE_API_BASE

COPY frontend/ ./
ENV VITE_API_BASE=${VITE_API_BASE}
RUN pnpm build

FROM --platform=$BUILDPLATFORM ${GO_IMAGE} AS backend-deps

ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG ALL_PROXY
ARG NO_PROXY
ARG http_proxy
ARG https_proxy
ARG all_proxy
ARG no_proxy

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,id=lykn-go-mod,target=/go/pkg/mod \
    go mod download

FROM backend-deps AS backend-builder

ARG TARGETOS
ARG TARGETARCH

COPY api ./api
COPY cmd ./cmd
COPY config ./config
COPY database ./database
COPY internal ./internal
COPY migrations ./migrations
COPY pkg ./pkg
RUN --mount=type=cache,id=lykn-go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=lykn-go-build,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w -buildid=" -o /out/lykn-server ./cmd/server

FROM ${RUNTIME_IMAGE} AS runtime

LABEL org.opencontainers.image.title="Lykn"
LABEL org.opencontainers.image.description="Lykn management application with bundled web UI"
LABEL org.opencontainers.image.licenses="Apache-2.0"

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S lykn && \
    adduser -S -D -H -h /app -s /sbin/nologin -G lykn lykn && \
    mkdir -p /app/config /app/frontend/dist && \
    chown -R lykn:lykn /app

COPY --from=backend-builder --chown=lykn:lykn /out/lykn-server /app/lykn-server
COPY --from=frontend-builder --chown=lykn:lykn /src/frontend/dist /app/frontend/dist

ENV GIN_MODE=release
ENV LYKN_CONFIG=/app/config/config.yaml
ENV LYKN_WEB_DIST=/app/frontend/dist

USER lykn

EXPOSE 8080

STOPSIGNAL SIGTERM

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8080/health >/dev/null || exit 1

ENTRYPOINT ["/app/lykn-server"]
