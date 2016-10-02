# Makefile to build fluentd Docker image

PROJECT=$(shell gcloud config list project | awk 'FNR==2 { print $$3 }')
VERSION=$(shell cat VERSION)

all: image

# Build a docker image for the local architecture
image:
	docker build -t asia.gcr.io/$(PROJECT)/aggremod-fluentd:$(VERSION) .

push:
	$(SHELL) -ec 'gcloud docker -- push asia.gcr.io/$(PROJECT)/aggremod-fluentd:$(VERSION)'
