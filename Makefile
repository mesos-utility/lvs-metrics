default: help 
$(shell go get -u github.com/tools/godep)

COMMIT := $(shell git rev-parse HEAD 2> /dev/null || true)
GITHUB_SRC := github.com/mesos-utility
MODULE := lvs-metrics
CURDIR_LINK := $(CURDIR)/Godeps/_workspace/src/$(GITHUB_SRC)
export GOPATH := $(CURDIR)/Godeps/_workspace

## Make bin for lvs-metrics.
bin: ${CURDIR_LINK}
	#./control build
	go build -i -ldflags "-X g.Commit=${COMMIT}" -o lvs-metrics .

## Get godep and restore dep.
godep:
	@go get -u github.com/tools/godep
	GO15VENDOREXPERIMENT=0 GOPATH=`godep path` godep restore

$(CURDIR_LINK):
	mkdir -p $(CURDIR_LINK)
	ln -sfn $(CURDIR) $(CURDIR_LINK)/$(MODULE)

## Get vet go tools.
vet:
	go get golang.org/x/tools/cmd/vet

## Validate this go project.
validate:
	script/validate-gofmt
	#go vet ./...

## Run test case for this go project.
test:
	go test -v ./...

## Clean everything (including stray volumes).
clean:
#	find . -name '*.created' -exec rm -f {} +
	-rm -rf var
	-rm -f lvs-metrics

help: # Some kind of magic from https://gist.github.com/rcmachado/af3db315e31383502660
	$(info Available targets)
	@awk '/^[a-zA-Z\-\_0-9]+:/ {                                   \
		nb = sub( /^## /, "", helpMsg );                             \
		if(nb == 0) {                                                \
			helpMsg = $$0;                                             \
			nb = sub( /^[^:]*:.* ## /, "", helpMsg );                  \
		}                                                            \
		if (nb)                                                      \
			printf "\033[1;31m%-" width "s\033[0m %s\n", $$1, helpMsg; \
	}                                                              \
	{ helpMsg = $$0 }'                                             \
	width=$$(grep -o '^[a-zA-Z_0-9]\+:' $(MAKEFILE_LIST) | wc -L)  \
	$(MAKEFILE_LIST)

