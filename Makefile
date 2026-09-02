

.PHONY: install-tools
install-tools:  ## Install pinned tool versions via mise
	mise install

.PHONY: build
build: install-tools  ## Build the binaries
	mise exec -- goreleaser release --snapshot --clean