# Makefile for generating Go code from proto files using buf and ConnectRPC

# Tools
BUF := buf
GO := go

.PHONY: all gen tools clean test test-unit test-integration test-internal

all: gen

# Install required tools
tools:
	$(GO) install github.com/bufbuild/buf/cmd/buf@latest


# Generate all proto files using buf
gen:
	cd proto && $(BUF) generate

# Generate code from filesystem proto
gen-filesystem:
	cd proto && $(BUF) generate --path tool/filesystem

# Generate code from process proto  
gen-process:
	cd proto && $(BUF) generate --path tool/process

# Remove generated files
clean:
	rm -rf pb/

clean-filesystem:
	rm -f pb/filesystem/filesystem.pb.go
	rm -rf pb/filesystem/filesystemconnect

clean-process:
	rm -f pb/process/process.pb.go
	rm -rf pb/process/processconnect

# Test commands
test:
	$(GO) test ./...

test-unit:
	$(GO) test -short ./...

test-integration:
	$(GO) test ./test/...

# Run internal tests (requires internal network access)
test-internal:
	$(GO) test -tags=internal ./test/...