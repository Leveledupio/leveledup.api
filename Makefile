ECR := 365769272576.dkr.ecr.us-west-2.amazonaws.com
AWS_REGION := ${AWS_DEFAULT_REGION}
CURRENT_DIR := $(shell pwd)
PROJECT := $(notdir $(CURRENT_DIR))
USER := $(notdir $(shell dirname $(CURRENT_DIR)))
CONTAINER_DIR := /go/src/github.com/$(USER)/$(PROJECT)
CONTAINER_DIR_CIRCLE := /go/src/github.com/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}
CIRCLECI := ${CIRCLECI}

VERSION := git-$(shell git rev-parse --short HEAD)

vendor:
	go get -u github.com/kardianos/govendor
	govendor init
	govendor get

clean:
	go clean

login:
	$(eval $(aws ecr get-login --no-include-email --region $(AWS_REGION)))

build:
	docker build -t $(ECR)/$(PROJECT):$(VERSION) .
	docker tag $(ECR)/$(PROJECT):$(VERSION) $(ECR)/$(PROJECT)\latest

push:
	make login
	docker push $(ECR)/$(PROJECT):$(VERSION)
	docker push $(ECR)/$(PROJECT):latest

deploy:
	make build
	make push

image:
ifeq ($(CIRCLECI), true)
	docker build --rm=false -t ${CIRCLE_PROJECT_REPONAME}:$(VERSION).${CIRCLE_BUILD_NUM} .
	docker tag -f ${CIRCLE_PROJECT_REPONAME}:$(VERSION).${CIRCLE_BUILD_NUM} ${CIRCLE_PROJECT_REPONAME}:latest
else
	docker build -t $(ECR)/$(PROJECT):latest .
endif