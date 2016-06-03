default: help

COMMIT := $(shell git rev-parse HEAD 2> /dev/null || true)

## Make bin for lvs-metrics.
bin:
	#./control build
	GOPATH=`godep path`
	go build -i -ldflags "-X g.Commit=${COMMIT}" -o lvs-metrics .

## Get godep and restore dep.
godep:
	@go get -u github.com/tools/godep
	GO15VENDOREXPERIMENT=0 GOPATH=`godep path` godep restore

## Get vet go tools.
vet:
	go get golang.org/x/tools/cmd/vet

## Validate this go project.
validate:
	script/validate-gofmt
	#go vet ./...

## Run test case for this go project.
test:
	go list ./... | grep -v 'vendor' | xargs -L1 go test -v

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

