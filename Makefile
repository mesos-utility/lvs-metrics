default: help 

NAME := lvs-metrics
COMMIT := $(shell git rev-parse HEAD 2> /dev/null || true)
GITHUB_SRC := github.com/mesos-utility
CURDIR_LINK := $(CURDIR)/Godeps/_workspace/src/$(GITHUB_SRC)
export GOPATH := $(CURDIR)/Godeps/_workspace


MKDIR	= mkdir
INSTALL	= install
BIN		= $(BUILD_ROOT)
MAN		= $(BIN)
VERSION	= $(shell git describe --tags --abbrev=0 2> /dev/null)
RELEASE	= 0
RPMSOURCEDIR	= $(shell rpm --eval '%_sourcedir')
RPMSPECDIR	= $(shell rpm --eval '%_specdir')
RPMBUILD = $(shell				\
	if [ -x /usr/bin/rpmbuild ]; then	\
		echo "/usr/bin/rpmbuild";	\
	else					\
		echo "/bin/rpm";		\
	fi )

## Make bin for lvs-metrics.
bin: ${CURDIR_LINK}
	go build -i -ldflags "-X github.com/mesos-utility/lvs-metrics/g.Version=${VERSION}" -o lvs-metrics .

## Get godep and restore dep.
godep:
	@go get -u github.com/tools/godep
	GO15VENDOREXPERIMENT=0 GOPATH=`godep path` godep restore

$(CURDIR_LINK):
	rm -rf $(CURDIR_LINK)
	mkdir -p $(CURDIR_LINK)
	ln -sfn $(CURDIR) $(CURDIR_LINK)/$(NAME)

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
	-rm -rf var
	-rm -f ${NAME}
	-rm -f ${NAME}-*.tar.gz
	-rm -f ${NAME}.spec


# Rpm install action for lvs-metrics rpm package build.
rpm-install:
	if [ ! -d $(BIN) ]; then $(MKDIR) -p $(BIN); fi
	$(INSTALL) -m 0755 $(NAME) $(BIN)
	$(INSTALL) -m 0644 cfg.json $(BIN)
	$(INSTALL) -m 0755 control $(BIN)
	[ -d $(MAN) ] || $(MKDIR) -p $(MAN)
	$(INSTALL) -m 0644 README.md $(MAN)

dist: clean
	sed -e "s/@@VERSION@@/$(VERSION)/g" \
		-e "s/@@RELEASE@@/$(RELEASE)/g" \
		< lvs-metrics.spec.in > lvs-metrics.spec
	rm -f $(NAME)-$(VERSION)
	rm -f cfg.json
	cp cfg.example.json cfg.json
	ln -s . $(NAME)-$(VERSION)
	tar czvf $(NAME)-$(VERSION).tar.gz			\
		--exclude CVS --exclude .git --exclude TAGS		\
		--exclude $(NAME)-$(VERSION)/$(NAME)-$(VERSION)	\
		--exclude $(NAME)-$(VERSION).tar.gz			\
		$(NAME)-$(VERSION)/*
	rm -f $(NAME)-$(VERSION)

rpms: dist
	cp $(NAME)-$(VERSION).tar.gz $(RPMSOURCEDIR)/
	cp $(NAME).spec $(RPMSPECDIR)/
	$(RPMBUILD) -ba $(RPMSPECDIR)/$(NAME).spec

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
