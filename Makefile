#
#  Makefile
#
#  The kickoff point for all project management commands.
#

GOCC := go

# Program version
VERSION := $(shell git describe --always --tags)

# Binary name for bintray
BIN_NAME=email

# Project owner for bintray
OWNER=gesquive

# Project name for bintray
PROJECT_NAME=email

# Project url used for builds
# examples: github.com, bitbucket.org
REPO_HOST_URL=github.com

# Grab the current commit
GIT_COMMIT=$(shell git rev-parse HEAD)

# Check if there are uncommited changes
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

# Use a local vendor directory for any dependencies; comment this out to
# use the global GOPATH instead
# GOPATH=$(PWD)

INSTALL_PATH=$(GOPATH)/${REPO_HOST_URL}/${OWNER}/${PROJECT_NAME}

FIND_DIST:=find * -type d -exec

default: build

.PHONY: help
help:
	@echo 'Management commands for ${PROJECT_NAME}:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	 awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo

.PHONY: build
build: ## Compile the project
	@echo "building ${OWNER} ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	${GOCC} build -ldflags "-X main.version=${VERSION} -X main.dirty=${GIT_DIRTY}" -o ${BIN_NAME}

.PHONY: install
install: build ## Install the binaries on this computer
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./${BIN_NAME} ${DESTDIR}/usr/local/bin/${BIN_NAME}

.PHONY: depends
depends: ## Download project dependencies
	${GOCC} get -u github.com/Masterminds/glide
	glide install

.PHONY: test
test: ## Run golang tests
	${GOCC} test ./...

.PHONY: clean
clean: ## Clean the directory tree of artifacts
	${GOCC} clean
	rm -f ./${BIN_NAME}.test
	rm -f ./${BIN_NAME}
	rm -rf ./dist

.PHONY: bootstrap-dist
bootstrap-dist:
	${GOCC} get -u github.com/mitchellh/gox

.PHONY: build-all
build-all: bootstrap-dist
	gox -verbose \
	-ldflags "-X main.version=${VERSION} -X main.dirty=${GIT_DIRTY}" \
	-os="linux darwin windows " \
	-arch="amd64 386" \
	-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" .

.PHONY: dist
dist: build-all ## Cross compile the full distribution
	cd dist && \
	$(FIND_DIST) cp ../LICENSE {} \; && \
	$(FIND_DIST) cp ../README.md {} \; && \
	$(FIND_DIST) tar -zcf ${PROJECT_NAME}-${VERSION}-{}.tar.gz {} \; && \
	$(FIND_DIST) zip -r ${PROJECT_NAME}-${VERSION}-{}.zip {} \; && \
	cd ..

.PHONY: fmt
fmt: ## Reformat the source tree with gofmt
	find . -name '*.go' -not -path './.vendor/*' -exec gofmt -w=true {} ';'

.PHONY: link
link: ## Symlink this source tree into the GOPATH
	if [ ! $(INSTALL_PATH) -ef . ]; then \
		mkdir -p `dirname $(INSTALL_PATH)`; \
		ln -s $(PWD) $(INSTALL_PATH); \
	fi
