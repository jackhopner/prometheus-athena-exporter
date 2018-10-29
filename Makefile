DOCKERREPO = movio
ARTIFACT = go-app
IMAGE = prometheus-athena-exporter
OS = $(shell uname | tr [:upper:] [:lower:])

all: build

build: GOOS ?= ${OS}
build: GOARCH ?= amd64
build: clean
		@echo "building..."; \
		GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -o ${ARTIFACT} -a .

clean:
		rm -f ${ARTIFACT}

image: clean
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${ARTIFACT} -a .
		${SUDO} docker build -t $(DOCKERREPO)/$(CONTAINER):$(TAG) .

image-and-push: build clean image
		${SUDO} docker push $(DOCKERREPO)/$(IMAGE):$(TAG)

test:
		go test -v -cover                                   \
			.                                               \

test-short:
		go test -v -short					\
			.                                               \

run: build
		./${ARTIFACT}

git-clean:
		@if [ "$$(git rev-parse --abbrev-ref HEAD)" != "master" ]; then \
			echo "You are not on master branch, please merge before you deploy"; \
			exit 1; \
		else \
			if ! git diff-index --quiet HEAD -- ; then \
				echo "Uncommitted changed, please commit or stash then deploy again"; \
				exit 1; \
			fi\
		fi
git-tag-push:
		git tag $(TAG)
		git push --tags

check-tag:
ifndef TAG
		$(error TAG is undefined)
endif
