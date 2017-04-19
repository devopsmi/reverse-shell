# These will be provided to the target
VERSION := 0.0.1
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-tags release -ldflags "-s -w -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
ARCH ?= `go env GOHOSTARCH`
GOOS ?= `go env GOOS`

all: agent master rendezvous

package: clean all
	cd bin && tar -zcvf reverse-shell-$(VERSION)-$(GOOS)-$(ARCH).tar.gz agent master rendezvous

agent: build_dir
	cd agents/go && go build $(LDFLAGS) -o ../../bin/agent

master: build_dir
	cd master && go build $(LDFLAGS) -o ../bin/master

rendezvous: build_dir
	cd rendezvous && go build $(LDFLAGS) -o ../bin/rendezvous

build_dir:
	mkdir -p bin

clean:
	rm -rf bin/*
