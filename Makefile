GO ?= go
CONFIG ?= config/config.yaml
BINARY ?= server
DEMO_OUT ?= tests/fixtures

.PHONY: help run test build clean fmt demo docker-build

help:
	@printf "Lykn server make targets:\n"
	@printf "  make run          Run the HTTP server with LYKN_CONFIG=%s\n" "$(CONFIG)"
	@printf "  make test         Run all Go tests\n"
	@printf "  make build        Build server binary to %s\n" "$(BINARY)"
	@printf "  make clean        Remove build output\n"
	@printf "  make fmt          Run gofmt on all tracked Go packages\n"
	@printf "  make demo         Generate complete demo fixtures in %s\n" "$(DEMO_OUT)"
	@printf "  make docker-build Build Docker Compose images\n"

run:
	LYKN_CONFIG=$(CONFIG) $(GO) run ./cmd/server

test:
	$(GO) test ./...

build:
	$(GO) build -o $(BINARY) ./cmd/server

clean:
	rm -f $(BINARY)
	rm -rf bin

fmt:
	$(GO) fmt ./...

demo:
	$(GO) run ./cmd/demo $(DEMO_OUT)

docker-build:
	docker compose build
