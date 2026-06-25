set shell := ["bash", "-c"]

golangci_lint_version := "v2.12.2"
golangci_lint_bin := ".bin/golangci-lint"

default:
  @just --list

install-golangci-lint:
  @set -euo pipefail
  @mkdir -p .bin
  @GOBIN="$(pwd)/.bin" go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@{{golangci_lint_version}}

lint:
  @set -euo pipefail
  @if [ ! -x "{{golangci_lint_bin}}" ] || [ "$$({{golangci_lint_bin}} version --format short 2>/dev/null || true)" != "{{golangci_lint_version}}" ]; then \
    just install-golangci-lint; \
  fi
  @{{golangci_lint_bin}} run ./...
