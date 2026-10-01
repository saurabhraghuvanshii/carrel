BINARY := carrel
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: build run test vet doctor check-packs dist clean site site-serve

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

run: build
	./$(BINARY)

test:
	go test ./...

vet:
	go vet ./...

doctor: build
	./$(BINARY) doctor

# Runs the reference solutions in tools/refs against every pack, in both languages.
check-packs:
	go run ./cmd/packcheck

# One archive per platform plus checksums.txt, named as the install scripts expect:
# carrel_<os>_<arch>.tar.gz (or .zip on Windows) holding the binary, LICENSE and README.md.
dist:
	rm -rf dist
	@set -e; for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; name=carrel_$${os}_$${arch}; dir=dist/$$name; \
		bin=carrel; [ $$os = windows ] && bin=carrel.exe; \
		mkdir -p $$dir; cp LICENSE README.md $$dir/; \
		echo "build $$name"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags "$(LDFLAGS)" -o $$dir/$$bin .; \
		if [ $$os = windows ]; then (cd $$dir && zip -q -X ../$$name.zip $$bin LICENSE README.md); \
		else tar -czf dist/$$name.tar.gz -C $$dir $$bin LICENSE README.md; fi; \
		rm -rf $$dir; \
	done
	cd dist && sha256sum carrel_* > checksums.txt

# The public website, built from site/src into site/dist.
site:
	go run ./tools/sitebuild -src site/src -out site/dist

site-serve:
	go run ./tools/sitebuild -src site/src -out site/dist -serve 127.0.0.1:8080

clean:
	rm -rf dist site/dist $(BINARY)
