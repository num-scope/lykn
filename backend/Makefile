GO ?= go
CONFIG ?= config/config.yaml
DOCKER_IMAGE ?= lykn-server:dev
DEMO_OUT ?= ../tests/fixtures

# Windows (cmd + ezwinports/mingw make) vs Unix shell
ifeq ($(OS),Windows_NT)
	BINARY ?= server.exe
	RM = del /Q
	RMDIR = rmdir /S /Q
	NULL = NUL
	SET_CONFIG = set "LYKN_CONFIG=$(CONFIG)"&&
else
	BINARY ?= server
	RM = rm -f
	RMDIR = rm -rf
	NULL = /dev/null
	SET_CONFIG = LYKN_CONFIG=$(CONFIG)
endif

.PHONY: help run test build clean fmt demo docker-build

help:
	@echo Lykn server make targets:
	@echo   make run          Run the HTTP server with LYKN_CONFIG=$(CONFIG)
	@echo   make test         Run all Go tests
	@echo   make build        Build server binary to $(BINARY)
	@echo   make clean        Remove build output
	@echo   make fmt          Run gofmt on all tracked Go packages
	@echo   make demo         Generate complete demo fixtures in $(DEMO_OUT)
	@echo   make docker-build Build backend Docker image $(DOCKER_IMAGE)
	@echo.
	@echo Windows PowerShell alternative:  .\dev.ps1 help

run:
	$(SET_CONFIG) $(GO) run ./cmd/server

test:
	$(GO) test ./...

build:
	$(GO) build -o $(BINARY) ./cmd/server

clean:
ifeq ($(OS),Windows_NT)
	-if exist $(BINARY) $(RM) $(BINARY)
	-if exist server $(RM) server
	-if exist bin $(RMDIR) bin
else
	-$(RM) $(BINARY)
	-$(RMDIR) bin
endif

fmt:
	$(GO) fmt ./...

demo:
	$(GO) run ./cmd/demo $(DEMO_OUT)

docker-build:
	docker build -t $(DOCKER_IMAGE) .
