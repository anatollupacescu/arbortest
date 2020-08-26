# installs the server binary locally
BIN_DIR := $(shell go env GOPATH)/bin

PACKR_TOOL ?= $(BIN_DIR)/packr2

$(PACKR_TOOL):
	go get github.com/gobuffalo/packr/v2/packr2

ARBOR_SERVER ?= $(BIN_DIR)/arbortest

$(ARBOR_SERVER): $(PACKR_TOOL)
	GO111MODULE=on $(PACKR_TOOL) install
	GO111MODULE=on $(PACKR_TOOL) clean

.PHONY: clean-server

clean-server:
	@rm $(ARBOR_SERVER)

.PHONY: install-server

install-server: $(ARBOR_SERVER) 

ARBOR_GEN ?= $(BIN_DIR)/arborgen

$(ARBOR_GEN):
	go install ./arborgen

.PHONY: install-gen

install-gen: $(ARBOR_GEN)

.PHONY: test

test: 
	go test -v -count=1 ./...
