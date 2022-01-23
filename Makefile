GOVERSION := $(shell go env GOVERSION)

BIN := bin

LINTER := $(BIN)/golangci-lint
LINTER_PACKAGE := github.com/golangci/golangci-lint/cmd/golangci-lint
LINTER_VERSION := v1.43.0
LINTER_VERSIONED := $(LINTER)_$(LINTER_VERSION)_$(GOVERSION)

lint: $(LINTER)
	$(LINTER) run

$(LINTER): $(LINTER_VERSIONED) | $(BIN)
	ln -sf $(abspath $<) $@

$(LINTER_VERSIONED): | $(BIN)
	find $(BIN) -name 'golangci-lint' -exec rm {} +
	find $(BIN) -name '$(notdir $(LINTER))_*' -exec rm {} +
	GOBIN=$(abspath $(BIN)) go install $(LINTER_PACKAGE)@$(LINTER_VERSION)
	mv $(BIN)/golangci-lint $@

$(BIN):
	mkdir -p $@
