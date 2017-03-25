# These will be provided to the target
VERSION := 0.0.1
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-tags release -ldflags "-s -w -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

all: agent master rendezvous

agent: bin_dir
	cd agents/go && go build $(LDFLAGS) -o ../../bin/agent

master: bin_dir
	cd master && go build $(LDFLAGS) -o ../bin/master

rendezvous: bin_dir
	cd rendezvous && go build $(LDFLAGS) -o ../bin/rendezvous

bin_dir:
	mkdir -p bin
