FROM golang:1.26.2 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/lykn-server ./cmd/server

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/lykn-server /app/lykn-server

EXPOSE 8080

ENTRYPOINT ["/app/lykn-server"]
