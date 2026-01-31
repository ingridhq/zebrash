.PHONY: test-arm64

# Run Go tests in Docker ARM64 container
test-arm64:
	docker run --rm --platform linux/arm64 -v $(PWD):/app -w /app golang:1.25 go test ./...
