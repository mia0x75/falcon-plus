CMD = agent aggregator graph hbs judge nodata transfer gateway api alarm updater exporter
TARGET = open-falcon
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
GOFMT ?= gofmt "-s"
GO_VERSION_MIN=1.10
VERSION := $(shell cat VERSION)
export GO111MODULE=on

all: fmt $(CMD) $(TARGET)

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w $(GOFILES)

install:
	@hash govendor > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	govendor sync

vet:
	go vet $(PACKAGES)

fmt:
	@bash ./genver.sh $(GO_VERSION_MIN)
	@$(GOFMT) -l -s -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	# get all go files and run go fmt on them
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

$(CMD):
	@echo "Compiling module $@ ..."
	@if [ $@ = "gateway" ]; then \
		sed -i -e "s/Transfer/Gateway/g" modules/transfer/g/g.go ; \
		find ./ -type f -name "g.go-e" -exec rm -f {} \; ; \
		go build -o bin/$@/falcon-$@ ./modules/transfer ; \
	elif [ $@ = "transfer" ]; then \
		sed -i -e "s/Gateway/Transfer/g" modules/transfer/g/g.go ; \
		find ./ -type f -name "g.go-e" -exec rm -f {} \; ; \
		go build -o bin/$@/falcon-$@ ./modules/$@ ; \
	else \
		go build -o bin/$@/falcon-$@ ./modules/$@ ; \
	fi
	@if [ "$(shell uname -s)" != "Darwin" ]; then strip bin/$@/falcon-$@; fi;

.PHONY: $(TARGET)
$(TARGET): $(GOFILES)
	@echo "Compiling module $@ ..."
	@go build -ldflags "-X main.GitCommit=`git rev-parse --short HEAD` -X main.Version=$(VERSION)" -o bin/open-falcon
	@if [ "$(shell uname -s)" != "Darwin" ]; then strip bin/$@; fi;

checkbin: bin/ config/

pack: checkbin
	@if [ -e out ] ; then rm -rf out; fi
	@mkdir out
	@$(foreach var,$(CMD),mkdir -p ./out/$(var);)
	@$(foreach var,$(CMD),cp ./config/$(var).json ./out/$(var)/$(var).json;)
	@$(foreach var,$(CMD),cp ./bin/$(var)/falcon-$(var) ./out/$(var);)
	@$(foreach var,$(CMD),upx --best ./out/$(var)/falcon-$(var)>/dev/null;)
	@cp -r ./modules/agent/public ./out/agent/
	@(cd ./out && ln -s ./agent/public/ ./public)
	@(cd ./out && mkdir -p ./agent/plugin && ln -s ./agent/plugin/ ./plugin)
	@cp -r ./modules/api/data ./out/api/
	@mkdir out/graph/data
	@bash ./config/confgen.sh
	@cp bin/$(TARGET) ./out/$(TARGET)
	@echo "Compressing executable binary for $(TARGET) ..."
	@upx --best ./out/$(TARGET)>/dev/null
	tar -C out -zcf open-falcon-v$(VERSION).tar.gz .
	@rm -rf out

.PHONY: test
test:
	@go test ./modules/api/test

clean:
	@rm -rf ./bin
	@rm -rf ./out
	@rm -rf ./$(TARGET)
	@rm -rf open-falcon-v$(VERSION).tar.gz

.PHONY: clean all agent aggregator graph hbs judge nodata transfer gateway api alarm updater exporter
