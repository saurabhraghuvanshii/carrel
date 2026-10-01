BINARY := carrel
LDFLAGS := -s -w

.PHONY: build run test vet doctor check-packs dist clean

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

# One binary per platform, no other files needed.
dist:
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/carrel-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/carrel-linux-arm64 .
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/carrel-darwin-arm64 .
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/carrel-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/carrel-windows-amd64.exe .

clean:
	rm -rf dist $(BINARY)
