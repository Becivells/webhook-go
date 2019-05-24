# This how we want to name the binary output
#
# use checkmake linter https://github.com/mrtazz/checkmake
# $ checkmake Makefile
# Copy https://github.com/XiaoMi/soar  author @martianzhang

BINARY=webhooks-go
GITREPO=github.com:Becivells/webhook-go
GOPATH ?= $(shell go env GOPATH)
# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif
PATH := ${GOPATH}/bin:$(PATH)
GCFLAGS=-gcflags "all=-trimpath=${GOPATH}"
LDFLAGS=-ldflags="-s -w"

# These are the values we want to pass for VERSION  and BUILD
BUILD_TIME=`date +%Y%m%d%H%M`
COMMIT_VERSION=`git rev-parse HEAD`

# colors compatible setting
COLOR_ENABLE=$(shell tput colors > /dev/null; echo $$?)
ifeq "$(COLOR_ENABLE)" "0"
CRED=$(shell echo "\033[91m")
CGREEN=$(shell echo "\033[92m")
CYELLOW=$(shell echo "\033[93m")
CEND=$(shell echo "\033[0m")
endif

HTTPPROXY=http://127.0.0.1:8118

.PHONY: all
all: | fmt build

.PHONY: go_version_check
GO_VERSION_MIN=1.11
# Parse out the x.y or x.y.z version and output a single value x*10000+y*100+z (e.g., 1.9 is 10900)
# that allows the three components to be checked in a single comparison.
VER_TO_INT:=awk '{split(substr($$0, match ($$0, /[0-9\.]+/)), a, "."); print a[1]*10000+a[2]*100+a[3]}'
go_version_check:
	@echo "$(CGREEN)Go version check ...$(CEND)"
	@if test $(shell go version | $(VER_TO_INT) ) -lt \
  	$(shell echo "$(GO_VERSION_MIN)" | $(VER_TO_INT)); \
  	then printf "go version $(GO_VERSION_MIN)+ required, found: "; go version; exit 1; \
		else echo "go version check pass";	fi

# Code format
.PHONY: fmt
fmt: go_version_check
	@echo "$(CGREEN)Run gofmt on all source files ...$(CEND)"
	@echo "gofmt -l -s -w ..."
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		gofmt -l -s -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

# Run golang test cases
.PHONY: test
test:
	@echo "$(CGREEN)Run all test cases ...$(CEND)"
	go test -race ./...
	@echo "test Success!"

# Rule golang test cases with `-update` flag
test-update:
	@echo "$(CGREEN)Run all test cases with -update flag ...$(CEND)"
	go test ./... -update
	@echo "test-update Success!"


.PHONY: cover
cover: test
	@echo "$(CGREEN)Run test cover check ...$(CEND)"
	go test -coverpkg=./... -coverprofile=coverage.data ./... | column -t
	go tool cover -html=coverage.data -o coverage.html
	go tool cover -func=coverage.data -o coverage.txt
	@tail -n 1 coverage.txt | awk '{sub(/%/, "", $$NF); \
		if($$NF < 80) \
			{print "$(CRED)"$$0"%$(CEND)"} \
		else if ($$NF >= 90) \
			{print "$(CGREEN)"$$0"%$(CEND)"} \
		else \
			{print "$(CYELLOW)"$$0"%$(CEND)"}}'

# Builds the project
build: fmt
	@echo "$(CGREEN)Building ...$(CEND)"
	@mkdir -p bin

	@ret=0 && for d in $$(go list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...); do \
		b=$$(basename $${d}) ; \
		go build ${GCFLAGS} -o bin/$${b} $$d || ret=$$? ; \
	done ; exit $$ret
	@echo "build Success!"

# Installs our project: copies binaries
install: build
	@echo "$(CGREEN)Install ...$(CEND)"
	go install ....
	@echo "install Success!"

.PHONY: release
release: build
	@echo "$(CGREEN)Cross platform building for release ...$(CEND)"
	@mkdir -p release
	@for GOOS in darwin linux windows; do \
		for GOARCH in amd64; do \
			for d in $$(go list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...); do \
				b=$$(basename $${d}) ; \
				if [ "$${GOOS}" = 'windows' ]; then\
				echo "Building $${b}.$${GOOS}-$${GOARCH}.exe ..."; \
				GOOS=$${GOOS} GOARCH=$${GOARCH} go build ${GCFLAGS} ${LDFLAGS} -v -o release/$${b}.$${GOOS}-$${GOARCH}.exe $$d 2>/dev/null ; \
				else \
				echo "Building $${b}.$${GOOS}-$${GOARCH} ..."; \
				GOOS=$${GOOS} GOARCH=$${GOARCH} go build ${GCFLAGS} ${LDFLAGS} -v -o release/$${b}.$${GOOS}-$${GOARCH} $$d 2>/dev/null ; \
				fi \
			done ; \
		done ;\
	done

.PHONY: init
init:
	@echo "$(CGREEN)go mod init $(GITREPO) ...$(CEND)"
	go mod init $(GITREPO)

.PHONY: tidy
tidy:
	@echo "$(CGREEN)go mod tidy $(GITREPO) ...$(CEND)"
	go mod tidy

.PHONY: vendor
verdor:
	@echo "$(CGREEN)go mod vendor $(GITREPO) ...$(CEND)"
	go mod vendor

.PHONY: verify
verify:
	@echo "$(CGREEN)go mod verify $(GITREPO) ...$(CEND)"
	go mod verify

.PHONY: graph
graph:
	@echo "$(CGREEN)go mod graph $(GITREPO) ...$(CEND)"

	go mod graph

.PHONY: edit
edit:
	@echo "$(CGREEN)go mod edit $(GITREPO) ...$(CEND)"
	go mod edit $(edit)

.PHONY: download
download:
	@echo "$(CGREEN)go mod download $(GITREPO) ...$(CEND)"
	go mod download

.PHONY: ptidy
ptidy: export http_proxy=$(HTTPPROXY)
ptidy: export https_proxy=$(HTTPPROXY)
ptidy:
	@echo "$(CGREEN)go mod tidy $(GITREPO) ...$(CEND)"
	go mod tidy

.PHONY: pdownload
pdownload: export http_proxy=$(HTTPPROXY)
pdownload: export https_proxy=$(HTTPPROXY)
pdownload:
	@echo "$(CGREEN)go mod download $(GITREPO) ...$(CEND)"
	go mod download

.PHONY: server
server: build
	@echo "$(CGREEN)go mod edit $(GITREPO) ...$(CEND)"
	bin/webhook-go -c webhooks.yaml

.PHONY: curl
curl:
	@echo "$(CGREEN)go mod edit $(GITREPO) ...$(CEND)"
	curl myip.ipip.net
# Cleans our projects: deletes binaries
.PHONY: clean

clean:
	@echo "$(CGREEN)Cleanup ...$(CEND)"
	go clean
	@for GOOS in darwin linux windows; do \
	    for GOARCH in 386 amd64; do \
	    if [ "$${GOOS}" = 'windows' ]; then\
			rm -f ${BINARY}.$${GOOS}-$${GOARCH}.exe ;\
		else\
			rm -f ${BINARY}.$${GOOS}-$${GOARCH};\
		fi\
		done ;\
	done
	rm -f ${BINARY} coverage.*
	find . -name "*.log" -delete
	git clean -fi

