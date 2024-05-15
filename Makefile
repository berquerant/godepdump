GOMOD = go mod
GOBUILD = go build -trimpath -v
GOTEST = go test -v -cover -race

ROOT = $(shell git rev-parse --show-toplevel)
BIN = dist/godepdump

.PHONY: $(BIN)
$(BIN):
	$(GOBUILD) -o $@

.PHONY: test
test:
	$(GOTEST) ./...

.PHONY: init
init:
	$(GOMOD) tidy

.PHONY: generate
generate: clean-generated
	go generate ./...

.PHONY: clean-generated
clean-generated:
	find . -name "*_generated.go" -type f -delete

.PHONY: vuln
vuln:
	go run golang.org/x/vuln/cmd/govulncheck ./...

.PHONY: vet
vet:
	go vet ./...

DOCKER_RUN = docker run --rm -v "$(ROOT)":/usr/src/myapp -w /usr/src/myapp
DOCKER_GO_IMAGE = golang:1.21
DOCKER_LINT_IMAGE = golangci/golangci-lint:v1.54.2

.PHONY: docker-test
docker-test:
	$(DOCKER_RUN) $(DOCKER_GO_IMAGE) $(GOTEST) ./...

.PHONY: docker-dist
docker-dist:
	$(DOCKER_RUN) $(DOCKER_GO_IMAGE) $(GOBUILD) -o $(BIN) $(CMD)
