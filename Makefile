

.PHONY: build
build:  ## Build the binaries
	goreleaser release --snapshot --clean