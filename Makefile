ECR := 365769272576.dkr.ecr.us-west-2.amazonaws.com
AWS_REGION := ${AWS_DEFAULT_REGION}
AWS_ACCESS_KEY_ID := ${AWS_ACCESS_KEY_ID}
AWS_SECRET_ACCESS_KEY := ${AWS_SECRET_ACCESS_KEY}
CONFIG_BUCKET := dev-leveledup-api-config
CURRENT_DIR := $(shell pwd)
PROJECT := $(notdir $(CURRENT_DIR))
USER := $(notdir $(shell dirname $(CURRENT_DIR)))
CONTAINER_DIR := /go/src/github.com/$(USER)/$(PROJECT)
CONTAINER_DIR_CIRCLE := /go/src/github.com/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}
CIRCLECI := ${CIRCLECI}
PORTS := "8080:8080"

make = docker run -p $(PORTS) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e SECRETS_BUCKET_NAME=$(CONFIG_BUCKET) $(ECR)/$(PROJECT):latest &
.PHONY: default run build

VERSION := git-$(shell git rev-parse --short HEAD)

all: clean deploy

test: vendor build run

vendor:
ifndef glide
	curl https://glide.sh/get | sh
endif
	glide install

clean:
	rm -rf vendor/
	go clean

login:
ifeq ($(CIRCLECI), true)
	~/.local/bin/aws ecr get-login --no-include-email --region $(AWS_REGION) > login.sh
else
	aws ecr get-login --no-include-email --region $(AWS_REGION) > login.sh
endif
	bash login.sh
	rm login.sh

build: vendor
	docker build -t $(ECR)/$(PROJECT):$(VERSION) .
	docker tag $(ECR)/$(PROJECT):$(VERSION) $(ECR)/$(PROJECT):latest

push: login
	docker push $(ECR)/$(PROJECT):$(VERSION)
	docker push $(ECR)/$(PROJECT):latest

deploy: build push

image:
ifeq ($(CIRCLECI), true)
	docker build --rm=false -t ${CIRCLE_PROJECT_REPONAME}:$(VERSION).${CIRCLE_BUILD_NUM} .
	docker tag -f ${CIRCLE_PROJECT_REPONAME}:$(VERSION).${CIRCLE_BUILD_NUM} ${CIRCLE_PROJECT_REPONAME}:latest
else
	docker build -t $(ECR)/$(PROJECT):latest .
endif

run:
	docker run -p $(PORTS) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e SECRETS_BUCKET_NAME=$(CONFIG_BUCKET) $(ECR)/$(PROJECT):latest
