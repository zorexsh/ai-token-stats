.PHONY: build test cross

build:
	go build -o bin/ai-token-stats ./cmd/ai-token-stats

test:
	go test ./...

# Cross-compiles the six release binaries. CGO stays off so the SQLite
# driver does not need a C toolchain on the target.
cross:
	@mkdir -p dist
	@for spec in darwin/arm64 darwin/amd64 linux/amd64 linux/arm64 windows/amd64 windows/arm64; do \
		os=$${spec%/*}; arch=$${spec#*/}; ext=""; \
		[ "$$os" = windows ] && ext=".exe"; \
		echo "$$os/$$arch"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -ldflags "-s -w" -o "dist/ai-token-stats-$$os-$$arch$$ext" ./cmd/ai-token-stats || exit 1; \
	done
