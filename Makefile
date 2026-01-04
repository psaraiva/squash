# Dirs
DIR_WEB_BIN=./bin/web
DIR_WEB_CMD=./cmd/wasm

# Files
FILE_WEB_BIN=squash.wasm
FILE_WEB_INDEX=index.html
FILE_WEB_JS=wasm_exec.js

# Tools deployment
TOOL_TINYGO_ROOT=$(shell tinygo env TINYGOROOT)
TOOL_WASM_EXEC=$(TOOL_TINYGO_ROOT)/targets/$(FILE_WEB_JS)

# Tools developer
TOOL_GOTEST=go test
TOOL_MOCKERY=mockery

# Docker
DOCKER_IMAGE=squash-game
DOCKER_TAG=latest
DOCKER_PORT=8080

.PHONY: web-deploy-local web-build web-copy-files web-serve-start web-clean go-mock go-test go-test-wasm go-test-all docker-build docker-run docker-stop docker-deploy docker-clean

web-deploy-local: web-copy-files web-build web-serve-start

web-copy-files:
	@echo "Copiando wasm_exec.js..."
	cp $(TOOL_WASM_EXEC) $(DIR_WEB_BIN)
	cp $(DIR_WEB_CMD)/$(FILE_WEB_INDEX) $(DIR_WEB_BIN)

web-build: web-copy-files
	@echo "Compilando Squash para Wasm..."
	tinygo build -o $(DIR_WEB_BIN)/$(FILE_WEB_BIN) -target wasm $(DIR_WEB_CMD)

web-serve-start:
	@echo "Iniciando servidor Go..."
	go run server.go

web-clean:
	rm -f $(DIR_WEB_BIN)/$(FILE_WEB_BIN)
	rm -f $(DIR_WEB_BIN)/$(FILE_WEB_JS)
	rm -f $(DIR_WEB_BIN)/$(FILE_WEB_INDEX)

go-mock:
	@echo "Gerando mocks para as interfaces de renderização..."
	rm -rf internal/ports/mocks/
	$(TOOL_MOCKERY) --name=Renderer --dir=./internal/ports --output=./internal/ports/mocks --case=underscore

go-test:
	@echo "Executando testes unitários com cobertura..."
	$(TOOL_GOTEST) -v -cover ./internal/app/... ./pkg/adapters/input/web/...

go-test-wasm:
	@echo "Executando testes WASM..."
	tinygo test -target wasm ./pkg/adapters/output/web/...

go-test-all: go-test go-test-wasm

# Docker commands
docker-build:
	@echo "Construindo imagem Docker..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	@echo "Executando container Docker..."
	docker run -d --name squash-container -p $(DOCKER_PORT):$(DOCKER_PORT) $(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "Acesse: http://localhost:$(DOCKER_PORT)"

docker-stop:
	@echo "Parando container Docker..."
	docker stop squash-container || true
	docker rm squash-container || true

docker-deploy: docker-stop docker-build docker-run

docker-clean: docker-stop
	@echo "Removendo imagem Docker..."
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true
