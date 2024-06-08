.PHONY: all
all: test build

.PHONY: test
test:
	go test ./...

vendor:
	go mod tidy -compat=1.22
	go mod vendor

.PHONY: nix
nix:
	nix-build -A niksnut

.PHONY: build
build:
	go build

.PHONY: nixpin
nixpin:
	#%nix-shell --run "npins -d ./build init"
	nix-shell --run "npins -d ./build update"

.PHONY: nixfmt
nixfmt:
	nix-shell --run "nixfmt *.nix"

