.PHONY: all
all: test build

.PHONY: test
test:
	go test ./...

.PHONY: vendor
vendor:
	go mod tidy -compat=1.23
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

.PHONY: fmt
fmt:
	nix-shell --run "nixfmt *.nix"
	nix-shell --run "prettier -w static/s.css"

.PHONY: shell
shell:
	nix-shell

.PHONY: check
check:
	make build && ./niksnut check

.PHONY: dev
dev:
	make build && ./niksnut -root=. httpd

.PHONY: run
run:
	make build && ./niksnut httpd
