# installs the server binary locally

BIN_DIR := $(shell go env GOPATH)/bin
PACKR_TOOL := $(BIN_DIR)/packr2

$(PACKR_TOOL):
	go get github.com/gobuffalo/packr/v2/packr2

.PHONY: install

install: $(PACKR_TOOL)
	GO111MODULE=on $(PACKR_TOOL) install
	GO111MODULE=on $(PACKR_TOOL) clean
