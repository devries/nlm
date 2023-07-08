BINARY := nlmserver
SOURCE := $(wildcard **/*.go) go.mod go.sum
TEMPLATES := $(wildcard nlmserver/templates/*.html)
ARTICLES := $(wildcard articles/*.txt)
STATIC := $(wildcard nlmserver/static/*)
ALLSOURCE := $(SOURCE) $(TEMPLATES) $(ARTICLES) $(STATIC)
.PHONY: clean build test
.DEFAULT: build

build/darwin/$(BINARY): $(ALLSOURCE)
	mkdir -p build/darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $@ ./$(BINARY)

build/darwinarm/$(BINARY): $(SOURCE)
	mkdir -p build/darwinarm
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $@ ./$(BINARY)

build/linux/$(BINARY): $(SOURCE)
	mkdir -p build/linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ ./$(BINARY)

build/linuxarmhf/$(BINARY): $(SOURCE)
	mkdir -p build/linuxarmhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o $@ ./$(BINARY)

build/linuxarm64/$(BINARY): $(SOURCE)
	mkdir -p build/linuxarm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $@ ./$(BINARY)

build/darwinuniversal/$(BINARY): build/darwin/$(BINARY) build/darwinarm/$(BINARY)
	mkdir -p build/darwinuniversal
	lipo -create -output $@ $^

build: build/darwinuniversal/$(BINARY) build/linux/$(BINARY) build/linuxarmhf/$(BINARY) build/linuxarm64/$(BINARY)

clean:
	rm -rf build || true

test:
	go test -cover .
