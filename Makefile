SHELL=/usr/bin/env bash

export CGO_CFLAGS_ALLOW=-D__BLST_PORTABLE__
export CGO_CFLAGS=-D__BLST_PORTABLE__

GOVERSION:=$(shell go version | cut -d' ' -f 3 | cut -d. -f 2)
ifeq ($(shell expr $(GOVERSION) \< 17), 1)
$(warning Your Golang version is go 1.$(GOVERSION))
$(error Update Golang to version to at least 1.17.9)
endif

ldflags=-X=github.com/diwufeiwen/venus-metrics/version.CurrentCommit=+git.$(subst -,.,$(shell git describe --always --match=NeVeRmAtCh --dirty 2>/dev/null || git rev-parse --short HEAD 2>/dev/null))
ifneq ($(strip $(LDFLAGS)),)
	    ldflags+=-extldflags=$(LDFLAGS)
endif

GOFLAGS+=-ldflags="$(ldflags)"

## variables
BINS:=

all: build
.PHONY: all

build:
	rm -rf ./venus-metrics
	go build $(GOFLAGS) -o venus-metrics ./cmd/

.PHONY: build
BINS+=venus-metrics

clean:
	rm -rf $(BINS)
.PHONY: clean

dist-clean:
	git clean -xdff
	git submodule deinit --all -f
.PHONY: dist-clean
